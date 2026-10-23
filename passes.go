package canvas3d

import (
	"errors"
	"fmt"

	"github.com/Zyko0/go-sdl3/sdl"
)

const hdrFormat = sdl.GPU_TEXTUREFORMAT_R16G16B16A16_FLOAT

type framePass struct {
	name string
	run  func(*Canvas) error
}

// ensureTargets makes the HDR color, depth and shadow
// textures match the window size.
func (canvas *Canvas) ensureTargets() error {
	rend := &canvas.rend
	device := canvas.gpu.device

	w := max(1, uint32(float32(rend.width)*rend.renderScale))
	h := max(1, uint32(float32(rend.height)*rend.renderScale))

	if rend.hdr == nil || rend.depth == nil || rend.renderW != w || rend.renderH != h {
		if rend.hdr != nil {
			device.ReleaseTexture(rend.hdr)
			rend.hdr = nil
		}
		if rend.depth != nil {
			device.ReleaseTexture(rend.depth)
			rend.depth = nil
		}

		hdr, err := device.CreateTexture(
			&sdl.GPUTextureCreateInfo{
				Type:              sdl.GPU_TEXTURETYPE_2D,
				Format:            hdrFormat,
				Usage:             sdl.GPU_TEXTUREUSAGE_COLOR_TARGET | sdl.GPU_TEXTUREUSAGE_SAMPLER,
				Width:             w,
				Height:            h,
				LayerCountOrDepth: 1,
				NumLevels:         1,
			},
		)
		if err != nil {
			return err
		}
		rend.hdr = hdr

		depth, err := device.CreateTexture(
			&sdl.GPUTextureCreateInfo{
				Type:              sdl.GPU_TEXTURETYPE_2D,
				Format:            depthFormat,
				Usage:             sdl.GPU_TEXTUREUSAGE_DEPTH_STENCIL_TARGET,
				Width:             w,
				Height:            h,
				LayerCountOrDepth: 1,
				NumLevels:         1,
			},
		)
		if err != nil {
			return err
		}
		rend.depth = depth
		rend.renderW, rend.renderH = w, h
	}

	size := uint32(rend.shadowSize)
	if rend.shadowTex == nil || rend.shadowTexSize != size {
		if rend.shadowTex != nil {
			device.ReleaseTexture(rend.shadowTex)
			rend.shadowTex = nil
		}
		shadow, err := device.CreateTexture(
			&sdl.GPUTextureCreateInfo{
				Type:              sdl.GPU_TEXTURETYPE_2D,
				Format:            depthFormat,
				Usage:             sdl.GPU_TEXTUREUSAGE_DEPTH_STENCIL_TARGET | sdl.GPU_TEXTUREUSAGE_SAMPLER,
				Width:             size,
				Height:            size,
				LayerCountOrDepth: 1,
				NumLevels:         1,
			},
		)
		if err != nil {
			return err
		}
		rend.shadowTex, rend.shadowTexSize = shadow, size
	}
	return nil
}

// shadowPass renders the scene's depth from the first light.
func (canvas *Canvas) shadowPass() error {
	rend := &canvas.rend
	canvas.updateShadowMatrix()

	pipeline, err := canvas.getPipeline(pipelineKey{shadow: true})
	if err != nil {
		return err
	}

	pass := rend.cmd.BeginRenderPass(
		nil,
		&sdl.GPUDepthStencilTargetInfo{
			Texture:        rend.shadowTex,
			ClearDepth:     1.0,
			LoadOp:         sdl.GPU_LOADOP_CLEAR,
			StoreOp:        sdl.GPU_STOREOP_STORE,
			StencilLoadOp:  sdl.GPU_LOADOP_DONT_CARE,
			StencilStoreOp: sdl.GPU_STOREOP_DONT_CARE,
		},
	)
	pass.BindGraphicsPipeline(pipeline)

	for _, item := range rend.draws {
		if item.material != nil && item.material.blend {
			continue
		}
		pass.BindVertexBuffers([]sdl.GPUBufferBinding{{Buffer: item.mesh.vertexBuffer}})
		pass.BindIndexBuffer(&sdl.GPUBufferBinding{Buffer: item.mesh.indexBuffer}, sdl.GPU_INDEXELEMENTSIZE_32BIT)

		mvp := Mul(rend.lightViewProj, item.transform)
		rend.cmd.PushVertexUniformData(0, asBytes(&mvp))
		pass.DrawIndexedPrimitives(item.mesh.indexCount, 1, 0, 0, 0)
	}
	pass.End()

	rend.shadowActive = true
	return nil
}

// scenePass draws every object, lit, into the HDR target.
func (canvas *Canvas) scenePass() error {
	rend := &canvas.rend

	pass := rend.cmd.BeginRenderPass(
		[]sdl.GPUColorTargetInfo{{
			Texture:    rend.hdr,
			ClearColor: rend.clearColor,
			LoadOp:     sdl.GPU_LOADOP_CLEAR,
			StoreOp:    sdl.GPU_STOREOP_STORE,
		}},
		&sdl.GPUDepthStencilTargetInfo{
			Texture:        rend.depth,
			ClearDepth:     1.0,
			LoadOp:         sdl.GPU_LOADOP_CLEAR,
			StoreOp:        sdl.GPU_STOREOP_DONT_CARE,
			StencilLoadOp:  sdl.GPU_LOADOP_DONT_CARE,
			StencilStoreOp: sdl.GPU_STOREOP_DONT_CARE,
		},
	)
	err := canvas.recordDraws(pass)
	pass.End()
	return err
}

// tonemapPass turns the HDR image into what goes on screen: exposure, tone
// curve, and (through the sRGB swapchain) the final color encoding. It also
// scales the image up or down when the render scale is not 1.
func (canvas *Canvas) tonemapPass() error {
	rend := &canvas.rend

	pipeline, err := canvas.getPipeline(pipelineKey{colorFormat: rend.colorFmt, tonemap: true})
	if err != nil {
		return err
	}

	pass := rend.cmd.BeginRenderPass(
		[]sdl.GPUColorTargetInfo{{
			Texture:    rend.swapchain,
			ClearColor: rend.clearColor,
			LoadOp:     sdl.GPU_LOADOP_CLEAR,
			StoreOp:    sdl.GPU_STOREOP_STORE,
		}},
		nil,
	)
	pass.BindGraphicsPipeline(pipeline)
	pass.BindFragmentSamplers([]sdl.GPUTextureSamplerBinding{{Texture: rend.hdr, Sampler: rend.hdrSampler}})

	params := [4]float32{rend.exposure, float32(rend.tonemapper), 0, 0}
	rend.cmd.PushFragmentUniformData(0, asBytes(&params))

	pass.DrawPrimitives(3, 1, 0, 0)
	pass.End()
	return nil
}

func (canvas *Canvas) createTonemapPipeline(key pipelineKey) (*sdl.GPUGraphicsPipeline, error) {
	vs, err := canvas.loadShader(tvs_name, sdl.GPU_SHADERSTAGE_VERTEX, 0, 0)
	if err != nil {
		return nil, err
	}
	defer canvas.gpu.device.ReleaseShader(vs)

	fs, err := canvas.loadShader(tfs_name, sdl.GPU_SHADERSTAGE_FRAGMENT, 1, 1)
	if err != nil {
		return nil, err
	}
	defer canvas.gpu.device.ReleaseShader(fs)

	pipeline, err := canvas.gpu.device.CreateGraphicsPipeline(
		&sdl.GPUGraphicsPipelineCreateInfo{
			VertexShader:   vs,
			FragmentShader: fs,
			PrimitiveType:  sdl.GPU_PRIMITIVETYPE_TRIANGLELIST,
			RasterizerState: sdl.GPURasterizerState{
				FillMode: sdl.GPU_FILLMODE_FILL,
				CullMode: sdl.GPU_CULLMODE_NONE,
			},
			TargetInfo: sdl.GPUGraphicsPipelineTargetInfo{
				ColorTargetDescriptions: []sdl.GPUColorTargetDescription{{Format: key.colorFormat}},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	canvas.rend.pipelines[key] = pipeline
	
	return pipeline, nil
}

func (canvas *Canvas) recordDraws(pass *sdl.GPURenderPass) error {
	rend := &canvas.rend

	cam := &rend.camera
	aspect := float32(rend.width) / float32(rend.height)
	viewProj := Mul(perspective(cam.FOV, aspect, cam.Near, cam.Far), lookAt(cam.Position, cam.Target))
	frame := canvas.buildFrameUniforms()

	pick := func(texture, fallback *Texture) sdl.GPUTextureSamplerBinding {
		if texture == nil {
			texture = fallback
		}
		return sdl.GPUTextureSamplerBinding{Texture: texture.texture, Sampler: rend.sampler}
	}

	var current pipelineKey
	bound := false

	for _, item := range rend.draws {
		material := item.material
		if material == nil {
			material = &defaultMaterial
		}

		key := pipelineKey{colorFormat: hdrFormat, doubleSided: material.doubleSided, blend: material.blend}
		if !bound || key != current {
			pipeline, err := canvas.getPipeline(key)
			if err != nil {
				return err
			}
			pass.BindGraphicsPipeline(pipeline)
			rend.cmd.PushFragmentUniformData(1, asBytes(&frame))
			current, bound = key, true
		}

		pass.BindVertexBuffers([]sdl.GPUBufferBinding{{Buffer: item.mesh.vertexBuffer}})
		pass.BindIndexBuffer(&sdl.GPUBufferBinding{Buffer: item.mesh.indexBuffer}, sdl.GPU_INDEXELEMENTSIZE_32BIT)

		draw := drawUniforms{mvp: Mul(viewProj, item.transform), model: item.transform}
		rend.cmd.PushVertexUniformData(0, asBytes(&draw))
		rend.cmd.PushFragmentUniformData(0, asBytes(&material.uniforms))

		bindings := [6]sdl.GPUTextureSamplerBinding{
			pick(material.baseColorTex, rend.whiteTexture),
			pick(material.metalRoughTex, rend.whiteTexture),
			pick(material.normalTex, rend.flatNormal),
			pick(material.occlusionTex, rend.whiteTexture),
			pick(material.emissiveTex, rend.whiteTexture),
			{Texture: rend.shadowTex, Sampler: rend.shadowSampler},
		}
		pass.BindFragmentSamplers(bindings[:])

		pass.DrawIndexedPrimitives(item.mesh.indexCount, 1, 0, 0, 0)
	}
	return nil
}

func (canvas *Canvas) buildPasses() []framePass {
	rend := &canvas.rend

	passes := make([]framePass, 0, 3)
	if rend.shadowsOn && len(rend.lights) > 0 && len(rend.draws) > 0 {
		passes = append(passes, framePass{"shadow", (*Canvas).shadowPass})
	}
	passes = append(
		passes,
		framePass{"scene", (*Canvas).scenePass},
		framePass{"tonemap", (*Canvas).tonemapPass},
	)
	return passes
}

func (canvas *Canvas) clearSwapchain() {
	rend := &canvas.rend
	pass := rend.cmd.BeginRenderPass(
		[]sdl.GPUColorTargetInfo{{
			Texture:    rend.swapchain,
			ClearColor: rend.clearColor,
			LoadOp:     sdl.GPU_LOADOP_CLEAR,
			StoreOp:    sdl.GPU_STOREOP_STORE,
		}},
		nil,
	)
	pass.End()
}

func (canvas *Canvas) EndFrame() error {
	rend := &canvas.rend
	if rend.cmd == nil {
		return errors.New("canvas3d: EndFrame called without a successful BeginFrame")
	}

	frameErr := rend.frameErr
	rend.frameErr = nil
	rend.shadowActive = false

	if frameErr == nil {
		canvas.sortDraws()
		for _, p := range canvas.buildPasses() {
			if err := p.run(canvas); err != nil {
				frameErr = fmt.Errorf("%s pass: %w", p.name, err)
				break
			}
		}
	}
	if frameErr != nil {
		canvas.clearSwapchain()
	}

	if err := rend.cmd.Submit(); err != nil && frameErr == nil {
		frameErr = err
	}
	rend.cmd, rend.swapchain = nil, nil
	return frameErr
}
