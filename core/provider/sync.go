package provider

type SyncSubscriptionMethod string

const (
	TypeSyncBatch  SyncSubscriptionMethod = "batch"
	TypeSyncSingle SyncSubscriptionMethod = "single"
)

func (sm SyncSubscriptionMethod) String() string {
	return string(sm)
}

const DefaultSyncSubscriptionMethod = TypeSyncBatch
