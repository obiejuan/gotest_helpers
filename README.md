# GoTest Helpers

Functions for setting up a sqlite3 database for unit tests.

I've been writing Java code for an awful long time.  When I started working in Go, two
of the things that I missed were the ability to run my unit tests against on a small, embedded
SQL database like [HyperSQL](http://hsqldb.org) and using [DBUnit](http://dbunit.sourceforge.net)
to load the database.  While Go doesn't have an embedded, SQL-compliant database like HyperSQL,
it does support [SQLite](https://www.sqlite.org) very well.  The purpose of the GoTest Helpers
library is to load data from SQL scripts into SQLite and build a consistent `database.sql.DB`
object for use in unit tests.

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

### License

GoTest Helpers is licensed under the Apache License 2.0