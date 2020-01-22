package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/lafriks/go-tiled"

	"github.com/jwlarocque/engine"
)

const mapPath = "resources/levels/test_map.tmx"

var player *Player

func update(screen *ebiten.Image) error {
	// TODO: state update here
	player.Update()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	ebitenutil.DebugPrint(screen, "Hello, World!")
	screen.DrawImage(player.GetImage(), player.GetRenderOpts())

	return nil
}

func main() {
	frameDuration, err := time.ParseDuration("0.08s")
	if err != nil {
		log.Fatal(err)
	}

	gameMap, err := tiled.LoadFromFile(mapPath)

	if err != nil {
		log.Fatal("Error parsing map")
	}

	log.Print(gameMap)

	player = &Player{
		engine.Orientation{},
		engine.Entity{
			Position:     engine.Vector2{X: 64.0, Y: 64.0},
			Velocity:     engine.Vector2{X: 0.0, Y: 0.0},
			CurrentState: 0,
			ImageProviders: map[int]engine.ImageProvider{
				0: engine.NewAnimationFromFolder("resources/entities/player/Idle", frameDuration),
				1: engine.NewAnimationFromFolder("resources/entities/player/Walk", frameDuration)}}}

	if err := ebiten.Run(update, 400, 240, 2, "Hello, World!"); err != nil {
		log.Fatal(err)
	}
}
