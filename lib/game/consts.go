package game

import "math"

const (
	FieldWidth  = 2000
	FieldHeight = 1500

	Radius = 30

	PortalRadius   = 75
	PortalCooldown = 5.0

	turnAngle     = math.Pi * 1.5
	moveTurnAngle = math.Pi * 1

	acceleration       = 10.0
	maxVelocity        = 300.0
	maxCollideVelocity = 500
	Braking            = 0.75

	wallElasticity  = 1.1
	BrickElasticity = 0.6

	friction = 0.9

	blinkDistance = 500
	BlinkDuration = 0.4
	BlinkCooldown = 2.0 + BlinkDuration

	HookCooldown         = 3.0
	MaxHookLength        = 500
	hookVelocity         = 750
	hookBackwardVelocity = 750
	hookDamage           = 50

	untouchableTime = 2.0

	MaxHP       = 100.0
	RespawnTime = 2.5
)
