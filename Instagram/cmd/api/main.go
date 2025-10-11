package main

import (
	"backend/db"
	"backend/migration"
	"backend/utils"

	"github.com/jmoiron/sqlx"
)

var (
	DB *sqlx.DB
)

func main() {
	migration.Migrate(DB)
	db.ConnectDB()
	utils.NewClient("amqp://guest:guest@localhost:5672/")

}
