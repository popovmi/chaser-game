package warsgame

const (
	FieldWidth  = 2000
	FieldHeight = 1500

	Radius       = 30
	PortalRadius = 75

	Top    = 0 + Radius
	Bottom = FieldHeight - Radius
	Left   = 0 + Radius
	Right  = FieldWidth - Radius

	Acceleration         = 0.5
	Braking              = 0.85
	Friction             = 0.05
	WallElasticity       = 1.2
	PlayerElasticity     = 1.2
	BrickElasticity      = 0.8
	MaxVelocity          = 6.5
	MaxCollideVelocity   = 10
	UntouchableTime      = 2000
	BlinkDistance        = 350
	BlinkCooldown        = 4000
	PortalCooldown       = 5000
	HookVelocity         = 20
	HookBackwardVelocity = 15
	HookDistance         = 700
	HookCooldown         = 5000

	TPS = 60

	LineSpacing     = 1.1
	TextFieldHeight = 30
	TextFieldWidth  = 282
	MaxTextLength   = 20
)
