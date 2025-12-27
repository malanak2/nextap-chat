package domain

import "database/sql"

var Db *sql.DB

func InitDb(conn string) error {
	var err error
	Db, err = sql.Open("postgres", conn)
	return err
}
