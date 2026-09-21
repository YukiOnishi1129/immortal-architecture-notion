package sqlc

import (
	"testing"
	"time"
)

func TestPgNullableHelpers(t *testing.T) {
	txt := "x"
	if v := pgNullableText(nil); v.Valid {
		t.Fatalf("expected invalid text for nil")
	}
	if v := pgNullableText(&txt); !v.Valid || v.String != txt {
		t.Fatalf("unexpected text: %+v", v)
	}
	now := time.Now().UTC()
	if v := pgNullableTime(nil); v.Valid {
		t.Fatalf("expected invalid time for nil")
	}
	if v := pgNullableTime(&now); !v.Valid || !v.Time.Equal(now) {
		t.Fatalf("unexpected time: %+v", v)
	}
}
