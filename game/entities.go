package game

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	playerScore = 0
	cpuScore    = 0
)

var (
	Yellow   = rl.NewColor(243, 213, 91, 255)
	Red      = rl.NewColor(255, 50, 50, 255)
	DarkRed  = rl.NewColor(200, 30, 30, 255)
	LightRed = rl.NewColor(255, 80, 80, 255)
)

type Ball struct {
	X, Y           float32
	SpeedX, SpeedY float32
	Radius         float32
}

func (b *Ball) Draw() {
	rl.DrawCircle(int32(b.X), int32(b.Y), b.Radius, Yellow)
}

func (b *Ball) Update(onAnyScore func()) {
	b.X += b.SpeedX
	b.Y += b.SpeedY

	// Wall Collision (Top/Bottom)
	if b.Y+b.Radius >= float32(rl.GetScreenHeight()) || b.Y-b.Radius <= 0 {
		b.SpeedY *= -1
	}

	// CPU Wins
	if b.X+b.Radius >= float32(rl.GetScreenWidth()) {
		cpuScore++
		b.Reset()
		onAnyScore()
	}

	// Player Wins
	if b.X-b.Radius <= 0 {
		playerScore++
		b.Reset()
		onAnyScore()
	}
}

func (b *Ball) Reset() {
	b.X = float32(rl.GetScreenWidth()) / 2
	b.Y = float32(rl.GetScreenHeight()) / 2

	speedChoices := []float32{-1, 1}
	b.SpeedX *= speedChoices[rl.GetRandomValue(0, 1)]
	b.SpeedY *= speedChoices[rl.GetRandomValue(0, 1)]
}

type Paddle struct {
	X, Y          float32
	Width, Height float32
	Speed         float32
}

func (p *Paddle) Draw() {
	rec := rl.Rectangle{X: p.X, Y: p.Y, Width: p.Width, Height: p.Height}
	rl.DrawRectangleRounded(rec, 0.8, 0, rl.White)
}

func (p *Paddle) LimitMovement() {
	if p.Y <= 0 {
		p.Y = 0
	}
	if p.Y+p.Height >= float32(rl.GetScreenHeight()) {
		p.Y = float32(rl.GetScreenHeight()) - p.Height
	}
}

func (p *Paddle) Update() {
	if rl.IsKeyDown(rl.KeyUp) {
		p.Y -= p.Speed
	}
	if rl.IsKeyDown(rl.KeyDown) {
		p.Y += p.Speed
	}
	p.LimitMovement()
}

type CpuPaddle struct {
	Paddle
}

func (c *CpuPaddle) Update(ballY float32) {
	if c.Y+c.Height/2 > ballY {
		c.Y -= c.Speed
	}
	if c.Y+c.Height/2 <= ballY {
		c.Y += c.Speed
	}
	c.LimitMovement()
}
