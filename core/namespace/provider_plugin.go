package namespace

import (
	"context"

	"github.com/goto/siren/core/provider"
)

type ConfigSyncer interface {
	SyncRuntimeConfig(ctx context.Context, namespaceID uint64, namespaceURN string, namespaceLabels map[string]string, prov provider.Provider) (map[string]string, error)
}
