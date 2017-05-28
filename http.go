package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseUrl string
	*http.Client
}

func NewClient(host string) *Client {
	return &Client{
		baseUrl: host,
		Client:  &http.Client{Timeout: time.Second * 30},
	}
}

func (c *Client) PostAdKeywords(ad *Ad, k *Keyword) error {
	body, err := json.Marshal(newAdWithKeywordJson(ad, k))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/ad_keywords", bytes.NewBuffer(body))
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

type adJson struct {
	H1       string `json:"h1"`
	H2       string `json:"h2"`
	Desc     string `json:"desc"`
	Path     string `json:"path"`
	Raw      string `json:"raw"`
	Rest     string `json:"rest"`
	Position int    `json:"position"`
}

type keywordJson struct {
	Value string `json:"value"`
}

type adWithKeywordJson struct {
	Ad      *adJson      `json:"ad"`
	Keyword *keywordJson `json:"keyword"`
}

func newAdWithKeywordJson(ad *Ad, keyword *Keyword) *adWithKeywordJson {
	return &adWithKeywordJson{
		Ad:      newAdJson(ad),
		Keyword: newKeywordJson(keyword),
	}
}

func newAdJson(ad *Ad) *adJson {
	return &adJson{
		H1:       ad.H1,
		H2:       ad.H2,
		Desc:     ad.Desc,
		Path:     ad.Path,
		Raw:      ad.GetRaw(),
		Rest:     ad.GetRest(),
		Position: ad.Position,
	}
}

func newKeywordJson(k *Keyword) *keywordJson {
	return &keywordJson{Value: k.Value}
}

type server struct {
	store *Store
}

func NewServer(store *Store) *server {
	return &server{store}
}

func (s *server) Listen(port int) {
	mux := http.NewServeMux()
	mux.Handle("/", root())
	mux.Handle("/ad_keywords", create())
	http.ListenAndServe(":"+strconv.Itoa(port), jsonContent(mux))
}

type createHandler struct{}

func (h *createHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeResponse(w, badRequest())
	}

	// read ad, read keyword from params
	// save ad, save ad_keyword
	writeResponse(w, &successResponse{status: 201})
}

func create() http.Handler {
	return &createHandler{}
}

type response interface {
	Status() int
	Body() interface{}
}

func writeResponse(w http.ResponseWriter, r response) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(r.Status())
	encoder.Encode(r.Body())
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

func badRequest() response {
	return &errorResponse{status: http.StatusBadRequest, Message: "Bad request"}
}

func notFound() response {
	return &errorResponse{status: http.StatusNotFound, Message: "Not found"}
}

type successResponse struct {
	body   interface{}
	status int
}

func (r *successResponse) Status() int {
	return r.status
}

func (r *successResponse) Body() interface{} {
	return r.Body
}

func root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, notFound())
	}
}

func jsonContent(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
