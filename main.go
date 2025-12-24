package main

import (
	"fmt"
	"gago/game"
)

func main() {
	fmt.Println("Starting the game")

	screenWidth := int32(1280)
	screenHeight := int32(800)

	g := game.NewGame(screenWidth, screenHeight)
	g.Run()
}
