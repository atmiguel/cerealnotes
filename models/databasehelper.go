package models

import (
	// "database/sql"
	"errors"
	"github.com/lib/pq"
)

// UniqueConstraintError is returned when a uniqueness constraint is violated during an insert.
var UniqueConstraintError = errors.New("postgres: unique constraint violation")

// QueryResultContainedMultipleRowsError is returned when a query unexpectedly returns more than one row.
var QueryResultContainedMultipleRowsError = errors.New("query result unexpectedly contained multiple rows")

// QueryResultContainedNoRowsError is returned when a query unexpectedly returns no rows.
var QueryResultContainedNoRowsError = errors.New("query result unexpectedly contained no rows")

func convertPostgresError(err error) error {
	const uniqueConstraintErrorCode = "23505"

	if postgresErr, ok := err.(*pq.Error); ok {
		if postgresErr.Code == uniqueConstraintErrorCode {
			return UniqueConstraintError
		}
	}

	return err
}

func (db *DB) execOneResult(sqlQuery string, object interface{}, args ...interface{}) error {

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return convertPostgresError(err)
	}
	defer rows.Close()

	foundResult := false
	for rows.Next() {

		if foundResult {
			return QueryResultContainedMultipleRowsError
		}

		if err := rows.Scan(object); err != nil {
			return convertPostgresError(err)
		}

		foundResult = true
	}

	if !foundResult {
		return QueryResultContainedNoRowsError
	}

	if err := rows.Err(); err != nil {
		return convertPostgresError(err)
	}

	return nil
}

func (db *DB) execNoResults(sqlQuery string, args ...interface{}) (int64, error) {

	res, err := db.Exec(sqlQuery, args...)
	if err != nil {
		return 0, convertPostgresError(err)
	}

	numAffected, err := res.RowsAffected()
	if err != nil {
		return 0, convertPostgresError(err)
	}

	return numAffected, nil
}
