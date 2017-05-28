package scraper

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	*sql.DB
}

func NewStore(url string) (*Store, error) {
	return nil, nil
}
