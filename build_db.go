package testutils

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type DatabaseBuilder struct {
	testdir string
}

func NewDatabaseBuilder(testdir string) *DatabaseBuilder {
	return &DatabaseBuilder{
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

func (d DatabaseBuilder) BuildDatabase(name string, sources ...string) (db *sql.DB, err error) {
	dbpath := fmt.Sprintf("%s/%s.db", d.testdir, name)
	os.Remove(dbpath)
	db, err = sql.Open("sqlite3", dbpath)
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
