package provider

type SyncMethod string

const (
	TypeSyncBatch  SyncMethod = "batch"
	TypeSyncSingle SyncMethod = "single"
)

func (sm SyncMethod) String() string {
	return string(sm)
}

const DefaultSyncMethod = TypeSyncBatch
