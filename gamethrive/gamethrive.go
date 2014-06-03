package gamethrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	libraryVersion   = "0.0.1"
	defaultBaseURL   = "https://gamethrive.com/api/v1/"
	defaultUserAgent = "go-gamethrive/" + libraryVersion
)

type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string

	Players       PlayersService
	Notifications NotificationsService
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	client := Client{
		client:    httpClient,
		BaseURL:   mustParse(url.Parse(defaultBaseURL)),
		UserAgent: defaultUserAgent,
	}
	client.Players = PlayersService{&client}
	client.Notifications = NotificationsService{&client}
	return &client
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path.Join(c.BaseURL.Path, urlStr))
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return resp, err
}

type ErrorResponse struct {
	*http.Response
	Errors []string `json: "errors"`
}

func (r ErrorResponse) Error() string {
	errsStr := strings.Join(r.Errors, "; ")
	return fmt.Sprintf("%s %s: (%d) %s",
		r.Response.Request.Method, r.Response.Request.URL.String(),
		r.StatusCode, errsStr)
}

func checkResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}
	err := ErrorResponse{Response: res}
	json.NewDecoder(res.Body).Decode(&err)
	return &err
}
