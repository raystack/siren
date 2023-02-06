package subscription

type Filter struct {
	NamespaceID       uint64
	Match             map[string]string
	NotificationMatch map[string]string
	SilenceID         string
	IDs               []int64
}
