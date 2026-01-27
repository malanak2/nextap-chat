package ports

import (
	"database/sql"
	"fmt"
	"os"
)

var Db *sql.DB

func InitDb() error {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_host"), os.Getenv("DB_port"), os.Getenv("DB_user"), os.Getenv("DB_pass"), os.Getenv("DB_name"))
	var err error
	Db, err = sql.Open("postgres", conn)
	return err
}
