package canvas3d

import (
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
	hdr                              *sdl.GPUTexture
	hdrSampler                       *sdl.GPUSampler
	renderW                          uint32
	renderH                          uint32
	renderScale                      float32
	shadowSize                       int
	shadowTexSize                    uint32
	shadowSoftness                   float32
	tonemapper                       Tonemapper
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

	hdrSampler, err := device.CreateSampler(&sdl.GPUSamplerCreateInfo{
		MinFilter:    sdl.GPU_FILTER_LINEAR,
		MagFilter:    sdl.GPU_FILTER_LINEAR,
		MipmapMode:   sdl.GPU_SAMPLERMIPMAPMODE_NEAREST,
		AddressModeU: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
		AddressModeV: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
		AddressModeW: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
	})
	if err != nil {
		return nil, err
	}
	canvas.rend.hdrSampler = hdrSampler

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
	canvas.SetQuality(QualityAuto)
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
	if canvas.rend.hdr != nil {
		canvas.gpu.device.ReleaseTexture(canvas.rend.hdr)
	}
	if canvas.rend.hdrSampler != nil {
		canvas.gpu.device.ReleaseSampler(canvas.rend.hdrSampler)
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
// window is minimized).
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

	if err := canvas.ensureTargets(); err != nil {
		canvas.rend.frameErr = err
	}
	return true
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
