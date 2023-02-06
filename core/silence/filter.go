package silence

type Filter struct {
	ID                string
	NamespaceID       uint64
	SubscriptionID    uint64
	Match             map[string]string
	SubscriptionMatch map[string]string
}
