package game

import "math"

const (
	FieldWidth  = 2000
	FieldHeight = 1500

	Radius = 30

	PortalRadius   = 75
	PortalCooldown = 5.0

	turnAngle     = math.Pi * 1.75
	moveTurnAngle = math.Pi * 1

	acceleration       = 10.0
	strafeAcceleration = 7.5
	maxVelocity        = 300.0
	maxCollideVelocity = 500
	maxStrafeVelocity  = 250

	wallElasticity  = 1.1
	BrickElasticity = 0.6

	friction       = 0.9
	strafeFriction = 0.7

	blinkDistance = 500
	BlinkDuration = 0.5
	BlinkCooldown = 4.0 + BlinkDuration

	HookCooldown         = 5.0
	hookDistance         = 750
	hookMinDistance      = 100
	hookVelocity         = 750
	hookBackwardVelocity = 750

	untouchableTime = 2.0
)
