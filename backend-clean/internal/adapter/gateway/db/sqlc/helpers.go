package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/generated"
	driverdb "immortal-architecture-cqrs/backend/internal/driver/db"
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

func pgNullableTime(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func timeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
