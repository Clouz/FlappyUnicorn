package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 800

	gravity         = 0.22
	jumpSpeed       = -6.5
	gentleFallLimit = 6

	obstacleSpeed   = 2.2
	obstacleGap     = 240
	obstacleSpacing = 260
	obstacleWidth   = 70

	playerX         = 80
	initialLives    = 3
	invulnerability = 45

	startCloudDrift  = 0.12
	gentleFloorGuard = 30
)

type GameState int

const (
	stateStart GameState = iota
	statePlaying
	stateGameOver
)

type Game struct {
	bg      *ebiten.Image
	unicorn *ebiten.Image

	playerY   float64
	playerVel float64

	obstacles []Obstacle
	score     int
	lives     int
	state     GameState
	invuln    int
	rnd       *rand.Rand

	backgroundOffset float64
}

type Obstacle struct {
	x      float64
	gapY   float64
	passed bool
	tint   color.NRGBA
}

func NewGame() *Game {
	bg, _, err := ebitenutil.NewImageFromFile("assets/gfx/background.png")
	if err != nil {
		log.Fatal(err)
	}
	uni, _, err := ebitenutil.NewImageFromFile("assets/gfx/unicorn_jump.png")
	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		bg:      bg,
		unicorn: uni,
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	g.reset(true)
	return g
}

func (g *Game) randomGap() float64 {
	min := 160.0
	max := float64(screenHeight) - 160.0
	return min + g.rnd.Float64()*(max-min)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.ErrFinished
	}

	switch g.state {
	case stateStart:
		if g.flapJustPressed() {
			g.reset(false)
			g.state = statePlaying
		}
	case statePlaying:
		g.backgroundOffset += startCloudDrift
		g.applyInput()
		g.applyPhysics()
		g.updateObstacles()
	case stateGameOver:
		if g.flapJustPressed() {
			g.reset(false)
			g.state = statePlaying
		}
	}

	return nil
}

func (g *Game) flapJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		return true
	}
	return false
}

func (g *Game) applyInput() {
	if g.flapJustPressed() {
		g.playerVel = jumpSpeed
	}
}

func (g *Game) applyPhysics() {
	if g.invuln > 0 {
		g.invuln--
	}

	g.playerVel += gravity
	if g.playerVel > gentleFallLimit {
		g.playerVel = gentleFallLimit
	}
	g.playerY += g.playerVel

	if g.playerY < 0 {
		g.playerY = 0
		g.playerVel = 0
	}

	floor := float64(screenHeight-gentleFloorGuard) - float64(g.unicorn.Bounds().Dy())
	if g.playerY > floor {
		g.playerY = floor
		g.handleCollision()
	}
}

func (g *Game) updateObstacles() {
	playerRect := image.Rect(int(playerX), int(g.playerY), int(playerX)+g.unicorn.Bounds().Dx(), int(g.playerY)+g.unicorn.Bounds().Dy())
	for i := range g.obstacles {
		o := &g.obstacles[i]
		o.x -= obstacleSpeed

		if !o.passed && o.x+float64(obstacleWidth) < playerX {
			o.passed = true
			g.score++
		}

		if o.x+float64(obstacleWidth) < 0 {
			o.x += float64(len(g.obstacles)) * obstacleSpacing
			o.gapY = g.randomGap()
			o.passed = false
			o.tint = pastelColor(g.rnd)
		}

		if g.collisionWith(*o, playerRect) {
			g.handleCollision()
		}
	}
}

func (g *Game) collisionWith(o Obstacle, playerRect image.Rectangle) bool {
	topRect := image.Rect(int(o.x), 0, int(o.x)+obstacleWidth, int(o.gapY)-obstacleGap/2)
	bottomRect := image.Rect(int(o.x), int(o.gapY)+obstacleGap/2, int(o.x)+obstacleWidth, screenHeight)
	return playerRect.Overlaps(topRect) || playerRect.Overlaps(bottomRect)
}

func (g *Game) handleCollision() {
	if g.invuln > 0 {
		return
	}
	g.playerVel = jumpSpeed / 2
	g.invuln = invulnerability
	g.lives--
	if g.lives <= 0 {
		g.state = stateGameOver
		g.invuln = 0
	}
}

func (g *Game) reset(initial bool) {
	g.playerY = screenHeight/2 - float64(g.unicorn.Bounds().Dy())/2
	g.playerVel = 0
	g.score = 0
	g.lives = initialLives
	g.invuln = 0
	if initial {
		g.state = stateStart
	}

	g.backgroundOffset = 0
	g.obstacles = g.obstacles[:0]
	for i := 0; i < 3; i++ {
		g.obstacles = append(g.obstacles, Obstacle{
			x:    float64(screenWidth + i*obstacleSpacing),
			gapY: g.randomGap(),
			tint: pastelColor(g.rnd),
		})
	}
}

func pastelColor(rnd *rand.Rand) color.NRGBA {
	palette := []color.NRGBA{
		{R: 0xFF, G: 0xC1, B: 0xCC, A: 0xFF},
		{R: 0xFF, G: 0xF4, B: 0xB0, A: 0xFF},
		{R: 0xB5, G: 0xE0, B: 0xF2, A: 0xFF},
		{R: 0xD8, G: 0xC4, B: 0xF2, A: 0xFF},
		{R: 0xC4, G: 0xED, B: 0xC7, A: 0xFF},
	}
	return palette[rnd.Intn(len(palette))]
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)

	pOp := &ebiten.DrawImageOptions{}
	if g.invuln > 0 && (g.invuln/5)%2 == 0 {
		pOp.ColorM.Scale(1, 1, 1, 0.6)
	}
	pOp.GeoM.Translate(playerX, g.playerY)
	screen.DrawImage(g.unicorn, pOp)

	for _, o := range g.obstacles {
		vector.DrawFilledRect(screen, float32(o.x), 0, float32(obstacleWidth), float32(o.gapY)-obstacleGap/2, o.tint)
		vector.DrawFilledRect(screen, float32(o.x), float32(o.gapY)+obstacleGap/2, float32(obstacleWidth), float32(screenHeight)-(float32(o.gapY)+obstacleGap/2), o.tint)
	}

	g.drawHUD(screen)
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	scaleX := float64(screenWidth) / float64(g.bg.Bounds().Dx())
	scaleY := float64(screenHeight) / float64(g.bg.Bounds().Dy())
	drift := math.Mod(g.backgroundOffset, float64(screenWidth))

	for i := -1; i <= 1; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(i*screenWidth)-drift, 0)
		screen.DrawImage(g.bg, op)
	}
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	switch g.state {
	case stateStart:
		ebitenutil.DebugPrintAt(screen, "Ciao amica dell'unicorno!", 40, 260)
		ebitenutil.DebugPrintAt(screen, "Tocca o premi SPAZIO per iniziare.", 40, 300)
		ebitenutil.DebugPrintAt(screen, "Accompagna Luna l'unicorno tra le nuvole!", 40, 340)
	case statePlaying:
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Stelle raccolte: %d", g.score), 20, 20)
	case stateGameOver:
		ebitenutil.DebugPrintAt(screen, "Bravissima! Vuoi riprovare?", 70, 300)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Stelle raccolte: %d", g.score), 120, 340)
		ebitenutil.DebugPrintAt(screen, "Premi SPAZIO o tocca per continuare", 60, 380)
	}

	hearts := ""
	for i := 0; i < g.lives; i++ {
		hearts += "❤ "
	}
	if hearts != "" {
		ebitenutil.DebugPrintAt(screen, hearts, 20, 60)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Flappy Unicorn")
	if err := ebiten.RunGame(NewGame()); err != nil && err != ebiten.ErrFinished {
		log.Fatal(err)
	}
}
