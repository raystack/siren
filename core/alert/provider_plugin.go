package alert

import (
	"context"
)

//go:generate mockery --name=AlertTransformer -r --case underscore --with-expecter --structname AlertTransformer --filename alert_transformer.go --output=./mocks
type AlertTransformer interface {
	TransformToAlerts(ctx context.Context, providerID uint64, namespaceID uint64, body map[string]any) ([]Alert, int, error)
}
