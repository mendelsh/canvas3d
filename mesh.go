package canvas3d

import (
	"errors"
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
)

type Vertex struct {
	Position [3]float32
	Normal   [3]float32
	UV       [2]float32
}

type MeshDesc struct {
	Vertices []Vertex
	Indices  []uint32
}

type Mesh struct {
	vertexBuffer *sdl.GPUBuffer
	indexBuffer  *sdl.GPUBuffer
	indexCount   uint32
}

func (canvas *Canvas) CreateMesh(desc *MeshDesc) (*Mesh, error) {
	if len(desc.Vertices) == 0 || len(desc.Indices) == 0 {
		return nil, errors.New("canvas3d: mesh needs vertices and indices")
	}

	vertexBytes := sliceBytes(desc.Vertices)
	indexBytes := sliceBytes(desc.Indices)

	vertexBuffer, err := canvas.uploadBuffer(sdl.GPU_BUFFERUSAGE_VERTEX, vertexBytes)
	if err != nil {
		return nil, err
	}
	indexBuffer, err := canvas.uploadBuffer(sdl.GPU_BUFFERUSAGE_INDEX, indexBytes)
	if err != nil {
		canvas.gpu.device.ReleaseBuffer(vertexBuffer)
		return nil, err
	}

	return &Mesh{
		vertexBuffer: vertexBuffer,
		indexBuffer:  indexBuffer,
		indexCount:   uint32(len(desc.Indices)),
	}, nil
}

func (canvas *Canvas) DestroyMesh(mesh *Mesh) {
	if mesh == nil {
		return
	}
	canvas.gpu.device.ReleaseBuffer(mesh.vertexBuffer)
	canvas.gpu.device.ReleaseBuffer(mesh.indexBuffer)
}

func sliceBytes[T any](s []T) []byte {
	if len(s) == 0 {
		return nil
	}
	var zero T
	return unsafe.Slice((*byte)(unsafe.Pointer(&s[0])), len(s)*int(unsafe.Sizeof(zero)))
}

func asBytes[T any](v *T) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(v)), unsafe.Sizeof(*v))
}

func (canvas *Canvas) uploadBuffer(usage sdl.GPUBufferUsageFlags, data []byte) (*sdl.GPUBuffer, error) {
	device := canvas.gpu.device

	buffer, err := device.CreateBuffer(&sdl.GPUBufferCreateInfo{Usage: usage, Size: uint32(len(data))})
	if err != nil {
		return nil, err
	}

	transfer, err := device.CreateTransferBuffer(&sdl.GPUTransferBufferCreateInfo{
		Usage: sdl.GPU_TRANSFERBUFFERUSAGE_UPLOAD,
		Size:  uint32(len(data)),
	})
	if err != nil {
		device.ReleaseBuffer(buffer)
		return nil, err
	}
	defer device.ReleaseTransferBuffer(transfer)

	mapped, err := device.MapTransferBuffer(transfer, false)
	if err != nil {
		device.ReleaseBuffer(buffer)
		return nil, err
	}
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&mapped))
	copy(unsafe.Slice((*byte)(ptr), len(data)), data)
	device.UnmapTransferBuffer(transfer)

	cmd, err := device.AcquireCommandBuffer()
	if err != nil {
		device.ReleaseBuffer(buffer)
		return nil, err
	}
	copyPass := cmd.BeginCopyPass()
	copyPass.UploadToGPUBuffer(
		&sdl.GPUTransferBufferLocation{TransferBuffer: transfer},
		&sdl.GPUBufferRegion{Buffer: buffer, Size: uint32(len(data))},
		false,
	)
	copyPass.End()
	if err := cmd.Submit(); err != nil {
		device.ReleaseBuffer(buffer)
		return nil, err
	}
	return buffer, nil
}
