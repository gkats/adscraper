package scraper

import (
	"database/sql"

	"github.com/gkats/scraper/keywords"
)

type Ad struct {
	ID        int64
	H1        string
	H2        string
	Path      string
	Desc      string
	Rest      sql.NullString
	Raw       sql.NullString
	Position  int
	CreatedAt string
	UpdatedAt string
}

func (ad *Ad) GetRaw() string {
	if ad.Raw.Valid {
		return ad.Raw.String
	}
	return ""
}

func (ad *Ad) SetRaw(s string) {
	ad.Raw = sql.NullString{String: s, Valid: true}
}

func (ad *Ad) GetRest() string {
	if ad.Rest.Valid {
		return ad.Rest.String
	}
	return ""
}

func (ad *Ad) SetRest(s string) {
	ad.Rest = sql.NullString{String: s, Valid: true}
}

type AdKeyword struct {
	ID            int64
	AdId          int64
	KeywordId     int64
	Position      int
	PositionCount int
	CreatedAt     string
	UpdatedAt     string
}

type AdWriter interface {
	Upsert(*Ad, *keywords.Keyword) error
}

func newAdKeyword(a *Ad, k *keywords.Keyword) *AdKeyword {
	return &AdKeyword{AdId: a.ID, KeywordId: k.ID, Position: a.Position}
}

func NewWriter(s Store) AdWriter {
	return &adsStore{s}
}

type adsStore struct {
	Store
}

func (s *adsStore) Upsert(ad *Ad, k *keywords.Keyword) error {
	if existing, err := s.findAdByH1H2Desc(ad.H1, ad.H2, ad.Desc); err != nil {
		return err
	} else if existing != nil {
		existing.Position = ad.Position
		return s.save(existing, k)
	}
	return s.save(ad, k)
}

func (s *adsStore) save(ad *Ad, k *keywords.Keyword) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}
	if ad.ID == 0 {
		err = tx.QueryRow(
			`
	    INSERT INTO ads (headline1, headline2, path, description, rest, raw)
	    VALUES($1, $2, $3, $4, $5, $6)
	    RETURNING id
	    `,
			ad.H1, ad.H2, ad.Path, ad.Desc, ad.GetRest(), ad.GetRaw(),
		).Scan(&ad.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	ak := newAdKeyword(ad, k)
	err = tx.QueryRow(
		`
    INSERT INTO ad_keywords (ad_id, keyword_id, position)
    VALUES($1, $2, $3)
    ON CONFLICT (ad_id, keyword_id, position)
    DO UPDATE SET position_count = EXCLUDED.position_count + 1
    RETURNING id
    `,
		ak.AdId, ak.KeywordId, ak.Position,
	).Scan(&ak.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *adsStore) findAdByH1H2Desc(h1 string, h2 string, desc string) (*Ad, error) {
	ad := &Ad{}
	err := s.QueryRow(
		`
    SELECT id, headline1, headline2, description, path, rest, raw, created_at, updated_at
    FROM ads
    WHERE headline1 = $1
    AND headline2 = $2
    AND description = $3
    `,
		h1, h2, desc,
	).Scan(
		&ad.ID, &ad.H1, &ad.H2, &ad.Desc, &ad.Path, &ad.Rest, &ad.Raw,
		&ad.CreatedAt, &ad.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return ad, err
}
