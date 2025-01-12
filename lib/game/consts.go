package game

import "math"

const (
	FieldWidth  = 3200
	FieldHeight = 2400

	Radius = 30

	PortalRadius     = 75
	TeleportDuration = 0.8
	PortalCooldown   = 3.0 + TeleportDuration

	turnAngle     = math.Pi * 1.5
	moveTurnAngle = math.Pi * 1

	acceleration       = 100.0
	maxVelocity        = 175.0
	boostAcceleration  = 100
	maxBoostVelocity   = 300.0
	maxCollideVelocity = 500.0
	Braking            = 0.5

	wallElasticity  = 5
	BrickElasticity = 0.6

	blinkDistance = 500
	BlinkDuration = 0.4
	BlinkCooldown = 2.0 + BlinkDuration

	HookCooldown         = 3.0
	MaxHookLength        = 400
	hookVelocity         = 500
	hookBackwardVelocity = 500
	hookDamage           = 50

	untouchableTime = 2.0

	MaxHP       = 100.0
	RespawnTime = 2.5
)
