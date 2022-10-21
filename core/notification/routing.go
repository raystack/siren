package notification

type RoutingMethod string

const (
	RoutingMethodReceiver    RoutingMethod = "receiver"
	RoutingMethodSubscribers RoutingMethod = "subscribers"
)

func (rm RoutingMethod) String() string {
	return string(rm)
}
