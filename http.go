package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gkats/httplog"
	"github.com/gkats/scraper/keywords"
	"github.com/gorilla/mux"
)

type Client struct {
	baseURL string
	*http.Client
}

func NewClient(host string) *Client {
	return &Client{
		baseURL: host,
		Client:  &http.Client{Timeout: time.Second * 30},
	}
}

func (c *Client) PostAdKeywords(ad *Ad, k *keywords.Keyword) error {
	body, err := json.Marshal(newAdWithKeywordJSON(ad, k))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/ad_keywords", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if resp, err := c.Do(req); err != nil {
		return err
	} else if resp.StatusCode > 399 {
		return fmt.Errorf("Got error response (%v)", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetKeywords() ([]*keywords.Keyword, error) {
	var kws []*keywords.Keyword

	req, err := http.NewRequest("GET", c.baseURL+"/keywords", nil)
	if err != nil {
		return kws, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return kws, err
	} else if resp.StatusCode > 399 {
		return kws, fmt.Errorf("Got error response (%v)", resp.StatusCode)
	}
	defer resp.Body.Close()

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return kws, err
	} else if err = json.Unmarshal(body, &kws); err != nil {
		return kws, err
	}

	return kws, nil
}

func (c *Client) PatchKeyword(id int64) error {
	req, err := http.NewRequest("PATCH", c.baseURL+"/keywords/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode > 399 {
		return fmt.Errorf("Got error response (%v)", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}

type adJSON struct {
	H1       string `json:"h1"`
	H2       string `json:"h2"`
	Desc     string `json:"desc"`
	Path     string `json:"path"`
	Raw      string `json:"raw"`
	Rest     string `json:"rest"`
	Position int    `json:"position"`
}

type keywordJSON struct {
	ID            int64  `json:"id"`
	Value         string `json:"value"`
	TimesScraped  int    `json:"timesScraped"`
	LastScrapedAt string `json:"lastScrapedAt"`
}

type adWithKeywordJSON struct {
	Ad      adJSON      `json:"ad"`
	Keyword keywordJSON `json:"keyword"`
}

func newAdWithKeywordJSON(ad *Ad, keyword *keywords.Keyword) *adWithKeywordJSON {
	return &adWithKeywordJSON{
		Ad:      newAdJSON(ad),
		Keyword: newKeywordJSON(keyword),
	}
}

func newAdJSON(ad *Ad) adJSON {
	return adJSON{
		H1:       ad.H1,
		H2:       ad.H2,
		Desc:     ad.Desc,
		Path:     ad.Path,
		Raw:      ad.GetRaw(),
		Rest:     ad.GetRest(),
		Position: ad.Position,
	}
}

func newKeywordJSON(k *keywords.Keyword) keywordJSON {
	return keywordJSON{
		ID:            k.ID,
		Value:         k.Value,
		TimesScraped:  k.TimesScraped,
		LastScrapedAt: k.LastScrapedAt,
	}
}

func (a *adJSON) ToAd() *Ad {
	ad := &Ad{
		H1: a.H1, H2: a.H2, Desc: a.Desc, Path: a.Path, Position: a.Position,
	}
	ad.SetRaw(a.Raw)
	ad.SetRest(a.Rest)
	return ad
}

func (k *keywordJSON) ToKeyword() *keywords.Keyword {
	return &keywords.Keyword{
		ID:            k.ID,
		Value:         k.Value,
		TimesScraped:  k.TimesScraped,
		LastScrapedAt: k.LastScrapedAt,
	}
}

type server struct {
	store  Store
	logger httplog.Logger
}

func NewServer(store Store) *server {
	logger := httplog.New(os.Stdout)
	return &server{store: store, logger: logger}
}

func (s *server) Listen(port int) {
	r := mux.NewRouter()
	r.Handle("/ad_keywords", create(s.store)).Methods("POST")
	r.Handle("/keywords", index(s.store)).Methods("GET")
	r.Handle("/keywords/{id}", update(s.store)).Methods("PATCH", "PUT")
	r.HandleFunc("/", root())
	http.Handle("/", r)
	http.ListenAndServe(":"+strconv.Itoa(port), httplog.WithLogging(jsonContent(r), s.logger))
}

type createHandler struct {
	adWriter       AdWriter
	keywordsWriter keywords.Writer
}

func (h *createHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := &adWithKeywordJSON{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		writeResponse(w, badRequest())
		return
	}

	if err := h.adWriter.Upsert(params.Ad.ToAd(), params.Keyword.ToKeyword()); err != nil {
		writeResponse(w, internalServerError())
		return
	}
	writeResponse(w, &successResponse{status: http.StatusCreated})
}

func create(s Store) http.Handler {
	return &createHandler{adWriter: NewWriter(s), keywordsWriter: keywords.NewWriter(s)}
}

type indexHandler struct {
	keywordsReader keywords.Reader
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	kws, err := h.keywordsReader.GetLeastScraped(20)
	if err != nil {
		writeResponse(w, internalServerError())
		return
	}

	kwsJSON := make([]keywordJSON, 0)
	for _, k := range kws {
		kwsJSON = append(kwsJSON, newKeywordJSON(&k))
	}
	writeResponse(w, ok(kwsJSON))
}

func index(s Store) http.Handler {
	return &indexHandler{keywordsReader: keywords.NewReader(s)}
}

type updateHandler struct {
	keywordsWriter keywords.Writer
}

func (h *updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeResponse(w, badRequest())
		return
	}

	if k, err := h.keywordsWriter.UpdateScraped(&keywords.Keyword{ID: int64(id)}); err != nil {
		writeResponse(w, internalServerError())
		return
	} else {
		writeResponse(w, ok(k))
	}
}

func update(s Store) http.Handler {
	return &updateHandler{keywordsWriter: keywords.NewWriter(s)}
}

type response interface {
	Status() int
	Body() interface{}
}

func writeResponse(w http.ResponseWriter, r response) {
	w.WriteHeader(r.Status())
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(r.Body()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type errorResponse struct {
	Message string `json:"message"`
	status  int
}

func (r *errorResponse) Status() int {
	return r.status
}

func (r *errorResponse) Body() interface{} {
	return r
}

func ok(body interface{}) response {
	return &successResponse{status: http.StatusOK, body: body}
}

func badRequest() response {
	return &errorResponse{status: http.StatusBadRequest, Message: "Bad request"}
}

func notFound() response {
	return &errorResponse{status: http.StatusNotFound, Message: "Not found"}
}

func internalServerError() response {
	return &errorResponse{status: http.StatusInternalServerError, Message: "Something went wrong"}
}

type successResponse struct {
	body   interface{}
	status int
}

func (r *successResponse) Status() int {
	return r.status
}

func (r *successResponse) Body() interface{} {
	return r.body
}

func root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, notFound())
	}
}

func jsonContent(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		h.ServeHTTP(w, r)
	})
}
