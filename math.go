package canvas3d

import "math"

// Mat4 is column-major.
type Mat4 = [16]float32

func Identity() Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

func Translation(x, y, z float32) Mat4 {
	m := Identity()
	m[12], m[13], m[14] = x, y, z
	return m
}

func RotationY(angle float32) Mat4 {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))
	m := Identity()
	m[0], m[8] = c, s
	m[2], m[10] = -s, c
	return m
}

// Mul returns a * b (b is applied first).
func Mul(a, b Mat4) Mat4 {
	var out Mat4
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			var sum float32
			for k := 0; k < 4; k++ {
				sum += a[k*4+row] * b[col*4+k]
			}
			out[col*4+row] = sum
		}
	}
	return out
}

func perspective(fov, aspect, near, far float32) Mat4 {
	f := float32(1.0 / math.Tan(float64(fov)*0.5))
	var m Mat4
	m[0] = f / aspect
	m[5] = f
	m[10] = far / (near - far)
	m[11] = -1
	m[14] = near * far / (near - far)
	return m
}

func lookAt(eye, target [3]float32) Mat4 {
	return lookAtUp(eye, target, [3]float32{0, 1, 0})
}

func lookAtUp(eye, target, up [3]float32) Mat4 {
	f := normalize3(sub3(target, eye))
	s := normalize3(cross3(f, up))
	u := cross3(s, f)

	m := Identity()
	m[0], m[4], m[8] = s[0], s[1], s[2]
	m[1], m[5], m[9] = u[0], u[1], u[2]
	m[2], m[6], m[10] = -f[0], -f[1], -f[2]
	m[12] = -dot3(s, eye)
	m[13] = -dot3(u, eye)
	m[14] = dot3(f, eye)
	return m
}

// orthographic builds a symmetric box projection, depth 0..1.
func orthographic(half, near, far float32) Mat4 {
	var m Mat4
	m[0] = 1 / half
	m[5] = 1 / half
	m[10] = 1 / (near - far)
	m[14] = near / (near - far)
	m[15] = 1
	return m
}

func sub3(a, b [3]float32) [3]float32 {
	return [3]float32{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func dot3(a, b [3]float32) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func cross3(a, b [3]float32) [3]float32 {
	return [3]float32{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

func normalize3(v [3]float32) [3]float32 {
	l := float32(math.Sqrt(float64(dot3(v, v))))
	if l == 0 {
		return v
	}
	return [3]float32{v[0] / l, v[1] / l, v[2] / l}
}

func Scale(x, y, z float32) Mat4 {
	m := Identity()
	m[0], m[5], m[10] = x, y, z
	return m
}
