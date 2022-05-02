package db

import (
	"database/sql"
	"fmt"

	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-sqlite3/driver"
	_ "modernc.org/sqlite"

	"github.com/volatiletech/sqlboiler/boil"
)

var DB *sql.DB

func Initialize() {
	var err error
	DB, err = sql.Open("sqlite", fmt.Sprintf("file:%s?_loc=UTC", "devel.db"))
	if err != nil {
		panic(err)
	}
	DB.Ping()
	boil.SetDB(DB)
}
