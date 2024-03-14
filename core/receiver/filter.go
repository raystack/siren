package receiver

type Filter struct {
	ReceiverIDs    []uint64
	MultipleLabels []map[string]string
	Expanded       bool
}
