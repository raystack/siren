package cortex

import "context"

type contextKeyType struct{}

var (
	// tenantContextKey is the key used for cortex.FromContext and
	// cortex.NewContext.
	tenantContextKey = contextKeyType(struct{}{})
)

// NewContext returns a new context.Context that carries the provided
// tenant ID.
func NewContext(ctx context.Context, tenantId string) context.Context {
	return context.WithValue(ctx, tenantContextKey, tenantId)
}

// FromContext returns the tenant ID from the context if present, and empty
// otherwise.
func FromContext(ctx context.Context) string {
	if t, ok := ctx.Value(tenantContextKey).(string); ok {
		return t
	}
	return ""
}
