package main

type Player struct {
	sprite anim.Sprite
	y, vel float64
}

func NewPlayer() *Player {
	spr := anim.New("unicorn_jump_2x2.png", 2, 2, 6)
	return &Player{sprite: spr, y: 400}
}

func (p *Player) Update(inputJump bool) {
	if inputJump {
		p.vel = -4
	}
	p.vel += 0.25
	p.y += p.vel
	p.sprite.Update()
}
