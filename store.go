package scraper

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store interface {
	Close() error
	QueryRow(string, ...interface{}) *sql.Row
	Query(string, ...interface{}) (*sql.Rows, error)
	Begin() (*sql.Tx, error)
}

type store struct {
	*sql.DB
}

func NewStore(url string) (Store, error) {
	db, err := sql.Open("postgres", "postgres://"+url+"?sslmode=require")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &store{DB: db}, nil
}
