package gotest_helpers

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	sqlite "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"sync"
)

type customFunc struct {
	name string
	impl interface{}
	pure bool
}

type DatabaseBuilder struct {
	id           string
	testdir      string
	registerOnce sync.Once
	funcs        []customFunc
	extensions   []string
}

func NewDatabaseBuilder(testdir string) *DatabaseBuilder {
	return &DatabaseBuilder{
		id:      uuid.NewV4().String(),
		testdir: testdir,
	}
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// skip initial whitespace and comments
	start := 0
	for start < len(data) {
		if data[start] == ' ' || data[start] == '\t' || data[start] == '\n' || data[start] == '\r' {
			start++
		} else if data[start] == '-' && data[start+1] == '-' {
			if i := bytes.IndexByte(data[start:], '\n'); i >= 0 || atEOF {
				start += i
			} else {
				return 0, nil, nil
			}
		} else {
			break
		}
	}

	if i := bytes.IndexByte(data[start:], ';'); i > 0 {
		return start + i + 1, data[start : start+i+1], nil
	}

	return 0, nil, nil
}

// RegisterFn adds a custom Go function to the database.  It may be called more than
// once to add multiple functions to the database.  Functions must be added before the
// first call to BuildDatabase, however -- subsequent function registrations will be
// ignored.
func (d *DatabaseBuilder) RegisterFn(name string, fn interface{}, pure bool) {
	d.funcs = append(d.funcs, customFunc{
		name: name,
		impl: fn,
		pure: pure,
	})
}

func (d *DatabaseBuilder) RegisterExtension(name string) {
	d.extensions = append(d.extensions, name)
}

func (d DatabaseBuilder) BuildDatabase(name string, sources ...string) (db *sql.DB, err error) {
	dbpath := fmt.Sprintf("%s/%s.db", d.testdir, name)
	os.Remove(dbpath)

	drivername := fmt.Sprintf("sqlite3_%s", d.id)

	d.registerOnce.Do(func() {
		sql.Register(drivername, &sqlite.SQLiteDriver{
			Extensions: d.extensions,
			ConnectHook: func(conn *sqlite.SQLiteConn) error {
				for _, f := range d.funcs {
					if err := conn.RegisterFunc(f.name, f.impl, f.pure); err != nil {
						log.Printf("error registering function %s: %s", name, err.Error())
						return err
					}
					log.Printf("registered function %s", f.name)
				}
				return nil
			},
		})
	})

	db, err = sql.Open(drivername, dbpath)
	if err != nil {
		return nil, err
	}

	for _, source := range sources {
		f, err := os.Open(source)
		if err != nil {
			log.Println(err.Error())
			break
		}
		log.Println("reading", source)
		scanner := bufio.NewScanner(f)
		scanner.Split(split)
		for scanner.Scan() {
			sql := scanner.Text()
			if _, err := db.Exec(sql); err != nil {
				db.Close()
				os.Remove(dbpath)

				return nil, err
			}
		}
	}

	return
}
