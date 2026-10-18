package main

import (
	"fmt"
	"os"

	"github.com/mendelsh/canvas3d"
	"github.com/mendelsh/canvas3d/model"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func run() error {
	canvas, err := canvas3d.CreateCanvas("canvas3d viewer", 1280, 720, false)
	if err != nil {
		return fmt.Errorf("canvas: %w", err)
	}
	defer canvas.Destroy()
	fmt.Println("driver:", canvas.Driver())

	canvas.SetClearColor(0.1, 0.2, 0.4)
	canvas.SetAmbient([3]float32{0.12, 0.12, 0.15})
	canvas.SetExposure(0.8)
	canvas.SetShadowArea([3]float32{0, 0, 0}, 6)

	// Optional model: assets/DamagedHelmet.glb, or pass a path as the first argument.
	modelPath := "assets/DamagedHelmet.glb"
	if len(os.Args) > 1 {
		modelPath = os.Args[1]
	}
	helmet, err := model.Load(canvas, modelPath)
	if err != nil {
		fmt.Println("no model loaded (put a .glb in assets/ or pass a path):", err)
		helmet = nil
	} else {
		defer helmet.Destroy(canvas)
	}

	// Meshes
	floor, err := canvas.CreateMesh(canvas3d.PlaneDesc(10))
	if err != nil {
		return fmt.Errorf("floor: %w", err)
	}
	defer canvas.DestroyMesh(floor)

	sphere, err := canvas.CreateMesh(canvas3d.SphereDesc(0.6, 48, 24))
	if err != nil {
		return fmt.Errorf("sphere: %w", err)
	}
	defer canvas.DestroyMesh(sphere)

	// Textures
	const checkerSize = 1024
	checker, err := canvas.CreateTexture(&canvas3d.TextureDesc{
		Width:   checkerSize,
		Height:  checkerSize,
		Pixels:  checkerboard(checkerSize, checkerSize/8),
		Mipmaps: true,
		SRGB:    true,
	})
	if err != nil {
		return fmt.Errorf("checkerboard: %w", err)
	}
	defer canvas.DestroyTexture(checker)

	// Materials (destroyed before the texture they use, since defers run in reverse)
	floorMat := canvas.CreateMaterial(&canvas3d.MaterialDesc{
		BaseColor:        [4]float32{0.5, 0.9, 0.5, 1},
		BaseColorTexture: checker,
		Roughness:        0.8,
	})
	chrome := canvas.CreateMaterial(&canvas3d.MaterialDesc{
		BaseColor: [4]float32{0.95, 0.93, 0.88, 1},
		Metallic:  1,
		Roughness: 0.2,
	})
	glass := canvas.CreateMaterial(&canvas3d.MaterialDesc{
		BaseColor: [4]float32{0.4, 0.6, 1, 0.35},
		Roughness: 0.1,
		AlphaMode: canvas3d.AlphaBlend,
	})
	defer canvas.DestroyMaterial(floorMat)
	defer canvas.DestroyMaterial(chrome)
	defer canvas.DestroyMaterial(glass)

	orbit := canvas3d.NewOrbitCamera([3]float32{0, 1, 0}, 7)
	rotating := true
	angle := float32(0)

	running := true
	for running {
		canvas3d.UpdateTime()
		dt := canvas3d.DeltaTime()

		var event canvas3d.Event
		for canvas3d.PollEvent(&event) {
			orbit.HandleEvent(&event)

			if event.Type == canvas3d.EVENT_QUIT {
				running = false
			}
			if event.Type == canvas3d.EVENT_KEY_DOWN {
				switch event.KeyboardEvent().Key {
				case canvas3d.K_ESCAPE:
					running = false
				case canvas3d.K_SPACE:
					rotating = !rotating
				case canvas3d.K_F:
					fmt.Printf("FPS: %.0f\n", 1/dt)
				}
			}
		}

		if !canvas.BeginFrame() {
			continue
		}

		cam := orbit.Camera()
		canvas.SetCamera(&cam)

		canvas.AddLight(&canvas3d.Light{ // main light, casts the shadows
			Direction: [3]float32{-0.4, -1, -0.6},
			Color:     [3]float32{1, 0.95, 0.9},
			Intensity: 1,
		})
		canvas.AddLight(&canvas3d.Light{ // weak cool fill from the other side
			Direction: [3]float32{0.7, -0.3, 0.5},
			Color:     [3]float32{0.4, 0.5, 0.9},
			Intensity: 0.35,
		})

		if rotating {
			angle += dt
		}

		// Everything sits on or above the floor at y = 0.
		canvas.DrawMesh(floor, floorMat, canvas3d.Identity())

		if helmet != nil {
			helmet.Draw(canvas, canvas3d.Mul(canvas3d.Translation(0, 1.2, 0), canvas3d.RotationY(angle)))
		}

		// radius 0.6, so its center is 0.6 above the floor
		canvas.DrawMesh(sphere, chrome, canvas3d.Translation(2.2, 0.6, 0))

		// scaled 1.6x, radius 0.96
		canvas.DrawMesh(sphere, glass, canvas3d.Mul(canvas3d.Translation(-2.2, 0.96, 0), canvas3d.Scale(1.6, 1.6, 1.6)))

		if err := canvas.EndFrame(); err != nil {
			return fmt.Errorf("frame: %w", err)
		}
	}
	return nil
}

func checkerboard(size, cell int) []byte {
	pixels := make([]byte, size*size*4)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			value := byte(255)
			if (x/cell+y/cell)%2 == 0 {
				value = 90
			}
			i := (y*size + x) * 4
			pixels[i], pixels[i+1], pixels[i+2], pixels[i+3] = value, value, value, 255
		}
	}
	return pixels
}
