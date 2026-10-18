package canvas3d

import (
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
)

const vs_name = "lit.vert"
const fs_name = "lit.frag"

const svs_name = "shadow.vert"
const sfs_name = "shadow.frag"

const depthFormat = sdl.GPU_TEXTUREFORMAT_D32_FLOAT

type pipelineKey struct {
	colorFormat sdl.GPUTextureFormat
	shadow      bool
	doubleSided bool
	blend       bool
}

func (canvas *Canvas) getPipeline(key pipelineKey) (*sdl.GPUGraphicsPipeline, error) {
	if pipeline, ok := canvas.rend.pipelines[key]; ok {
		return pipeline, nil
	}

	vsName, fsName := vs_name, fs_name
	var fsUniforms, fsSamplers uint32 = 2, 6
	cull := sdl.GPU_CULLMODE_BACK
	targets := sdl.GPUGraphicsPipelineTargetInfo{
		ColorTargetDescriptions: []sdl.GPUColorTargetDescription{{Format: key.colorFormat}},
		HasDepthStencilTarget:   true,
		DepthStencilFormat:      depthFormat,
	}
	if key.doubleSided {
		cull = sdl.GPU_CULLMODE_NONE
	}
	if key.blend {
		targets.ColorTargetDescriptions = []sdl.GPUColorTargetDescription{{
			Format: key.colorFormat,
			BlendState: sdl.GPUColorTargetBlendState{
				EnableBlend:         true,
				SrcColorBlendfactor: sdl.GPU_BLENDFACTOR_SRC_ALPHA,
				DstColorBlendfactor: sdl.GPU_BLENDFACTOR_ONE_MINUS_SRC_ALPHA,
				ColorBlendOp:        sdl.GPU_BLENDOP_ADD,
				SrcAlphaBlendfactor: sdl.GPU_BLENDFACTOR_ONE,
				DstAlphaBlendfactor: sdl.GPU_BLENDFACTOR_ONE_MINUS_SRC_ALPHA,
				AlphaBlendOp:        sdl.GPU_BLENDOP_ADD,
			},
		}}
	}
	if key.shadow {
		vsName, fsName = svs_name, sfs_name
		fsUniforms, fsSamplers = 0, 0
		cull = sdl.GPU_CULLMODE_NONE // keeps thin and one-sided geometry casting shadows
		targets.ColorTargetDescriptions = nil
	}

	vs, err := canvas.loadShader(vsName, sdl.GPU_SHADERSTAGE_VERTEX, 1, 0)
	if err != nil {
		return nil, err
	}
	defer canvas.gpu.device.ReleaseShader(vs)

	fs, err := canvas.loadShader(fsName, sdl.GPU_SHADERSTAGE_FRAGMENT, fsUniforms, fsSamplers)
	if err != nil {
		return nil, err
	}
	defer canvas.gpu.device.ReleaseShader(fs)

	var v Vertex
	pipeline, err := canvas.gpu.device.CreateGraphicsPipeline(
		&sdl.GPUGraphicsPipelineCreateInfo{
			VertexShader:   vs,
			FragmentShader: fs,
			PrimitiveType:  sdl.GPU_PRIMITIVETYPE_TRIANGLELIST,
			VertexInputState: sdl.GPUVertexInputState{
				VertexBufferDescriptions: []sdl.GPUVertexBufferDescription{{
					Slot:      0,
					Pitch:     uint32(unsafe.Sizeof(v)),
					InputRate: sdl.GPU_VERTEXINPUTRATE_VERTEX,
				}},
				VertexAttributes: []sdl.GPUVertexAttribute{
					{Location: 0, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT3, Offset: uint32(unsafe.Offsetof(v.Position))},
					{Location: 1, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT3, Offset: uint32(unsafe.Offsetof(v.Normal))},
					{Location: 2, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT2, Offset: uint32(unsafe.Offsetof(v.UV))},
				},
			},
			RasterizerState: sdl.GPURasterizerState{
				FillMode:  sdl.GPU_FILLMODE_FILL,
				CullMode:  cull,
				FrontFace: sdl.GPU_FRONTFACE_COUNTER_CLOCKWISE,
			},
			DepthStencilState: sdl.GPUDepthStencilState{
				EnableDepthTest:  true,
				EnableDepthWrite: !key.blend,
				CompareOp:        sdl.GPU_COMPAREOP_LESS,
			},
			TargetInfo: targets,
		},
	)
	if err != nil {
		return nil, err
	}

	canvas.rend.pipelines[key] = pipeline
	return pipeline, nil
}
