package game

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	MenuState GameState = iota
	PlayState
	CountdownState
	GameOverState
)

type Game struct {
	ScreenWidth   int32
	ScreenHeight  int32
	State         GameState
	Ball          *Ball
	Player        *Paddle
	CPU           *CpuPaddle
	Winner        string
	Countdown     int
	LastScoreTime time.Time
}

func NewGame(width, height int32) *Game {
	ball := &Ball{
		X:      float32(width) / 2,
		Y:      float32(height) / 2,
		SpeedX: 7,
		SpeedY: 7,
		Radius: 20,
	}

	player := &Paddle{
		Width:  25,
		Height: 120,
		Speed:  6,
	}
	player.X = float32(width) - player.Width - 10
	player.Y = float32(height)/2 - player.Height/2

	cpu := &CpuPaddle{
		Paddle: Paddle{
			Width:  25,
			Height: 120,
			Speed:  6,
			X:      10,
		},
	}
	cpu.Y = float32(height)/2 - cpu.Height/2

	return &Game{
		ScreenWidth:  width,
		ScreenHeight: height,
		State:        MenuState,
		Ball:         ball,
		Player:       player,
		CPU:          cpu,
		Countdown:    1,
	}
}

func (g *Game) Run() {
	rl.InitWindow(g.ScreenWidth, g.ScreenHeight, "My Pong Game!")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		g.Update()
		g.Draw()
	}
}

func (g *Game) Update() {
	switch g.State {
	case MenuState:
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.State = PlayState
		}
	case PlayState:
		g.Ball.Update(
			func() {
				g.State = CountdownState
				g.LastScoreTime = time.Now()
			},
		)
		g.Player.Update()
		g.CPU.Update(g.Ball.Y)

		ballPos := rl.Vector2{X: g.Ball.X, Y: g.Ball.Y}
		playerRec := rl.Rectangle{X: g.Player.X, Y: g.Player.Y, Width: g.Player.Width, Height: g.Player.Height}
		if rl.CheckCollisionCircleRec(ballPos, g.Ball.Radius, playerRec) {
			g.Ball.SpeedX *= -1
		}

		cpuRec := rl.Rectangle{X: g.CPU.X, Y: g.CPU.Y, Width: g.CPU.Width, Height: g.CPU.Height}
		if rl.CheckCollisionCircleRec(ballPos, g.Ball.Radius, cpuRec) {
			g.Ball.SpeedX *= -1
		}

		if playerScore >= 10 {
			g.Winner = "Player"
			g.State = GameOverState
		}
		if cpuScore >= 10 {
			g.Winner = "CPU"
			g.State = GameOverState
		}

	case CountdownState:
		if time.Since(g.LastScoreTime).Seconds() >= 1 {
			g.Countdown--
			g.LastScoreTime = time.Now()
		}
		if g.Countdown < 0 {
			g.State = PlayState
			g.Countdown = 1
		}

	case GameOverState:
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.Reset()
			g.State = MenuState
		}
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	g.DrawCourt()
	if g.State == PlayState || g.State == CountdownState {
		g.Ball.Draw()
		g.CPU.Draw()
		g.Player.Draw()
	}
	switch g.State {
	case MenuState:
		g.DrawTextCenter("Press ENTER to start", 60, rl.White)
	case PlayState:
		g.Ball.Draw()
		g.CPU.Draw()
		g.Player.Draw()
		g.DrawScores()
	case CountdownState:
		g.DrawScores()
		g.DrawTextCenter(fmt.Sprintf("%d", g.Countdown), 120, rl.White)
	case GameOverState:
		g.DrawTextCenter(fmt.Sprintf("%s wins!", g.Winner), 80, rl.White)
		g.DrawTextCenter("Press ENTER to restart", 40, rl.LightGray, 100)
	}
}

func (g *Game) DrawCourt() {
	rl.ClearBackground(DarkRed)
	rl.DrawRectangle(g.ScreenWidth/2, 0, g.ScreenWidth/2, g.ScreenHeight, Red)
	rl.DrawCircle(g.ScreenWidth/2, g.ScreenHeight/2, 150, LightRed)
	rl.DrawLine(g.ScreenWidth/2, 0, g.ScreenWidth/2, g.ScreenHeight, rl.White)
}

func (g *Game) DrawScores() {
	rl.DrawText(fmt.Sprint(cpuScore), g.ScreenWidth/4-20, 20, 80, rl.White)
	rl.DrawText(fmt.Sprint(playerScore), 3*g.ScreenWidth/4-20, 20, 80, rl.White)
}

func (g *Game) DrawTextCenter(text string, fontSize int32, color rl.Color, yOffset ...int32) {
	textSize := rl.MeasureText(text, fontSize)
	x := (g.ScreenWidth - textSize) / 2
	y := g.ScreenHeight/2 - fontSize/2
	if len(yOffset) > 0 {
		y += yOffset[0]
	}
	rl.DrawText(text, x, y, fontSize, color)
}

func (g *Game) Reset() {
	g.Ball.Reset()
	playerScore = 0
	cpuScore = 0
	g.Winner = ""
}
