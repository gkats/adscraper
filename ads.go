package scraper

import (
	"database/sql"
)

type Ad struct {
	Id        int64
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
