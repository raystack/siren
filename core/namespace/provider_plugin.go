package namespace

import (
	"context"

	"github.com/goto/siren/core/provider"
)

//go:generate mockery --name=ConfigSyncer -r --case underscore --with-expecter --structname ConfigSyncer --filename config_syncer.go --output=./mocks
type ConfigSyncer interface {
	SyncRuntimeConfig(ctx context.Context, namespaceID uint64, namespaceURN string, namespaceLabels map[string]string, prov provider.Provider) (map[string]string, error)
}
