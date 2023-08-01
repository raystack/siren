package receiver

type Filter struct {
	ReceiverIDs []uint64
	Labels      map[string]string
	Expanded    bool
}
