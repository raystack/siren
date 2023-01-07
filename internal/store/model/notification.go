package model

import (
	"database/sql"
	"time"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/pkg/pgc"
)

type Notification struct {
	ID            string                 `db:"id"`
	NamespaceID   sql.NullInt64          `db:"namespace_id"`
	Type          string                 `db:"type"`
	Data          pgc.StringInterfaceMap `db:"data"`
	Labels        pgc.StringStringMap    `db:"labels"`
	ValidDuration pgc.TimeDuration       `db:"valid_duration"`
	Template      sql.NullString         `db:"template"`
	CreatedAt     time.Time              `db:"created_at"`
}

func (n *Notification) FromDomain(d notification.Notification) {
	n.ID = d.ID
	n.Type = d.Type
	n.Data = d.Data
	n.Labels = d.Labels
	n.ValidDuration = pgc.TimeDuration(d.ValidDuration)

	if d.NamespaceID == 0 {
		n.NamespaceID = sql.NullInt64{Valid: false}
	} else {
		n.NamespaceID = sql.NullInt64{Int64: int64(d.NamespaceID), Valid: true}
	}

	if d.Template == "" {
		n.Template = sql.NullString{Valid: false}
	} else {
		n.Template = sql.NullString{String: d.Template, Valid: true}
	}

	n.CreatedAt = d.CreatedAt
}

func (n *Notification) ToDomain() notification.Notification {
	return notification.Notification{
		ID:            n.ID,
		NamespaceID:   uint64(n.NamespaceID.Int64),
		Type:          n.Type,
		Data:          n.Data,
		Labels:        n.Labels,
		ValidDuration: time.Duration(n.ValidDuration),
		Template:      n.Template.String,
		CreatedAt:     n.CreatedAt,
	}
}
