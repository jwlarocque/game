// TODO: un-spaghettify this

package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/jwlarocque/engine"
)

const (
	drag          = 0.85
	playerAccel   = 1.0
	velocityLimit = 1.0 // TODO: rename this, only applies to input
	stickTheshold = 0.2
	gravity       = 0.35

	keyJump      = ebiten.KeyW
	keyMoveLeft  = ebiten.KeyA
	keyMoveRight = ebiten.KeyD
)

type Player struct {
	jumping     bool
	jumpTime    time.Time
	orientation engine.Orientation
	engine.Entity
}

// TODO: write a player constructor to fill in defaults

// TODO: Go "enums" suck.  Find a better way.
type state int

const (
	idle state = iota
	walk
)

func (p *Player) Update() {
	p.Velocity = p.Velocity.Scale(drag)
	if p.Velocity.Magnitude() < stickTheshold {
		p.Velocity.X = 0
		p.Velocity.Y = 0
	}
	p.Velocity = p.getInputVelocity()
	log.Print(p.Velocity)
	log.Print("")
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

// TODO: this really should be in the engine package?
// GetRenderOpts provides the transforms on the entity's image (position, orientation, etc.)
func (p Player) GetRenderOpts() *ebiten.DrawImageOptions {
	opts := ebiten.DrawImageOptions{}
	if p.orientation.Horizontal {
		opts.GeoM.Scale(-1, 1)
		imageWidth, _ := p.GetImage().Size()
		opts.GeoM.Translate(float64(imageWidth), 0)
	}
	opts.GeoM.Translate(p.Position.X, p.Position.Y)
	return &opts
}

// TODO: move this to its own file
func (p *Player) getInputVelocity() engine.Vector2 {
	newVelocity := engine.Vector2{p.Velocity.X, p.Velocity.Y}

	inputAccelX := 0.0
	if ebiten.IsKeyPressed(keyMoveLeft) {
		inputAccelX -= playerAccel
	}
	if ebiten.IsKeyPressed(keyMoveRight) {
		inputAccelX += playerAccel
	}

	// set player orientation
	if inputAccelX < 0 {
		p.orientation.Horizontal = true
	} else if inputAccelX > 0 {
		p.orientation.Horizontal = false
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
			newVelocity.Y -= 5.0
		}
	}

	log.Print(fmt.Sprintf("X: %.6f", inputAccelX))
	log.Print(newVelocity)
	return newVelocity
}

func (p Player) IsJumping() bool {
	return p.jumping
}
