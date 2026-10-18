package canvas3d

import (
	"errors"
	"math/bits"
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
)

type TextureDesc struct {
	Width, Height int
	Pixels        []byte
	Mipmaps       bool
	SRGB          bool
}

type Texture struct {
	texture       *sdl.GPUTexture
	width, height int
}

func mipLevels(width, height int) uint32 {
	size := uint(max(width, height))
	return uint32(bits.Len(size))
}

func (canvas *Canvas) CreateTexture(desc *TextureDesc) (*Texture, error) {
	if desc.Width <= 0 || desc.Height <= 0 || len(desc.Pixels) != desc.Width*desc.Height*4 {
		return nil, errors.New("canvas3d: texture pixels must be exactly Width*Height*4 bytes (RGBA8)")
	}
	device := canvas.gpu.device

	levels := uint32(1)
	usage := sdl.GPU_TEXTUREUSAGE_SAMPLER
	if desc.Mipmaps {
		levels = mipLevels(desc.Width, desc.Height)
		usage |= sdl.GPU_TEXTUREUSAGE_COLOR_TARGET
	}

	format := sdl.GPU_TEXTUREFORMAT_R8G8B8A8_UNORM
	if desc.SRGB {
		format = sdl.GPU_TEXTUREFORMAT_R8G8B8A8_UNORM_SRGB
	}

	texture, err := device.CreateTexture(
		&sdl.GPUTextureCreateInfo{
			Type:              sdl.GPU_TEXTURETYPE_2D,
			Format:            format,
			Usage:             usage,
			Width:             uint32(desc.Width),
			Height:            uint32(desc.Height),
			LayerCountOrDepth: 1,
			NumLevels:         levels,
		},
	)
	if err != nil {
		return nil, err
	}

	transfer, err := device.CreateTransferBuffer(&sdl.GPUTransferBufferCreateInfo{
		Usage: sdl.GPU_TRANSFERBUFFERUSAGE_UPLOAD,
		Size:  uint32(len(desc.Pixels)),
	})
	if err != nil {
		device.ReleaseTexture(texture)
		return nil, err
	}
	defer device.ReleaseTransferBuffer(transfer)

	mapped, err := device.MapTransferBuffer(transfer, false)
	if err != nil {
		device.ReleaseTexture(texture)
		return nil, err
	}
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&mapped))
	copy(unsafe.Slice((*byte)(ptr), len(desc.Pixels)), desc.Pixels)
	device.UnmapTransferBuffer(transfer)

	cmd, err := device.AcquireCommandBuffer()
	if err != nil {
		device.ReleaseTexture(texture)
		return nil, err
	}
	copyPass := cmd.BeginCopyPass()
	copyPass.UploadToGPUTexture(
		&sdl.GPUTextureTransferInfo{TransferBuffer: transfer},
		&sdl.GPUTextureRegion{Texture: texture, W: uint32(desc.Width), H: uint32(desc.Height), D: 1},
		false,
	)
	copyPass.End()
	if desc.Mipmaps {
		cmd.GenerateMipmapsForGPUTexture(texture)
	}
	if err := cmd.Submit(); err != nil {
		device.ReleaseTexture(texture)
		return nil, err
	}

	return &Texture{texture: texture, width: desc.Width, height: desc.Height}, nil
}

func (canvas *Canvas) DestroyTexture(texture *Texture) {
	if texture == nil {
		return
	}
	canvas.gpu.device.ReleaseTexture(texture.texture)
}
