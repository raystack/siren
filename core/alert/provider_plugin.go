package alert

import (
	"context"
)

type AlertTransformer interface {
	TransformToAlerts(ctx context.Context, providerID uint64, namespaceID uint64, body map[string]any) ([]Alert, int, error)
}
