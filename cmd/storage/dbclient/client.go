package dbclient

import (
	"context"
	"database/sql"
	_ "github.com/microsoft/go-mssqldb"
)

func Pool() (*sql.DB, error) {
	c, err := getConfig()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlserver", c.Uri)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
