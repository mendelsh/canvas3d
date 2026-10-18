package canvas3d

import (
	"errors"
	"runtime"
	"sort"
	"sync"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

func init() {
	runtime.LockOSThread()
}

var sdlState struct {
	sync.Mutex
	users  int
	unload func()
}

func acquireSDL() error {
	sdlState.Lock()
	defer sdlState.Unlock()

	if sdlState.users == 0 {
		unload := binsdl.Load().Unload
		if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
			unload()
			return err
		}
		sdlState.unload = unload
	}
	sdlState.users++
	return nil
}

type windowCtx struct {
	window *sdl.Window
}

type gpuCtx struct {
	device  *sdl.GPUDevice
	claimed bool
}

type drawItem struct {
	mesh      *Mesh
	material  *Material
	transform [16]float32
}

type Camera struct {
	Position [3]float32
	Target   [3]float32
	FOV      float32 // radians
	Near     float32
	Far      float32
}

const shadowMapSize = 2048

type renderer struct {
	cmd                              *sdl.GPUCommandBuffer
	swapchain                        *sdl.GPUTexture
	width                            uint32
	height                           uint32
	depth                            *sdl.GPUTexture
	depthW                           uint32
	depthH                           uint32
	colorFmt                         sdl.GPUTextureFormat
	clearColor                       sdl.FColor
	pipelines                        map[pipelineKey]*sdl.GPUGraphicsPipeline
	camera                           Camera
	draws                            []drawItem
	lights                           []Light
	skyZenith, skyHorizon, skyGround [3]float32
	flatNormal                       *Texture
	frameErr                         error
	sampler                          *sdl.GPUSampler
	whiteTexture                     *Texture
	exposure                         float32
	shadowTex                        *sdl.GPUTexture
	shadowSampler                    *sdl.GPUSampler
	shadowsOn                        bool
	shadowActive                     bool
	shadowCenter                     [3]float32
	shadowRadius                     float32
	lightViewProj                    Mat4
}

func releaseSDL() {
	sdlState.Lock()
	defer sdlState.Unlock()

	sdlState.users--
	if sdlState.users == 0 {
		sdl.Quit()
		sdlState.unload()
		sdlState.unload = nil
	}
}

type Canvas struct {
	win  windowCtx
	gpu  gpuCtx
	rend renderer
	live bool
}

func CreateCanvas(title string, width, height int, fullscreen bool) (*Canvas, error) {
	if err := acquireSDL(); err != nil {
		return nil, err
	}

	canvas := &Canvas{live: true}
	canvas.rend.clearColor = sdl.FColor{R: 0.1, G: 0.1, B: 0.12, A: 1.0}

	success := false
	defer func() {
		if !success {
			canvas.Destroy()
		}
	}()

	device, err := sdl.CreateGPUDevice(
		sdl.GPU_SHADERFORMAT_SPIRV|sdl.GPU_SHADERFORMAT_DXIL|sdl.GPU_SHADERFORMAT_MSL,
		true,
		"",
	)
	if err != nil {
		return nil, err
	}
	canvas.gpu.device = device

	flags := sdl.WINDOW_RESIZABLE
	if fullscreen {
		flags |= sdl.WINDOW_FULLSCREEN
	}
	window, err := sdl.CreateWindow(title, width, height, flags)
	if err != nil {
		return nil, err
	}
	canvas.win.window = window

	if err := device.ClaimWindow(window); err != nil {
		return nil, err
	}
	canvas.gpu.claimed = true

	if device.WindowSupportsSwapchainComposition(window, sdl.GPU_SWAPCHAINCOMPOSITION_SDR_LINEAR) {
		err := device.SetSwapchainParameters(window, sdl.GPU_SWAPCHAINCOMPOSITION_SDR_LINEAR, sdl.GPU_PRESENTMODE_VSYNC)
		if err != nil {
			return nil, err
		}
	}
	canvas.rend.colorFmt = device.SwapchainTextureFormat(window)

	canvas.rend.pipelines = make(map[pipelineKey]*sdl.GPUGraphicsPipeline)
	canvas.rend.colorFmt = device.SwapchainTextureFormat(window)

	sampler, err := device.CreateSampler(
		&sdl.GPUSamplerCreateInfo{
			MinFilter:        sdl.GPU_FILTER_LINEAR,
			MagFilter:        sdl.GPU_FILTER_LINEAR,
			MipmapMode:       sdl.GPU_SAMPLERMIPMAPMODE_LINEAR,
			AddressModeU:     sdl.GPU_SAMPLERADDRESSMODE_REPEAT,
			AddressModeV:     sdl.GPU_SAMPLERADDRESSMODE_REPEAT,
			AddressModeW:     sdl.GPU_SAMPLERADDRESSMODE_REPEAT,
			EnableAnisotropy: true,
			MaxAnisotropy:    8,
		},
	)
	if err != nil {
		return nil, err
	}
	canvas.rend.sampler = sampler

	white, err := canvas.CreateTexture(&TextureDesc{Width: 1, Height: 1, Pixels: []byte{255, 255, 255, 255}})
	if err != nil {
		return nil, err
	}
	canvas.rend.whiteTexture = white

	shadowTex, err := device.CreateTexture(
		&sdl.GPUTextureCreateInfo{
			Type:              sdl.GPU_TEXTURETYPE_2D,
			Format:            depthFormat,
			Usage:             sdl.GPU_TEXTUREUSAGE_DEPTH_STENCIL_TARGET | sdl.GPU_TEXTUREUSAGE_SAMPLER,
			Width:             shadowMapSize,
			Height:            shadowMapSize,
			LayerCountOrDepth: 1,
			NumLevels:         1,
		},
	)
	if err != nil {
		return nil, err
	}
	canvas.rend.shadowTex = shadowTex

	shadowSampler, err := device.CreateSampler(
		&sdl.GPUSamplerCreateInfo{
			MinFilter:    sdl.GPU_FILTER_NEAREST,
			MagFilter:    sdl.GPU_FILTER_NEAREST,
			MipmapMode:   sdl.GPU_SAMPLERMIPMAPMODE_NEAREST,
			AddressModeU: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
			AddressModeV: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
			AddressModeW: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
		},
	)
	if err != nil {
		return nil, err
	}
	canvas.rend.shadowSampler = shadowSampler

	canvas.rend.shadowsOn = true
	canvas.rend.shadowRadius = 6

	canvas.rend.camera = Camera{
		Position: [3]float32{0, 0, 3},
		Target:   [3]float32{0, 0, 0},
		FOV:      1.0,
		Near:     0.1,
		Far:      100,
	}

	canvas.SetSky(
		[3]float32{0.25, 0.40, 0.70},
		[3]float32{0.65, 0.70, 0.75},
		[3]float32{0.15, 0.13, 0.12},
	)

	flat, err := canvas.CreateTexture(&TextureDesc{Width: 1, Height: 1, Pixels: []byte{128, 128, 255, 255}})
	if err != nil {
		return nil, err
	}
	canvas.rend.flatNormal = flat

	canvas.rend.exposure = 1.0

	success = true
	return canvas, nil
}

func (canvas *Canvas) Destroy() {
	if canvas == nil || !canvas.live {
		return
	}
	canvas.live = false

	if canvas.gpu.device != nil {
		canvas.gpu.device.WaitForIdle()
	}
	if canvas.gpu.claimed {
		canvas.gpu.device.ReleaseWindow(canvas.win.window)
	}
	if canvas.win.window != nil {
		canvas.win.window.Destroy()
	}
	for _, pipeline := range canvas.rend.pipelines {
		canvas.gpu.device.ReleaseGraphicsPipeline(pipeline)
	}
	if canvas.rend.depth != nil {
		canvas.gpu.device.ReleaseTexture(canvas.rend.depth)
	}
	if canvas.rend.whiteTexture != nil {
		canvas.DestroyTexture(canvas.rend.whiteTexture)
	}
	if canvas.rend.flatNormal != nil {
		canvas.DestroyTexture(canvas.rend.flatNormal)
	}
	if canvas.rend.sampler != nil {
		canvas.gpu.device.ReleaseSampler(canvas.rend.sampler)
	}
	if canvas.rend.shadowTex != nil {
		canvas.gpu.device.ReleaseTexture(canvas.rend.shadowTex)
	}
	if canvas.rend.shadowSampler != nil {
		canvas.gpu.device.ReleaseSampler(canvas.rend.shadowSampler)
	}
	if canvas.gpu.device != nil {
		canvas.gpu.device.Destroy()
	}
	releaseSDL()
}

func (canvas *Canvas) Driver() string {
	return canvas.gpu.device.Driver()
}

func (canvas *Canvas) SetClearColor(r, g, b float32) {
	canvas.rend.clearColor = sdl.FColor{R: r, G: g, B: b, A: 1.0}
}

func (canvas *Canvas) ensureDepth() error {
	rend := &canvas.rend
	if rend.depth != nil && rend.depthW == rend.width && rend.depthH == rend.height {
		return nil
	}
	if rend.depth != nil {
		canvas.gpu.device.ReleaseTexture(rend.depth)
		rend.depth = nil
	}
	depth, err := canvas.gpu.device.CreateTexture(
		&sdl.GPUTextureCreateInfo{
			Type:              sdl.GPU_TEXTURETYPE_2D,
			Format:            depthFormat,
			Usage:             sdl.GPU_TEXTUREUSAGE_DEPTH_STENCIL_TARGET,
			Width:             rend.width,
			Height:            rend.height,
			LayerCountOrDepth: 1,
			NumLevels:         1,
		},
	)
	if err != nil {
		return err
	}
	rend.depth, rend.depthW, rend.depthH = depth, rend.width, rend.height
	return nil
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

		key := pipelineKey{colorFormat: rend.colorFmt, doubleSided: material.doubleSided, blend: material.blend}
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

func (canvas *Canvas) sortDraws() {
	rend := &canvas.rend
	cam := rend.camera.Position

	material := func(item drawItem) *Material {
		if item.material == nil {
			return &defaultMaterial
		}
		return item.material
	}

	dist := func(t Mat4) float32 {
		dx, dy, dz := t[12]-cam[0], t[13]-cam[1], t[14]-cam[2]
		return dx*dx + dy*dy + dz*dz
	}

	sort.SliceStable(rend.draws, func(i, j int) bool {
		a, b := rend.draws[i], rend.draws[j]
		ma, mb := material(a), material(b)
		if ma.blend != mb.blend {
			return !ma.blend
		}
		if ma.blend {
			return dist(a.transform) > dist(b.transform)
		}
		if ma.doubleSided != mb.doubleSided {
			return !ma.doubleSided
		}
		return ma.id < mb.id
	})
}

// BeginFrame returns false when there is nothing to draw to (e.g. the
// window is minimized). In that case skip the frame and do not call EndFrame.
func (canvas *Canvas) BeginFrame() bool {
	cmd, err := canvas.gpu.device.AcquireCommandBuffer()
	if err != nil {
		return false
	}

	swapchain, err := cmd.WaitAndAcquireGPUSwapchainTexture(canvas.win.window)
	if err != nil || swapchain == nil {
		cmd.Submit()
		return false
	}

	canvas.rend.cmd = cmd
	canvas.rend.swapchain = swapchain.Texture
	canvas.rend.width = swapchain.Width
	canvas.rend.height = swapchain.Height
	canvas.rend.draws = canvas.rend.draws[:0]
	canvas.rend.lights = canvas.rend.lights[:0]

	if err := canvas.ensureDepth(); err != nil {
		canvas.rend.frameErr = err // reported by EndFrame; frame goes out without draws
	}
	return true
}

func (canvas *Canvas) EndFrame() error {
	rend := &canvas.rend
	if rend.cmd == nil {
		return errors.New("canvas3d: EndFrame called without a successful BeginFrame")
	}

	frameErr := rend.frameErr
	rend.frameErr = nil

	colorTarget := sdl.GPUColorTargetInfo{
		Texture:    rend.swapchain,
		ClearColor: rend.clearColor,
		LoadOp:     sdl.GPU_LOADOP_CLEAR,
		StoreOp:    sdl.GPU_STOREOP_STORE,
	}

	var depthTarget *sdl.GPUDepthStencilTargetInfo
	if rend.depth != nil {
		depthTarget = &sdl.GPUDepthStencilTargetInfo{
			Texture:        rend.depth,
			ClearDepth:     1.0,
			LoadOp:         sdl.GPU_LOADOP_CLEAR,
			StoreOp:        sdl.GPU_STOREOP_DONT_CARE,
			StencilLoadOp:  sdl.GPU_LOADOP_DONT_CARE,
			StencilStoreOp: sdl.GPU_STOREOP_DONT_CARE,
		}
	}

	canvas.sortDraws()

	rend.shadowActive = false
	if depthTarget != nil && len(rend.draws) > 0 && rend.shadowsOn && len(rend.lights) > 0 {
		if err := canvas.recordShadowPass(); err != nil {
			if frameErr == nil {
				frameErr = err
			}
		} else {
			rend.shadowActive = true
		}
	}

	pass := rend.cmd.BeginRenderPass([]sdl.GPUColorTargetInfo{colorTarget}, depthTarget)
	if depthTarget != nil && len(rend.draws) > 0 {
		if err := canvas.recordDraws(pass); err != nil && frameErr == nil {
			frameErr = err
		}
	}
	pass.End()

	if err := rend.cmd.Submit(); err != nil && frameErr == nil {
		frameErr = err
	}
	rend.cmd, rend.swapchain = nil, nil
	return frameErr
}

func PollEvent(event *Event) bool {
	return sdl.PollEvent(event)
}

func (canvas *Canvas) DrawMesh(mesh *Mesh, material *Material, transform [16]float32) {
	canvas.rend.draws = append(canvas.rend.draws, drawItem{mesh, material, transform})
}

func (canvas *Canvas) SetCamera(camera *Camera) {
	canvas.rend.camera = *camera
}
