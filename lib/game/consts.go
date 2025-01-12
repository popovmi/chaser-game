package game

import "math"

const (
	FieldWidth  = 3200
	FieldHeight = 2400

	Radius = 30

	PortalRadius     = 75
	TeleportDuration = 0.8
	PortalCooldown   = 3.0 + TeleportDuration

	turnAngle     = math.Pi * 1
	moveTurnAngle = math.Pi * 0.85

	acceleration       = 600.0
	maxVelocity        = 300.0
	boostAcceleration  = 600
	maxBoostVelocity   = 600.0
	maxCollideVelocity = 650.0
	Braking            = 0.5

	wallElasticity  = 1.05
	BrickElasticity = 0.6

	blinkDistance = 500
	BlinkDuration = 0.4
	BlinkCooldown = 2.0 + BlinkDuration

	HookCooldown         = 3.0
	MaxHookLength        = 400
	hookVelocity         = 700
	hookBackwardVelocity = 700
	hookDamage           = 50

	untouchableTime = 2.0

	MaxHP       = 100.0
	RespawnTime = 2.5
)
