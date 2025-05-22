package error_utils

import (
	"database/sql"
	"log"
	"strings"
)

func HasNoRow(err error) bool {
	return strings.Contains(err.Error(), "no rows")
}

func IsDuplicate(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}

func InvalidInput(err error) bool {
	return strings.Contains(err.Error(), "invalid input")
}

func IsInvalidForeignKey(err error) bool {
	return strings.Contains(err.Error(), "foreign key")
}

func HasNoAffectedRow(sql sql.Result) bool {
	affectedRow, err := sql.RowsAffected()
	if err != nil {
		log.Println("has-no-affected-row-err", err.Error())
		return false
	}

	return affectedRow == 0
}
