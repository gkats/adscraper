package keywords

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Writer interface {
	Upsert(*Keyword) (*Keyword, error)
	UpdateScraped(*Keyword) (*Keyword, error)
}

func NewWriter(s Store) Writer {
	return &repository{s}
}

type Reader interface {
	GetLeastScraped(int) ([]Keyword, error)
}

func NewReader(s Store) Reader {
	return &repository{s}
}

type ReaderWriter interface {
	Reader
	Writer
}

func NewReaderWriter(s Store) ReaderWriter {
	return &repository{s}
}

type Keyword struct {
	Id            int64
	Value         string
	TimesScraped  int
	CreatedAt     string
	UpdatedAt     string
	LastScrapedAt string
}

func New(value string) *Keyword {
	return &Keyword{Value: value}
}

type Store interface {
	Close() error
	QueryRow(string, ...interface{}) *sql.Row
	Query(string, ...interface{}) (*sql.Rows, error)
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

type store struct {
	*sql.DB
}

type repository struct {
	Store
}

func (r *repository) Upsert(k *Keyword) (*Keyword, error) {
	err := r.Store.QueryRow(
		`
    INSERT INTO keywords (value)
    VALUES ($1)
    ON CONFLICT (value) DO UPDATE SET value = $2
    RETURNING id, created_at, updated_at, times_scraped
    `,
		k.Value, k.Value,
	).Scan(&k.Id, &k.CreatedAt, &k.UpdatedAt, &k.TimesScraped)

	return k, err
}

func (r *repository) GetLeastScraped(limit int) ([]Keyword, error) {
	ks := make([]Keyword, 0)

	rows, err := r.Store.Query(
		`
	   SELECT id, value, times_scraped, last_scraped_at, created_at, updated_at
	   FROM keywords
	   ORDER BY times_scraped ASC
	   LIMIT $1
	   `,
		limit,
	)
	defer rows.Close()
	if err != nil {
		return ks, err
	}

	k := Keyword{}
	for rows.Next() {
		rows.Scan(
			&k.Id, &k.Value, &k.TimesScraped, &k.LastScrapedAt, &k.CreatedAt, &k.UpdatedAt,
		)
		ks = append(ks, k)
	}
	return ks, nil
}

func (r *repository) UpdateScraped(k *Keyword) (*Keyword, error) {
	err := r.Store.QueryRow(
		`
    UPDATE keywords
    SET times_scraped = times_scraped + 1, last_scraped_at = NOW()
    WHERE id = $1
    RETURNING times_scraped
    `,
		k.Id,
	).Scan(&k.TimesScraped)

	return k, err
}
