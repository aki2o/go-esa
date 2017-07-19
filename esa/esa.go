package esa

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// defaultBaseURL esa API の host
	defaultBaseURL = "https://api.esa.io"
)

// Client esa API クライアント
type Client struct {
	Client   *http.Client
	apiKey   string
	baseURL  string
	Team     *TeamService
	Stats    *StatsService
	Post     *PostService
	Comment  *CommentService
	Members  *MembersService
	Category *CategoryService
	User     *UserService
}

// NewClient esa APIのClientを生成する
func NewClient(apikey string) *Client {
	c := &Client{}
	c.Client = http.DefaultClient
	c.apiKey = apikey
	c.baseURL = defaultBaseURL
	c.Team = &TeamService{client: c}
	c.Stats = &StatsService{client: c}
	c.Post = &PostService{client: c}
	c.Comment = &CommentService{client: c}
	c.Members = &MembersService{client: c}
	c.Category = &CategoryService{client: c}
	c.User = &UserService{client: c}

	return c
}

func (c *Client) createURL(esaURL string) string {
	return c.baseURL + esaURL + "?access_token=" + c.apiKey
}

func (c *Client) post(esaURL string, bodyType string, body io.Reader, v interface{}) (resp *http.Response, err error) {
	res, err := c.Client.Post(c.createURL(esaURL), bodyType, body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(res.StatusCode))
	}

	if err := responseUnmarshal(res.Body, v); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) patch(esaURL string, bodyType string, body io.Reader, v interface{}) (resp *http.Response, err error) {
	path := c.createURL(esaURL)
	req, err := http.NewRequest("PATCH", path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", bodyType)
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(res.StatusCode))
	}

	if err := responseUnmarshal(res.Body, v); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) delete(esaURL string) (resp *http.Response, err error) {
	path := c.createURL(esaURL)
	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return nil, errors.New(http.StatusText(res.StatusCode))
	}

	return res, nil
}

func (c *Client) get(esaURL string, query url.Values, v interface{}) (resp *http.Response, err error) {
	path := c.createURL(esaURL)
	queries := query.Encode()
	if len(queries) != 0 {
		path += "&" + queries
	}

	res, err := c.Client.Get(path)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(res.StatusCode))
	}

	if err := responseUnmarshal(res.Body, v); err != nil {
		return nil, err
	}

	return res, err
}

func responseUnmarshal(body io.ReadCloser, v interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
