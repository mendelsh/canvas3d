package canvas3d

import (
	"embed"
	"errors"

	"github.com/Zyko0/go-sdl3/sdl"
)

//go:embed shaders
var shaderFS embed.FS

type shaderFormat struct {
	flag sdl.GPUShaderFormat
	ext  string
}

var shaderFormats = []shaderFormat{
	{sdl.GPU_SHADERFORMAT_DXIL, ".dxil"},
	{sdl.GPU_SHADERFORMAT_SPIRV, ".spv"},
	{sdl.GPU_SHADERFORMAT_MSL, ".msl"},
}

func (canvas *Canvas) loadShader(name string, stage sdl.GPUShaderStage, uniformBuffers, samplers uint32) (*sdl.GPUShader, error) {
	supported := canvas.gpu.device.ShaderFormats()

	for _, format := range shaderFormats {
		if supported&format.flag == 0 {
			continue
		}
		code, err := shaderFS.ReadFile("shaders/" + name + format.ext)
		if err != nil {
			return nil, err
		}
		entry := "main"
		if format.flag == sdl.GPU_SHADERFORMAT_MSL {
			entry = "main0" // SPIR-V-cross style output renames main
		}
		return canvas.gpu.device.CreateGPUShader(
			&sdl.GPUShaderCreateInfo{
				Code:              code,
				Entrypoint:        entry,
				Format:            format.flag,
				Stage:             stage,
				NumUniformBuffers: uniformBuffers,
				NumSamplers:       samplers,
			},
		)
	}
	return nil, errors.New("canvas3d: no supported shader format for this device")
}
