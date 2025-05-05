package main

import "github.com/hajimehoshi/ebiten/v2"

type Manager struct {
	current State
	states  map[StateID]State
}

func NewManager() *Manager {
	m := &Manager{states: make(map[StateID]State)}
	m.states[Ready] = NewReady(m)
	m.states[Play] = game.NewPlay(m)
	m.states[GameOver] = NewGameOver(m)
	m.current = m.states[Ready]
	return m
}
func (m *Manager) Update() error              { return m.current.Update() }
func (m *Manager) Draw(s *ebiten.Image)       { m.current.Draw(s) }
func (m *Manager) Layout(w, h int) (int, int) { return 480, 800 }
func (m *Manager) Switch(to StateID)          { m.current = m.states[to] }
