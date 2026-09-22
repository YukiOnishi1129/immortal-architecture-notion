package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"immortal-architecture-notion/backend/internal/adapter/gateway/db/sqlc/generated"
	driverdb "immortal-architecture-notion/backend/internal/driver/db"
)

func toUUID(str string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{}, err
	}
	var id pgtype.UUID
	id.Bytes = parsed
	id.Valid = true
	return id, nil
}

func uuidToString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	val, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		return ""
	}
	return val.String()
}

func timestamptzToTime(t pgtype.Timestamptz) time.Time {
	return t.Time
}

func nullableTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// textToStringPtr keeps SQL NULL distinct from an empty string.
// Used for optional columns the domain models as *string.
func textToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

// timestamptzToTimePtr returns nil for SQL NULL rather than the zero time.
func timestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	v := t.Time
	return &v
}

func queriesForContext(ctx context.Context, q *generated.Queries) *generated.Queries {
	if tx := driverdb.TxFromContext(ctx); tx != nil {
		return q.WithTx(tx)
	}
	return q
}

func pgNullableText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// pgTextFromString maps an empty string to SQL NULL.
// Used for optional columns the domain models as a plain string.
func pgTextFromString(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func pgNullableTime(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func timeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
