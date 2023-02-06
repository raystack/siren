package model

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/odpf/siren/core/log"
)

type NotificationLog struct {
	ID             string         `db:"id"`
	NamespaceID    sql.NullInt64  `db:"namespace_id"`
	NotificationID string         `db:"notification_id"`
	SubscriptionID uint64         `db:"subscription_id"`
	ReceiverID     sql.NullInt64  `db:"receiver_id"`
	AlertIDs       pq.Int64Array  `db:"alert_ids"`
	SilenceIDs     pq.StringArray `db:"silence_ids"`
	CreatedAt      time.Time      `db:"created_at"`
}

func (ns *NotificationLog) FromDomain(d log.Notification) {
	ns.ID = d.ID

	if d.NamespaceID == 0 {
		ns.NamespaceID = sql.NullInt64{Valid: false}
	} else {
		ns.NamespaceID = sql.NullInt64{Valid: true, Int64: int64(d.NamespaceID)}
	}

	ns.NotificationID = d.NotificationID
	ns.SubscriptionID = d.SubscriptionID
	ns.AlertIDs = pq.Int64Array(d.AlertIDs)
	ns.SilenceIDs = pq.StringArray(d.SilenceIDs)

	if d.ReceiverID == 0 {
		ns.ReceiverID = sql.NullInt64{Valid: false}
	} else {
		ns.ReceiverID = sql.NullInt64{Valid: true, Int64: int64(d.ReceiverID)}
	}

	ns.CreatedAt = d.CreatedAt
}

func (ns *NotificationLog) ToDomain() log.Notification {
	return log.Notification{
		ID:             ns.ID,
		NamespaceID:    uint64(ns.NamespaceID.Int64),
		NotificationID: ns.NotificationID,
		SubscriptionID: ns.SubscriptionID,
		ReceiverID:     uint64(ns.ReceiverID.Int64),
		AlertIDs:       ns.AlertIDs,
		SilenceIDs:     ns.SilenceIDs,
		CreatedAt:      ns.CreatedAt,
	}
}
