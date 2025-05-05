package main

import "github.com/hajimehoshi/ebiten/v2"

type Obstacle struct {
	x float64
}

func (o *Obstacle) Update() {
	o.x -= 3
	if o.x < -120 {
		o.x += 600
	} // ricicla
}

func (o *Obstacle) Draw(dst *ebiten.Image) {
	// usa sprite-sheet: prima metà arcobaleno, seconda fiore
}
