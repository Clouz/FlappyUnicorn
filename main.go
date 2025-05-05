package main

import (
	"internal/state"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func main() {
	ebiten.SetWindowSize(480, 800)
	ebiten.SetWindowTitle("Flappy Unicorn")
	ebiten.RunGame(state.NewManager())
}
