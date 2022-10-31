package cortex

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestContext(t *testing.T) {
	t.Run("should return passed tenant id if exist in context", func(t *testing.T) {
		var (
			passedTenantID = "some-tenant-id"
			ctx            = newContextWithTenantID(context.Background(), passedTenantID)
			actualTenantID = tenantIDFromContext(ctx)
		)
		if !cmp.Equal(passedTenantID, actualTenantID) {
			t.Fatalf("actual is \"%+v\" but expected was \"%+v\"", actualTenantID, passedTenantID)
		}
	})

	t.Run("should return empty tenant id if not exist in context", func(t *testing.T) {
		actual := tenantIDFromContext(context.Background())
		if actual != "" {
			t.Fatalf("actual is \"%+v\" but expected was \"%+v\"", actual, "")
		}
	})

	t.Run("should return empty tenant id if context is nil", func(t *testing.T) {
		actual := tenantIDFromContext(context.TODO())
		if actual != "" {
			t.Fatalf("actual is \"%+v\" but expected was \"%+v\"", actual, "")
		}
	})
}
