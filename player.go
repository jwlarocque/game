package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/jwlarocque/engine"
)

const drag = 0.8
const accel = 1.0
const velocityLimit = 1.0
const stickLimit = 0.2

type Player struct {
	orientation engine.Orientation
	engine.Entity
}

// TODO: Go "enums" suck.  Find a better way.
type state int

const (
	idle state = iota
	walk
)

func (p *Player) Update() {
	movementVector := engine.GetMovementVector()

	p.Velocity = getNewVelocity(p.Velocity, movementVector)
	p.Position = p.Position.Add(p.Velocity)

	if movementVector.X > 0 {
		p.orientation.Horizontal = false
	} else if movementVector.X < 0 {
		p.orientation.Horizontal = true
	}

	if math.Abs(p.Velocity.X) > 0 {
		p.CurrentState = int(walk)
	} else {
		p.CurrentState = int(idle)
	}

	// TODO: temp - limit Y to screen bottom
	if p.Position.Y > 180.0 {
		p.Position.Y = 180.0
		p.Velocity.Y = 0.0
	}
}

// TODO: this really should be in the engine package
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

func getNewVelocity(oldVelocity engine.Vector2, movementVector engine.Vector2) engine.Vector2 {
	newVelocity := oldVelocity.Scale(drag)
	newVelocity = newVelocity.Add(movementVector.Scale(accel))
	if newVelocity.Magnitude() > velocityLimit {
		newVelocity = newVelocity.Normalize().Scale(velocityLimit)
	} else if newVelocity.Magnitude() < stickLimit {
		newVelocity.X = 0
		newVelocity.Y = 0
	}
	return newVelocity
}
