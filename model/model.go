package model

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/url"
	"os"
	"path/filepath"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"

	"github.com/mendelsh/canvas3d"
)

type primitive struct {
	mesh     *canvas3d.Mesh
	material *canvas3d.Material // nil means the canvas default
}

type node struct {
	prims []primitive
	world canvas3d.Mat4
}

type Model struct {
	nodes     []node
	meshes    []*canvas3d.Mesh
	materials []*canvas3d.Material
	textures  []*canvas3d.Texture
}

// Draw submits every primitive of the model. transform is applied on top of
// the node transforms stored in the file.
func (m *Model) Draw(canvas *canvas3d.Canvas, transform canvas3d.Mat4) {
	for _, n := range m.nodes {
		world := canvas3d.Mul(transform, n.world)
		for _, p := range n.prims {
			canvas.DrawMesh(p.mesh, p.material, world)
		}
	}
}

// Destroy releases everything the model created. Call it before canvas.Destroy.
func (m *Model) Destroy(canvas *canvas3d.Canvas) {
	for _, mesh := range m.meshes {
		canvas.DestroyMesh(mesh)
	}
	for _, material := range m.materials {
		canvas.DestroyMaterial(material)
	}
	for _, texture := range m.textures {
		canvas.DestroyTexture(texture)
	}
	*m = Model{}
}

type textureKey struct {
	index int
	srgb  bool
}

type loader struct {
	canvas    *canvas3d.Canvas
	doc       *gltf.Document
	dir       string
	model     *Model
	meshes    map[int][]primitive
	materials map[int]*canvas3d.Material
	textures  map[textureKey]*canvas3d.Texture
}

func Load(canvas *canvas3d.Canvas, path string) (*Model, error) {
	doc, err := gltf.Open(path)
	if err != nil {
		return nil, err
	}
	if len(doc.Scenes) == 0 {
		return nil, errors.New("model: file has no scenes")
	}

	l := &loader{
		canvas:    canvas,
		doc:       doc,
		dir:       filepath.Dir(path),
		model:     &Model{},
		meshes:    map[int][]primitive{},
		materials: map[int]*canvas3d.Material{},
		textures:  map[textureKey]*canvas3d.Texture{},
	}

	ok := false
	defer func() {
		if !ok {
			l.model.Destroy(canvas)
		}
	}()

	sceneIndex := 0
	if doc.Scene != nil {
		sceneIndex = *doc.Scene
	}
	for _, root := range doc.Scenes[sceneIndex].Nodes {
		if err := l.addNode(root, canvas3d.Identity()); err != nil {
			return nil, err
		}
	}

	ok = true
	return l.model, nil
}

func (l *loader) addNode(index int, parent canvas3d.Mat4) error {
	n := l.doc.Nodes[index]
	world := canvas3d.Mul(parent, localMatrix(n))

	if n.Mesh != nil {
		prims, err := l.meshPrimitives(*n.Mesh)
		if err != nil {
			return err
		}
		l.model.nodes = append(l.model.nodes, node{prims: prims, world: world})
	}
	for _, child := range n.Children {
		if err := l.addNode(child, world); err != nil {
			return err
		}
	}
	return nil
}

func (l *loader) meshPrimitives(index int) ([]primitive, error) {
	if prims, ok := l.meshes[index]; ok {
		return prims, nil
	}

	var prims []primitive
	for _, src := range l.doc.Meshes[index].Primitives {
		if src.Mode != gltf.PrimitiveTriangles {
			continue // points, lines and strips are not supported
		}
		prim, err := l.buildPrimitive(src)
		if err != nil {
			return nil, fmt.Errorf("model: mesh %d: %w", index, err)
		}
		prims = append(prims, prim)
	}
	l.meshes[index] = prims
	return prims, nil
}

func (l *loader) buildPrimitive(src *gltf.Primitive) (primitive, error) {
	doc := l.doc

	posIndex, ok := src.Attributes[gltf.POSITION]
	if !ok {
		return primitive{}, errors.New("primitive has no POSITION")
	}
	positions, err := modeler.ReadPosition(doc, doc.Accessors[posIndex], nil)
	if err != nil {
		return primitive{}, err
	}

	var indices []uint32
	if src.Indices != nil {
		indices, err = modeler.ReadIndices(doc, doc.Accessors[*src.Indices], nil)
		if err != nil {
			return primitive{}, err
		}
	} else {
		indices = make([]uint32, len(positions))
		for i := range indices {
			indices[i] = uint32(i)
		}
	}

	var normals [][3]float32
	if idx, ok := src.Attributes[gltf.NORMAL]; ok {
		if normals, err = modeler.ReadNormal(doc, doc.Accessors[idx], nil); err != nil {
			return primitive{}, err
		}
	} else {
		normals = computeNormals(positions, indices)
	}

	var uvs [][2]float32
	if idx, ok := src.Attributes[gltf.TEXCOORD_0]; ok {
		if uvs, err = modeler.ReadTextureCoord(doc, doc.Accessors[idx], nil); err != nil {
			return primitive{}, err
		}
	}

	vertices := make([]canvas3d.Vertex, len(positions))
	for i := range vertices {
		vertices[i].Position = positions[i]
		vertices[i].Normal = normals[i]
		if i < len(uvs) {
			vertices[i].UV = uvs[i]
		}
	}

	mesh, err := l.canvas.CreateMesh(&canvas3d.MeshDesc{Vertices: vertices, Indices: indices})
	if err != nil {
		return primitive{}, err
	}
	l.model.meshes = append(l.model.meshes, mesh)

	material, err := l.material(src.Material)
	if err != nil {
		return primitive{}, err
	}
	return primitive{mesh: mesh, material: material}, nil
}

func (l *loader) material(index *int) (*canvas3d.Material, error) {
	if index == nil {
		return nil, nil
	}
	if material, ok := l.materials[*index]; ok {
		return material, nil
	}

	src := l.doc.Materials[*index]
	desc := &canvas3d.MaterialDesc{
		BaseColor: [4]float32{1, 1, 1, 1},
		Metallic:  1, // glTF defaults when there is no pbrMetallicRoughness block
		Roughness: 1,
	}

	if pbr := src.PBRMetallicRoughness; pbr != nil {
		c := pbr.BaseColorFactorOrDefault()
		desc.BaseColor = [4]float32{float32(c[0]), float32(c[1]), float32(c[2]), float32(c[3])}
		desc.Metallic = float32(pbr.MetallicFactorOrDefault())
		desc.Roughness = float32(pbr.RoughnessFactorOrDefault())

		var err error
		if pbr.BaseColorTexture != nil {
			if desc.BaseColorTexture, err = l.texture(int(pbr.BaseColorTexture.Index), true); err != nil {
				return nil, err
			}
		}
		if pbr.MetallicRoughnessTexture != nil {
			if desc.MetallicRoughnessTexture, err = l.texture(int(pbr.MetallicRoughnessTexture.Index), false); err != nil {
				return nil, err
			}
		}
	}

	var err error
	if nt := src.NormalTexture; nt != nil && nt.Index != nil {
		if desc.NormalTexture, err = l.texture(int(*nt.Index), false); err != nil {
			return nil, err
		}
		desc.NormalScale = float32(nt.ScaleOrDefault())
	}
	if ot := src.OcclusionTexture; ot != nil && ot.Index != nil {
		if desc.OcclusionTexture, err = l.texture(int(*ot.Index), false); err != nil {
			return nil, err
		}
		desc.OcclusionStrength = float32(ot.StrengthOrDefault())
	}
	if src.EmissiveTexture != nil {
		if desc.EmissiveTexture, err = l.texture(int(src.EmissiveTexture.Index), true); err != nil {
			return nil, err
		}
	}
	desc.EmissiveFactor = [3]float32{float32(src.EmissiveFactor[0]), float32(src.EmissiveFactor[1]), float32(src.EmissiveFactor[2])}

	desc.DoubleSided = src.DoubleSided
	switch src.AlphaMode {
	case gltf.AlphaMask:
		desc.AlphaMode = canvas3d.AlphaMask
		desc.AlphaCutoff = float32(src.AlphaCutoffOrDefault())
	case gltf.AlphaBlend:
		desc.AlphaMode = canvas3d.AlphaBlend
	}

	material := l.canvas.CreateMaterial(desc)
	l.model.materials = append(l.model.materials, material)
	l.materials[*index] = material
	return material, nil
}

func (l *loader) texture(index int, srgb bool) (*canvas3d.Texture, error) {
	key := textureKey{index, srgb}
	if texture, ok := l.textures[key]; ok {
		return texture, nil
	}

	src := l.doc.Textures[index]
	if src.Source == nil {
		return nil, fmt.Errorf("model: texture %d has no image source", index)
	}
	data, err := l.imageBytes(*src.Source)
	if err != nil {
		return nil, err
	}

	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("model: decoding image %d: %w", *src.Source, err)
	}
	// NRGBA keeps alpha un-premultiplied, which is what the GPU texture expects.
	pixels := image.NewNRGBA(decoded.Bounds())
	draw.Draw(pixels, pixels.Bounds(), decoded, decoded.Bounds().Min, draw.Src)

	texture, err := l.canvas.CreateTexture(&canvas3d.TextureDesc{
		Width:   pixels.Bounds().Dx(),
		Height:  pixels.Bounds().Dy(),
		Pixels:  pixels.Pix,
		Mipmaps: true,
	})
	if err != nil {
		return nil, err
	}
	l.model.textures = append(l.model.textures, texture)
	l.textures[key] = texture
	return texture, nil
}

func (l *loader) imageBytes(index int) ([]byte, error) {
	img := l.doc.Images[index]
	switch {
	case img.BufferView != nil: // typical for .glb
		view := l.doc.BufferViews[*img.BufferView]
		buffer := l.doc.Buffers[view.Buffer]
		return buffer.Data[view.ByteOffset : view.ByteOffset+view.ByteLength], nil
	case img.IsEmbeddedResource(): // data: URI
		return img.MarshalData()
	case img.URI != "":
		name, err := url.PathUnescape(img.URI)
		if err != nil {
			return nil, err
		}
		return os.ReadFile(filepath.Join(l.dir, name))
	}
	return nil, fmt.Errorf("model: image %d has no data", index)
}

// localMatrix returns the node's transform, either its matrix or T*R*S.
func localMatrix(n *gltf.Node) canvas3d.Mat4 {
	if m := n.MatrixOrDefault(); m != gltf.DefaultMatrix {
		var out canvas3d.Mat4
		for i, v := range m {
			out[i] = float32(v)
		}
		return out
	}

	t := n.TranslationOrDefault()
	q := n.RotationOrDefault() // x, y, z, w
	s := n.ScaleOrDefault()
	x, y, z, w := float32(q[0]), float32(q[1]), float32(q[2]), float32(q[3])
	sx, sy, sz := float32(s[0]), float32(s[1]), float32(s[2])

	return canvas3d.Mat4{
		(1 - 2*(y*y+z*z)) * sx, 2 * (x*y + z*w) * sx, 2 * (x*z - y*w) * sx, 0,
		2 * (x*y - z*w) * sy, (1 - 2*(x*x+z*z)) * sy, 2 * (y*z + x*w) * sy, 0,
		2 * (x*z + y*w) * sz, 2 * (y*z - x*w) * sz, (1 - 2*(x*x+y*y)) * sz, 0,
		float32(t[0]), float32(t[1]), float32(t[2]), 1,
	}
}

// builds smooth normals for primitives that ship without any.
func computeNormals(positions [][3]float32, indices []uint32) [][3]float32 {
	normals := make([][3]float32, len(positions))
	for i := 0; i+2 < len(indices); i += 3 {
		a, b, c := positions[indices[i]], positions[indices[i+1]], positions[indices[i+2]]
		e1 := [3]float32{b[0] - a[0], b[1] - a[1], b[2] - a[2]}
		e2 := [3]float32{c[0] - a[0], c[1] - a[1], c[2] - a[2]}
		face := [3]float32{
			e1[1]*e2[2] - e1[2]*e2[1],
			e1[2]*e2[0] - e1[0]*e2[2],
			e1[0]*e2[1] - e1[1]*e2[0],
		}
		for _, idx := range indices[i : i+3] {
			normals[idx][0] += face[0]
			normals[idx][1] += face[1]
			normals[idx][2] += face[2]
		}
	}
	for i, n := range normals {
		l := float32(math.Sqrt(float64(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])))
		if l > 0 {
			normals[i] = [3]float32{n[0] / l, n[1] / l, n[2] / l}
		} else {
			normals[i] = [3]float32{0, 1, 0}
		}
	}
	return normals
}
