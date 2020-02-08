// TODO: un-spaghettify this

package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/jwlarocque/engine"
	"github.com/jwlarocque/engine/mechanism"
	"github.com/jwlarocque/engine/r2"
)

const (
	drag          = 0.8
	playerAccel   = 1.0
	velocityLimit = 1.0 // TODO: rename this, only applies to input
	stickTheshold = 0.2
	gravity       = 0.3

	keyJump      = ebiten.KeyW
	keyMoveLeft  = ebiten.KeyA
	keyMoveRight = ebiten.KeyD
)

type Player struct {
	jumping        bool
	jumpTime       time.Time
	CurrentState   int
	ImageProviders map[int]engine.ImageProvider
	engine.Orientation
	mechanism.Collider
}

// TODO: write a player constructor to fill in defaults

// TODO: Go "enums" suck.  Find a better way.
type state int

const (
	idle state = iota
	walk
)

// TODO: decouple from framerate
func (p *Player) Update() {
	p.Velocity = p.Velocity.Scale(drag)
	if p.Velocity.Magnitude() < stickTheshold {
		p.Velocity.X = 0
		p.Velocity.Y = 0
	}
	p.Velocity = p.getInputVelocity()
	p.Velocity.Y += gravity
	p.Position = p.Position.Add(p.Velocity)

	if math.Abs(p.Velocity.X) > 0 {
		p.CurrentState = int(walk)
	} else {
		p.CurrentState = int(idle)
	}

	// TODO: temp - limit Y to screen bottom
	if p.Position.Y > 180.0 {
		p.Position.Y = 180.0
		p.Velocity.Y = 0.0
		if p.jumping && !ebiten.IsKeyPressed(keyJump) {
			p.jumping = false
		}
	}
}

func (p Player) GetImage() *ebiten.Image {
	return p.ImageProviders[p.CurrentState].GetImage()
}

// TODO: this really should be in the engine package?
// GetRenderOpts provides the transforms on the entity's image (position, orientation, etc.)
func (p Player) GetRenderOpts() *ebiten.DrawImageOptions {
	opts := ebiten.DrawImageOptions{}
	if p.HorizFlip {
		opts.GeoM.Scale(-1, 1)
		imageWidth, _ := p.GetImage().Size()
		opts.GeoM.Translate(float64(imageWidth), 0)
	}
	opts.GeoM.Translate(p.Position.X, p.Position.Y)
	return &opts
}

// TODO: move this to its own file
func (p *Player) getInputVelocity() r2.Vector {
	newVelocity := r2.Vector{p.Velocity.X, p.Velocity.Y}

	inputAccelX := 0.0
	if ebiten.IsKeyPressed(keyMoveLeft) {
		inputAccelX -= playerAccel
	}
	if ebiten.IsKeyPressed(keyMoveRight) {
		inputAccelX += playerAccel
	}

	// set player orientation
	if inputAccelX < 0 {
		p.HorizFlip = true
	} else if inputAccelX > 0 {
		p.HorizFlip = false
	}

	// add inputAccelX to velocity
	if math.Abs(newVelocity.X+inputAccelX) <= velocityLimit {
		newVelocity.X += inputAccelX
	} else if newVelocity.X+inputAccelX > velocityLimit {
		newVelocity.X = velocityLimit
	} else if newVelocity.X+inputAccelX < -velocityLimit {
		newVelocity.X = -velocityLimit
	}

	// TODO: better jumping (jumping in games is a strange and arcane thing)
	if ebiten.IsKeyPressed(keyJump) {
		if !p.jumping {
			p.jumpTime = time.Now()
			p.jumping = true
			newVelocity.Y -= 8.0
		}
	}

	return newVelocity
}

func (p Player) IsJumping() bool {
	return p.jumping
}
