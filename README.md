# GoTest Helpers

Functions for setting up a sqlite3 database for unit tests.

I've been writing Java code for an awful long time.  When I started working in Go, two
of the things that I missed were the ability to run my unit tests against on a small, embedded
SQL database like [HyperSQL](http://hsqldb.org) and using [DBUnit](http://dbunit.sourceforge.net)
to load the database.  While Go doesn't have an embedded, SQL-compliant database like HyperSQL,
it does support [SQLite](https://www.sqlite.org) very well.  The purpose of the GoTest Helpers
library is to load data from SQL scripts into SQLite and build a consistent `database.sql.DB`
object for use in unit tests.

This library will work best when your database code is aimed at the lowest common denominator
of SQL conformance.  Extended types like PostgreSQL's `JSONB` and GIS types will be difficult
to emulate in SQLite, for example.  See the section titled "Common Workarounds" below for a list
of problems I've found (and sometimes overcome) when testing my code against SQLite using this
library.

### Basic Usage

See the `examples` folder for more detailed use cases.

```go
// Create a DatabaseBuilder
builder := gotest_helpers.NewDatabaseBuilder("testdata")

// Ask it to build a *sql.DB with the provided SQL scripts
db, err := builder.BuildDatabase(
	"dbname",
	"testdata/schema.sql",
	"testdata/data.sql",
)
if err != nil {
	log.Println("error building dataset", err.Error())
	return
}
defer db.Close()

// use the *sql.DB as normal
```

### Custom Functions

You can register custom Go functions in the database driver using the `RegisterFn` func
on `DatabaseBuilder`.  If your SQL queries use functions that are native to your target
database but not present in SQLite, you can provide implementations of these functions in
Go.  For example, SQLite doesn't contain most of the Postgres
[JSON functions](http://www.postgresql.org/docs/9.4/static/functions-json.html).  You
can implement `array_to_json` in Go and register it with the `DatabaseBuilder` like this:

```go
func my_array_to_json(text string) string {
	// coincidentally, sqlite returns arrays as ["el1", "el2",...]
	return text
}

builder := gotest_helpers.NewDatabaseBuilder("testdata")
builder.RegisterFn("array_to_json", my_array_to_json, true)
db, err := builder.BuildDatabase(...)
```

Functions must be registered prior to the call to `BuildDatabase`.  Make sure you review
the [RegisterFunc](https://godoc.org/github.com/mattn/go-sqlite3#SQLiteConn.RegisterFunc)
method of the underlying SQLite driver and have a good grasp of SQLite's type system in 
order to write effective functions.

### Common Workarounds

* LastInsertID support

SQLite supports getting the primary key value of an inserted row through the `LastInsertID`
function.  The `lib/pq` driver for Postgres, which is what I use for my production apps, does
not.  Here's a bit of conditional code which gets the ID for both databases:

```go
tx := db.Begin()
if result, err := tx.Exec("INSERT INTO employees (first_name, last_name, hire_date) VALUES ($1, $2, $3)", "Paul", "Mietz Egli", "2010-06-21"); err == nil {
	if id, err = result.LastInsertId(); err != nil {
		if err = tx.QueryRow("SELECT currval('employees_emp_no_seq')").Scan(&id); err != nil {
			log.Println(err.Error())
			tx.Rollback()
      return
		}
	}
}
tx.Commit()
```

### License

GoTest Helpers is licensed under the Apache License 2.0
