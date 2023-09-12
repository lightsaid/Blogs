package dbrepo_test

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/lightsaid/blogs/dbrepo"
	_ "github.com/mattn/go-sqlite3"
)

var testRepo dbrepo.Repository
var testDB *sqlx.DB

func TestMain(m *testing.M) {
	db, err := sqlx.Connect("sqlite3", "../db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	testDB = db
	testRepo = *dbrepo.NewRepository(db)

	os.Exit(m.Run())
}
