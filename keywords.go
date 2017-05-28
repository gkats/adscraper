package scraper

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
