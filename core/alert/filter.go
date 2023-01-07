package alert

type Filter struct {
	ResourceName string
	ProviderID   uint64
	NamespaceID  uint64
	StartTime    int64
	EndTime      int64
}
