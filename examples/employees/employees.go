package main

import (
	"database/sql"
	"log"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db}
}

func (s Storage) SalariesBelow(salary uint) map[int]int {
	rows, err := s.db.Query("select emp_no, max(salary) max_salary from salaries group by emp_no having max(salary) < $1", salary)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	result := make(map[int]int)
	for rows.Next() {
		var emp_no, salary int
		if err := rows.Scan(&emp_no, &salary); err == nil {
			result[emp_no] = salary
		} else {
			log.Println(err.Error())
		}
	}
	return result
}

func main() {
	log.Println("hello, world!")
}
