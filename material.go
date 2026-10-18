package canvas3d

type AlphaMode int

const (
	AlphaOpaque AlphaMode = iota // alpha is ignored
	AlphaMask                    // pixels with alpha below AlphaCutoff are discarded
	AlphaBlend                   // blended, drawn after opaque objects, far to near
)

// / Zero values mean "default" for: BaseColor (white), Roughness (1),
// / NormalScale (1), OcclusionStrength (1). Metallic zero is a plain non-metal.
type MaterialDesc struct {
	BaseColor                [4]float32
	BaseColorTexture         *Texture // sRGB
	Metallic                 float32
	Roughness                float32
	MetallicRoughnessTexture *Texture // linear: G = roughness, B = metallic
	NormalTexture            *Texture // linear
	NormalScale              float32
	OcclusionTexture         *Texture // linear: R
	OcclusionStrength        float32
	EmissiveFactor           [3]float32
	EmissiveTexture          *Texture // sRGB
	DoubleSided              bool
	AlphaMode                AlphaMode
	AlphaCutoff              float32 // AlphaMask only; 0 means 0.5
}

// Layout must match the Material cbuffer in the fragment shader (64 bytes).
type materialUniforms struct {
	baseColor [4]float32
	emissive  [4]float32
	params    [4]float32 // metallic, roughness, normal scale, occlusion strength
	alpha     [4]float32
}

type Material struct {
	id            uint32
	uniforms      materialUniforms
	baseColorTex  *Texture
	metalRoughTex *Texture
	normalTex     *Texture
	occlusionTex  *Texture
	emissiveTex   *Texture
	doubleSided   bool
	blend         bool
}

var materialCount uint32

var defaultMaterial = Material{
	uniforms: materialUniforms{
		baseColor: [4]float32{1, 1, 1, 1},
		params:    [4]float32{0, 1, 1, 1},
		alpha:     [4]float32{0, 0.5},
	},
}

func orDefault(value, fallback float32) float32 {
	if value == 0 {
		return fallback
	}
	return value
}

func (canvas *Canvas) CreateMaterial(desc *MaterialDesc) *Material {
	materialCount++
	cutoff := orDefault(desc.AlphaCutoff, 0.5)

	color := desc.BaseColor
	if color == [4]float32{} {
		color = [4]float32{1, 1, 1, 1}
	}
	return &Material{
		uniforms: materialUniforms{
			baseColor: color,
			emissive:  [4]float32{desc.EmissiveFactor[0], desc.EmissiveFactor[1], desc.EmissiveFactor[2], 0},
			params: [4]float32{
				desc.Metallic,
				orDefault(desc.Roughness, 1),
				orDefault(desc.NormalScale, 1),
				orDefault(desc.OcclusionStrength, 1),
			},
			alpha: [4]float32{float32(desc.AlphaMode), cutoff, 0, 0},
		},
		id:            materialCount,
		baseColorTex:  desc.BaseColorTexture,
		metalRoughTex: desc.MetallicRoughnessTexture,
		normalTex:     desc.NormalTexture,
		occlusionTex:  desc.OcclusionTexture,
		emissiveTex:   desc.EmissiveTexture,
		doubleSided:   desc.DoubleSided,
		blend:         desc.AlphaMode == AlphaBlend,
	}
}

// No GPU resources yet to destroy.
func (canvas *Canvas) DestroyMaterial(material *Material) {}
