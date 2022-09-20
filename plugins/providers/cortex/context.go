package cortex

import "context"

type contextKeyType struct{}

var (
	// tenantContextKey is the key used for cortex.FromContext and
	// cortex.NewContext.
	tenantContextKey = contextKeyType(struct{}{})
)

// NewContextWithTenantID returns a new context.Context
// that carries the provided tenant ID.
func NewContextWithTenantID(ctx context.Context, tenantId string) context.Context {
	return context.WithValue(ctx, tenantContextKey, tenantId)
}

// TenantIDFromContext returns the tenant ID from the context
// if present, and empty otherwise.
func TenantIDFromContext(ctx context.Context) string {
	if t, ok := ctx.Value(tenantContextKey).(string); ok {
		return t
	}
	return ""
}
