package canvas3d

import "math"

// appendQuad adds one quad facing n. u × v must equal n so the winding is
// counter-clockwise when seen from outside.
func appendQuad(desc *MeshDesc, center, n, u, v [3]float32) {
	base := uint32(len(desc.Vertices))
	corners := [4][2]float32{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}}
	for _, c := range corners {
		desc.Vertices = append(desc.Vertices, Vertex{
			Position: [3]float32{
				center[0] + u[0]*c[0] + v[0]*c[1],
				center[1] + u[1]*c[0] + v[1]*c[1],
				center[2] + u[2]*c[0] + v[2]*c[1],
			},
			Normal: n,
			UV:     [2]float32{(c[0] + 1) / 2, (c[1] + 1) / 2},
		})
	}
	desc.Indices = append(desc.Indices, base, base+1, base+2, base, base+2, base+3)
}

// CubeDesc returns a cube centered at the origin with the given edge length.
func CubeDesc(size float32) *MeshDesc {
	h := size / 2
	desc := &MeshDesc{}
	x := [3]float32{h, 0, 0}
	y := [3]float32{0, h, 0}
	z := [3]float32{0, 0, h}
	neg := func(a [3]float32) [3]float32 { return [3]float32{-a[0], -a[1], -a[2]} }

	appendQuad(desc, x, [3]float32{1, 0, 0}, y, z)
	appendQuad(desc, neg(x), [3]float32{-1, 0, 0}, z, y)
	appendQuad(desc, y, [3]float32{0, 1, 0}, z, x)
	appendQuad(desc, neg(y), [3]float32{0, -1, 0}, x, z)
	appendQuad(desc, z, [3]float32{0, 0, 1}, x, y)
	appendQuad(desc, neg(z), [3]float32{0, 0, -1}, y, x)
	return desc
}

// PlaneDesc returns a square in the XZ plane facing +Y, centered at the origin.
func PlaneDesc(size float32) *MeshDesc {
	h := size / 2
	desc := &MeshDesc{}
	appendQuad(desc, [3]float32{}, [3]float32{0, 1, 0}, [3]float32{0, 0, h}, [3]float32{h, 0, 0})
	return desc
}

// SphereDesc returns a UV sphere centered at the origin.
func SphereDesc(radius float32, segments, rings int) *MeshDesc {
	desc := &MeshDesc{}
	for ring := 0; ring <= rings; ring++ {
		theta := math.Pi * float64(ring) / float64(rings)
		for seg := 0; seg <= segments; seg++ {
			phi := 2 * math.Pi * float64(seg) / float64(segments)
			n := [3]float32{
				float32(math.Sin(theta) * math.Cos(phi)),
				float32(math.Cos(theta)),
				float32(math.Sin(theta) * math.Sin(phi)),
			}
			desc.Vertices = append(desc.Vertices, Vertex{
				Position: [3]float32{n[0] * radius, n[1] * radius, n[2] * radius},
				Normal:   n,
				UV:       [2]float32{float32(seg) / float32(segments), float32(ring) / float32(rings)},
			})
		}
	}
	stride := uint32(segments + 1)
	for ring := range rings {
		for seg := range segments {
			a := uint32(ring)*stride + uint32(seg)
			b := a + stride
			c := b + 1
			d := a + 1
			desc.Indices = append(desc.Indices, a, c, b, a, d, c)
		}
	}
	return desc
}
