package canvas3d

// Quality is a bundle of rendering options. Most programs never need to touch
// anything else: the default is QualityAuto and everything just works.
type Quality int

const (
	// QualityAuto picks for you. For now it is the same as QualityMedium;
	// later it will adapt to the machine while the program runs.
	QualityAuto Quality = iota
	QualityLow
	QualityMedium
	QualityHigh
)

// Tonemapper is the curve that squeezes bright HDR values into the range the
// screen can show.
type Tonemapper int

const (
	TonemapACES    Tonemapper = iota // punchy filmic curve (default)
	TonemapNeutral                   // keeps base colors closer to the textures
	TonemapClamp                     // no curve, just clips; handy for comparing
)

// SetQuality applies a preset. Call it first, then override single options.
func (canvas *Canvas) SetQuality(quality Quality) {
	rend := &canvas.rend
	switch quality {
	case QualityLow:
		rend.shadowSize, rend.shadowSoftness, rend.renderScale = 1024, 1.5, 0.75
	case QualityHigh:
		rend.shadowSize, rend.shadowSoftness, rend.renderScale = 4096, 3.0, 1.0
	default: // QualityMedium and QualityAuto
		rend.shadowSize, rend.shadowSoftness, rend.renderScale = 2048, 2.0, 1.0
	}
}

// SetRenderScale renders the scene at a fraction of the window size (or more
// than 1 for supersampling) and scales it to fit. Lower values are faster.
func (canvas *Canvas) SetRenderScale(scale float32) {
	canvas.rend.renderScale = min(max(scale, 0.25), 2.0)
}

// SetShadowResolution sets the shadow map size in pixels (square).
func (canvas *Canvas) SetShadowResolution(size int) {
	canvas.rend.shadowSize = min(max(size, 256), 8192)
}

// SetShadowSoftness sets how wide the shadow edge blur is, in shadow map texels.
func (canvas *Canvas) SetShadowSoftness(texels float32) {
	canvas.rend.shadowSoftness = min(max(texels, 0), 8)
}

func (canvas *Canvas) SetTonemapper(tonemapper Tonemapper) {
	canvas.rend.tonemapper = tonemapper
}
