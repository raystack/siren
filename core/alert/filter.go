package alert

type Filter struct {
	ResourceName string
	ProviderID   uint64
	StartTime    uint64
	EndTime      uint64
}
