package notification

// HandlerOption is an option to customize handler creation
type HandlerOption func(*Handler)

// HandlerWithBatchSize sets created handler with the specified batch size
func HandlerWithBatchSize(bs int) HandlerOption {
	return func(w *Handler) {
		w.batchSize = bs
	}
}

// HandlerWithIdentifier sets created handler with the specified batch size
func HandlerWithIdentifier(identifier string) HandlerOption {
	return func(w *Handler) {
		w.identifier = identifier
	}
}
