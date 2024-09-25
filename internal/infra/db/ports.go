package db

import "database/sql"

// Executor abstracts the database operations needed for saving accounts
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}
