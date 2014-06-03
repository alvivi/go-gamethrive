package gamethrive

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func setup() (server *httptest.Server, mux *http.ServeMux, client *Client) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	return
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	if c.client != http.DefaultClient {
		t.Error("NewClient default client must be http default client")
	}
	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), defaultBaseURL)
	}
	if c.UserAgent != defaultUserAgent {
		t.Errorf("NewClient UserAgent = %v, want %v", c.UserAgent, defaultUserAgent)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)
	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody := &struct{ Bar string }{Bar: "rocks"}
	outBody := `{"Bar":"rocks"}` + "\n"
	req, _ := c.NewRequest("GET", inURL, inBody)
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v) Body = %v, want %v", inBody, string(body), outBody)
	}
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, want %v", userAgent, c.UserAgent)
	}
	contentType := req.Header.Get("Content-Type")
	if expected := "application/json"; contentType != expected {
		t.Errorf("NewRequest() Content-Type = %v, want %v", contentType, expected)
	}
}

func TestDo(t *testing.T) {
	server, mux, client := setup()
	defer server.Close()
	type foo struct {
		Bar string
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; r.Method != m {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprintf(w, `{"Bar":"rocks"}`)
	})
	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	client.Do(req, body)
	want := &foo{"rocks"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	server, mux, client := setup()
	defer server.Close()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})
	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)
	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"errors":["app_id not found."]}`)),
	}
	err := checkResponse(res).(*ErrorResponse)
	if err == nil {
		t.Errorf("Expected error response.")
	}
	want := &ErrorResponse{
		Response: res,
		Errors:   []string{"app_id not found."},
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := checkResponse(res).(*ErrorResponse)
	if err == nil {
		t.Errorf("Expected error response.")
	}
	want := &ErrorResponse{Response: res}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}
