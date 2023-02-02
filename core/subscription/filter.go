package subscription

type Filter struct {
	NamespaceID uint64
	Labels      map[string]string
}
