package canvas3d

import "github.com/Zyko0/go-sdl3/sdl"

const maxLights = 4

type Light struct {
	Direction [3]float32
	Color     [3]float32
	Intensity float32
}

// Layout must match the Frame cbuffer in the fragment shader (288 bytes).
type frameUniforms struct {
	cameraPos     [4]float32
	skyZenith     [4]float32
	skyHorizon    [4]float32
	skyGround     [4]float32
	lightCount    [4]float32
	direction     [maxLights][4]float32
	color         [maxLights][4]float32
	lightViewProj Mat4
	shadowParams  [4]float32
}

type drawUniforms struct {
	mvp   Mat4
	model Mat4
}

// sets the gradient the environment lighting and reflections come from.
func (canvas *Canvas) SetSky(zenith, horizon, ground [3]float32) {
	canvas.rend.skyZenith = zenith
	canvas.rend.skyHorizon = horizon
	canvas.rend.skyGround = ground
}

// SetAmbient is a uniform sky: flat ambient light, no visible gradient.
func (canvas *Canvas) SetAmbient(color [3]float32) {
	canvas.SetSky(color, color, color)
}

func (canvas *Canvas) AddLight(light *Light) {
	if len(canvas.rend.lights) >= maxLights {
		return
	}
	canvas.rend.lights = append(canvas.rend.lights, *light)
}

func (canvas *Canvas) buildFrameUniforms() frameUniforms {
	rend := &canvas.rend
	pad := func(v [3]float32) [4]float32 { return [4]float32{v[0], v[1], v[2], 0} }

	out := frameUniforms{
		cameraPos:  [4]float32{rend.camera.Position[0], rend.camera.Position[1], rend.camera.Position[2], rend.exposure},
		skyZenith:  pad(rend.skyZenith),
		skyHorizon: pad(rend.skyHorizon),
		skyGround:  pad(rend.skyGround),
		lightCount: [4]float32{float32(len(rend.lights))},
	}

	for i, light := range rend.lights {
		dir := normalize3(light.Direction)
		out.direction[i] = [4]float32{dir[0], dir[1], dir[2], 0}

		color := light.Color
		if color == [3]float32{} {
			color = [3]float32{1, 1, 1}
		}
		intensity := light.Intensity
		if intensity == 0 {
			intensity = 1
		}
		out.color[i] = [4]float32{color[0] * intensity, color[1] * intensity, color[2] * intensity, 0}
	}

	if rend.shadowActive {
		texelWorld := 2 * rend.shadowRadius / shadowMapSize
		out.lightViewProj = rend.lightViewProj
		out.shadowParams = [4]float32{1, 1.0 / shadowMapSize, 0.0005, texelWorld * 1.5}
	}

	return out
}

// scales scene brightness before tone mapping. 1 is neutral.
func (canvas *Canvas) SetExposure(exposure float32) {
	canvas.rend.exposure = exposure
}

// SetShadows turns the first light's shadow on or off.
func (canvas *Canvas) SetShadows(enabled bool) {
	canvas.rend.shadowsOn = enabled
}

// SetShadowArea sets the region that gets shadows: a box of half-size radius
// around center. A smaller area gives sharper shadows.
func (canvas *Canvas) SetShadowArea(center [3]float32, radius float32) {
	canvas.rend.shadowCenter = center
	canvas.rend.shadowRadius = radius
}

func (canvas *Canvas) updateShadowMatrix() {
	rend := &canvas.rend
	dir := normalize3(rend.lights[0].Direction)
	r := rend.shadowRadius
	c := rend.shadowCenter

	eye := [3]float32{c[0] - dir[0]*r*2, c[1] - dir[1]*r*2, c[2] - dir[2]*r*2}
	up := [3]float32{0, 1, 0}
	if dir[1] > 0.99 || dir[1] < -0.99 {
		up = [3]float32{0, 0, 1}
	}
	view := lookAtUp(eye, c, up)
	rend.lightViewProj = Mul(orthographic(r, 0.1, r*4), view)
}

func (canvas *Canvas) recordShadowPass() error {
	rend := &canvas.rend
	canvas.updateShadowMatrix()

	pipeline, err := canvas.getPipeline(pipelineKey{shadow: true})
	if err != nil {
		return err
	}

	pass := rend.cmd.BeginRenderPass(nil, &sdl.GPUDepthStencilTargetInfo{
		Texture:        rend.shadowTex,
		ClearDepth:     1.0,
		LoadOp:         sdl.GPU_LOADOP_CLEAR,
		StoreOp:        sdl.GPU_STOREOP_STORE,
		StencilLoadOp:  sdl.GPU_LOADOP_DONT_CARE,
		StencilStoreOp: sdl.GPU_STOREOP_DONT_CARE,
	})
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
	return nil
}
