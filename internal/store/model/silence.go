package model

import (
	"database/sql"
	"time"

	"github.com/raystack/siren/core/silence"
	"github.com/raystack/siren/pkg/pgc"
)

type Silence struct {
	ID               string                 `db:"id"`
	NamespaceID      uint64                 `db:"namespace_id"`
	Type             string                 `db:"type"`
	TargetID         sql.NullInt64          `db:"target_id"`
	TargetExpression pgc.StringInterfaceMap `db:"target_expression"`
	Creator          sql.NullString         `db:"creator"`
	Comment          sql.NullString         `db:"comment"`
	CreatedAt        time.Time              `db:"created_at"`
	DeletedAt        sql.NullTime           `db:"deleted_at"`
}

func (s *Silence) FromDomain(sil silence.Silence) {
	s.ID = sil.ID
	s.NamespaceID = sil.NamespaceID
	s.Type = sil.Type

	if sil.TargetID == 0 {
		s.TargetID = sql.NullInt64{Valid: false}
	} else {
		s.TargetID = sql.NullInt64{Int64: int64(sil.TargetID), Valid: true}
	}

	s.TargetExpression = pgc.StringInterfaceMap(sil.TargetExpression)

	if sil.Creator == "" {
		s.Creator = sql.NullString{Valid: false}
	} else {
		s.Creator = sql.NullString{String: sil.Creator, Valid: true}
	}

	if sil.Comment == "" {
		s.Comment = sql.NullString{Valid: false}
	} else {
		s.Comment = sql.NullString{String: sil.Comment, Valid: true}
	}

	s.CreatedAt = sil.CreatedAt

	if sil.DeletedAt.IsZero() {
		s.DeletedAt = sql.NullTime{Valid: false}
	} else {
		s.DeletedAt = sql.NullTime{Time: sil.DeletedAt, Valid: true}
	}
}

func (s *Silence) ToDomain() *silence.Silence {
	return &silence.Silence{
		ID:               s.ID,
		NamespaceID:      s.NamespaceID,
		Type:             s.Type,
		TargetID:         uint64(s.TargetID.Int64),
		TargetExpression: s.TargetExpression,
		CreatedAt:        s.CreatedAt,
		DeletedAt:        s.DeletedAt.Time,
	}
}
