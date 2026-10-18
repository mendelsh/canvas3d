package main

import (
	"fmt"

	"github.com/mendelsh/canvas3d"
)

func main() {
	canvas, err := canvas3d.CreateCanvas("hello", 800, 600, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer canvas.Destroy()

	cube, err := canvas.CreateMesh(canvas3d.CubeDesc(1))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer canvas.DestroyMesh(cube)

	red := canvas.CreateMaterial(&canvas3d.MaterialDesc{
		BaseColor: [4]float32{0.9, 0.2, 0.2, 1},
	})
	defer canvas.DestroyMaterial(red)

	angle := float32(0)
	running := true
	for running {
		var event canvas3d.Event
		for canvas3d.PollEvent(&event) {
			if event.Type == canvas3d.EVENT_QUIT {
				running = false
			}
		}

		if !canvas.BeginFrame() {
			continue
		}
		angle += 0.01
		canvas.AddLight(&canvas3d.Light{Direction: [3]float32{-0.4, -1, -0.6}})
		canvas.DrawMesh(cube, red, canvas3d.RotationY(angle))
		if err := canvas.EndFrame(); err != nil {
			fmt.Println(err)
			return
		}
	}
}
