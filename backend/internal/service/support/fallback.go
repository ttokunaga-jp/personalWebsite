package support

import (
	"errors"
	"strings"

	mysqlerr "github.com/go-sql-driver/mysql"

	"github.com/takumi/personal-website/internal/repository"
)

// ShouldFallback determines whether a repository error warrants falling back to canonical in-memory data.
func ShouldFallback(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, repository.ErrNotFound) {
		return true
	}

	var mysqlErr *mysqlerr.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1044, // ER_DBACCESS_DENIED_ERROR
			1049, // ER_BAD_DB_ERROR
			1146: // ER_NO_SUCH_TABLE
			return true
		}
	}

	const sqliteTableMissing = "no such table"
	if strings.Contains(strings.ToLower(err.Error()), sqliteTableMissing) {
		return true
	}

	return false
}
