package persistence

import (
	"database/sql"
)

type Database struct {
	*Employees
}

func New(pool *sql.DB) *Database {
	return &Database{
		Employees: NewEmployeePersistence(pool),
	}
}
