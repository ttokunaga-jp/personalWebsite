package mysql

import (
	"context"
	"fmt"
	"strings"

	_ "embed"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//go:embed schema.sql
var schemaDDL string

func applySchema(ctx context.Context, db *sqlx.DB) error {
	if db == nil {
		return fmt.Errorf("apply schema: nil database handle")
	}

	statements := splitStatements(schemaDDL)
	for _, statement := range statements {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, statement); err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				switch mysqlErr.Number {
				case 1050, // table already exists
					1051, // unknown table (when dropping)
					1060, // duplicate column
					1061, // duplicate key name
					1091: // can't drop; check exists
					continue
				}
			}
			return fmt.Errorf("apply schema statement: %w", err)
		}
	}
	return nil
}

func splitStatements(sql string) []string {
	builder := strings.Builder{}
	statements := make([]string, 0, 16)

	inSingleQuote := false
	inDoubleQuote := false
	for _, r := range sql {
		switch r {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case ';':
			if !inSingleQuote && !inDoubleQuote {
				statements = append(statements, builder.String())
				builder.Reset()
				continue
			}
		}
		builder.WriteRune(r)
	}

	if tail := strings.TrimSpace(builder.String()); tail != "" {
		statements = append(statements, tail)
	}

	return statements
}
