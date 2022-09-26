package cortex_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/plugins/providers/cortex"
)

func TestContext(t *testing.T) {
	t.Run("should return passed tenant id if exist in context", func(t *testing.T) {
		var (
			passedTenantID = "some-tenant-id"
			ctx            = cortex.NewContextWithTenantID(context.Background(), passedTenantID)
			actualTenantID = cortex.TenantIDFromContext(ctx)
		)
		if !cmp.Equal(passedTenantID, actualTenantID) {
			t.Fatalf("actual is \"%+v\" but expected was \"%+v\"", actualTenantID, passedTenantID)
		}
	})

	t.Run("should return empty tenant id if not exist in context", func(t *testing.T) {
		actual := cortex.TenantIDFromContext(context.Background())
		if actual != "" {
			t.Fatalf("actual is \"%+v\" but expected was \"%+v\"", actual, "")
		}
	})

	t.Run("should return empty tenant id if context is nil", func(t *testing.T) {
		actual := cortex.TenantIDFromContext(context.TODO())
		if actual != "" {
			t.Fatalf("actual is \"%+v\" but expected was \"%+v\"", actual, "")
		}
	})
}
