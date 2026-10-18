package canvas3d

import "math"

type OrbitCamera struct {
	Target   [3]float32
	Distance float32
	Yaw      float32 // radians, around the Y axis
	Pitch    float32 // radians, positive looks down on the target
	FOV      float32 // radians
	Near     float32
	Far      float32

	rotating bool
	panning  bool
}

func NewOrbitCamera(target [3]float32, distance float32) *OrbitCamera {
	return &OrbitCamera{
		Target:   target,
		Distance: distance,
		Pitch:    0.3,
		FOV:      1.0,
		Near:     0.1,
		Far:      100,
	}
}

// HandleEvent feeds a polled event into the camera. Call it for every event.
func (o *OrbitCamera) HandleEvent(event *Event) {
	switch event.Type {
	case EVENT_MOUSE_BUTTON_DOWN, EVENT_MOUSE_BUTTON_UP:
		button := event.MouseButtonEvent()
		down := event.Type == EVENT_MOUSE_BUTTON_DOWN
		switch button.Button {
		case 1: // left
			o.rotating = down
		case 3: // right
			o.panning = down
		}

	case EVENT_MOUSE_MOTION:
		motion := event.MouseMotionEvent()
		dx, dy := float64(motion.Xrel), float64(motion.Yrel)

		if o.rotating {
			o.Yaw -= float32(dx * 0.005)
			o.Pitch += float32(dy * 0.005)
			const limit = 1.5 // just under a quarter turn, avoids flipping over the top
			o.Pitch = float32(math.Max(-limit, math.Min(limit, float64(o.Pitch))))
		}
		if o.panning {
			sy, cy := float32(math.Sin(float64(o.Yaw))), float32(math.Cos(float64(o.Yaw)))
			sp, cp := float32(math.Sin(float64(o.Pitch))), float32(math.Cos(float64(o.Pitch)))
			right := [3]float32{cy, 0, -sy}
			up := [3]float32{-sy * sp, cp, -cy * sp}
			scale := o.Distance * 0.0015
			for i := 0; i < 3; i++ {
				o.Target[i] += (-right[i]*float32(dx) + up[i]*float32(dy)) * scale
			}
		}

	case EVENT_MOUSE_WHEEL:
		wheel := event.MouseWheelEvent()
		o.Distance *= float32(math.Pow(0.9, float64(wheel.Y)))
		o.Distance = float32(math.Max(0.2, float64(o.Distance)))
	}
}

// Camera returns the current view, ready for canvas.SetCamera.
func (o *OrbitCamera) Camera() Camera {
	sy, cy := float32(math.Sin(float64(o.Yaw))), float32(math.Cos(float64(o.Yaw)))
	sp, cp := float32(math.Sin(float64(o.Pitch))), float32(math.Cos(float64(o.Pitch)))
	return Camera{
		Position: [3]float32{
			o.Target[0] + o.Distance*cp*sy,
			o.Target[1] + o.Distance*sp,
			o.Target[2] + o.Distance*cp*cy,
		},
		Target: o.Target,
		FOV:    o.FOV,
		Near:   o.Near,
		Far:    o.Far,
	}
}
