package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/jwlarocque/engine"
	"github.com/jwlarocque/engine/tiled"
)

const mapPath = "resources/maps/test_map.tmx"

var player *Player
var tileset *tiled.Tileset

func update(screen *ebiten.Image) error {
	// TODO: state update here
	player.Update()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	isJumping := "not jumping"
	if player.IsJumping() {
		isJumping = "jumping!"
	}
	ebitenutil.DebugPrint(screen, isJumping)
	screen.DrawImage(player.GetImage(), player.GetRenderOpts())
	screen.DrawImage(tileset.TilesImage, &ebiten.DrawImageOptions{})

	return nil
}

func main() {
	frameDuration, err := time.ParseDuration("0.08s")
	if err != nil {
		log.Fatal(err)
	}

	tileset = tiled.NewTilesetFromFile("resources/tiles/cavesofgallet_tiles.tsx")

	player = &Player{
		false,
		time.Now(),
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
