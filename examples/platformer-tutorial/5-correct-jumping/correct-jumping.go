package main

import (
	"image/color"

	"github.com/oakmound/oak/collision"

	"github.com/oakmound/oak/physics"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

// Collision labels
const (
	// The only collision label we need for this demo is 'ground',
	// indicating something we shouldn't be able to fall or walk through
	Ground collision.Label = 1
)

func main() {
	oak.Add("platformer", func(string, interface{}) {

		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)

		render.Draw(char.R)

		char.Speed = physics.NewVector(3, 3)

		fallSpeed := .1

		char.Bind(func(id int, nothing interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)
			// Move left and right with A and D
			if oak.IsDown(key.A) {
				char.ShiftX(-char.Speed.X())
			}
			if oak.IsDown(key.D) {
				char.ShiftX(char.Speed.X())
			}
			oldY := char.Y()
			char.ShiftY(char.Delta.Y())
			hit := collision.HitLabel(char.Space, Ground)

			// If we've moved in y value this frame and in the last frame,
			// we were below what we're trying to hit, we are still falling
			if hit != nil && !(oldY != char.Y() && oldY+char.H > hit.Y()) {
				// Correct our y if we started falling into the ground
				char.SetY(hit.Y() - char.H)
				char.Delta.SetY(0)
				// Jump with Space
				if oak.IsDown(key.Spacebar) {
					char.Delta.ShiftY(-char.Speed.Y())
				}
			} else {
				// Fall if there's no ground
				char.Delta.ShiftY(fallSpeed)
			}
			return 0
		}, event.Enter)

		ground := entities.NewSolid(0, 400, 500, 20,
			render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
			nil, 0)
		ground.UpdateLabel(Ground)

		render.Draw(ground.R)

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})
	oak.Init("platformer")
}
