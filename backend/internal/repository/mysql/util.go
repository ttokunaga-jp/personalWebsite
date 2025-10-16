package mysql

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
)

func toLocalizedText(ja, en sql.NullString) model.LocalizedText {
	return model.LocalizedText{
		Ja: strings.TrimSpace(ja.String),
		En: strings.TrimSpace(en.String),
	}
}

func nullableString(val sql.NullString) string {
	if val.Valid {
		return strings.TrimSpace(val.String)
	}
	return ""
}

func nullString(value string) sql.NullString {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

func nullInt(value *int) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*value), Valid: true}
}

func nullTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: (*value).UTC(), Valid: true}
}

func nullableTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

func rollbackOnError(tx *sqlx.Tx, err *error) {
	if tx == nil {
		return
	}
	if err == nil {
		_ = tx.Rollback()
		return
	}
	if *err != nil {
		_ = tx.Rollback()
	}
}
