package main

import (
	"database/sql"
	"github.com/pegli/gotest_helpers"
	"log"
	"os"
	"testing"
)

var shared_db *sql.DB

func TestMain(m *testing.M) {
	builder := gotest_helpers.NewDatabaseBuilder("testdata")

	db, err := builder.BuildDatabase(
		"test1",
		"testdata/schema.sql",
		"testdata/departments.sql",
		"testdata/employees.sql",
		"testdata/dept_emp.sql",
		"testdata/dept_manager.sql",
		"testdata/titles.sql",
		"testdata/salaries1.sql",
	)
	if err != nil {
		log.Println("error building dataset", err.Error())
		return
	}
	defer db.Close()

	shared_db = db

	os.Exit(m.Run())
}

func TestSalariesBelow(t *testing.T) {
	storage := NewStorage(shared_db)

	salaries := storage.SalariesBelow(50000)
	if salaries == nil {
		t.Fatalf("SalariesBelow returned nil")
	}
	if len(salaries) != 10 {
		t.Fatalf("Expected 10; got %d", len(salaries))
	}
}
