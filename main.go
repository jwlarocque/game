package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/jwlarocque/engine"
	"github.com/jwlarocque/engine/mechanism"
	"github.com/jwlarocque/engine/r2"
	"github.com/jwlarocque/engine/tiled"
)

const mapPath = "resources/maps/test_map.tmx"

var player *Player
var levelMap *tiled.Map
var mapImageOpts *ebiten.DrawImageOptions

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
	screen.DrawImage(levelMap.Image, mapImageOpts)
	screen.DrawImage(player.GetImage(), player.GetRenderOpts())

	return nil
}

func main() {
	frameDuration, err := time.ParseDuration("0.08s")
	if err != nil {
		log.Fatal(err)
	}

	//tileset = tiled.NewTilesetFromFile("resources/tiles/cavesofgallet_tiles.tsx")
	levelMap = tiled.NewMapFromJSON("resources/maps/test_map.json")
	mapImageOpts = &ebiten.DrawImageOptions{}
	mapImageOpts.GeoM.Scale(2.0, 2.0)

	player = &Player{
		true,
		time.Now(),
		0,
		map[int]engine.ImageProvider{
			0: engine.NewAnimationFromFolder("resources/entities/player/Idle", frameDuration),
			1: engine.NewAnimationFromFolder("resources/entities/player/Walk", frameDuration)},
		engine.Orientation{},
		mechanism.Collider{
			Position: r2.Vector{X: 64.0, Y: 64.0},
			Velocity: r2.Vector{X: 0.0, Y: 0.0}}}

	if err := ebiten.Run(update, 400, 240, 2, "Hello, World!"); err != nil {
		log.Fatal(err)
	}
}
