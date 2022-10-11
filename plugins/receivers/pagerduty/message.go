package pagerduty

type EvenAction string

const (
	EvenActionTrigger     EvenAction = "trigger"
	EvenActionAcknowledge EvenAction = "acknowledge"
	EvenActionResolve     EvenAction = "resolve"
)
