package health

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func TestDBChecker(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		t.Errorf("Failed to create temp folder\n%s", err.Error())
	}
	defer os.RemoveAll(dir)

	db := prepareDB(dir, t)
	defer db.Close()
	chk := NewDBChecker("main-db", db, time.Duration(1 * time.Second))

	s := chk.Test()
	t.Log(s)
	if s.Status != CheckerStatusOK {
		t.Errorf("Should return status 'OK', but returned '%s'", s.Status)
	}
}

func TestDBCheckerTimeout(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		t.Errorf("Failed to create temp folder\n%s", err.Error())
	}
	defer os.RemoveAll(dir)

	db := prepareDB(dir, t)
	defer db.Close()
	chk := NewDBCheckerCustomQuery("main-db", db, time.Duration(1 * time.Nanosecond), "SELECT SLEEP(100)")

	s := chk.Test()
	t.Log(s)
	if s.Status != CheckerStatusNOK {
		t.Errorf("Should return status 'DOWN', but returned '%s'", s.Status)
	}
}

func prepareDB(dir string, t *testing.T) *sql.DB {
	script := `CREATE TABLE "log_temperature" (
		"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		"date"  datetime NOT NULL,
		"temperature"  varchar(20) NOT NULL
	)`

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/sqlite-database.db", dir))
	if err != nil {
		t.Errorf("Failed to open DB connection\n%s", err.Error())
	}

	if _, err := db.Exec(script); err != nil {
		t.Errorf("Failed to prepare database\n%s", err.Error())
	}

	return db
}