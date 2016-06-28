package schemaregistry

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type D func(req *http.Request) (*http.Response, error)

func (d D) Do(req *http.Request) (*http.Response, error) {
	return d(req)
}

func TestSubjects(t *testing.T) {
	d := D(func(req *http.Request) (*http.Response, error) {
		if req.Method != "GET" {
			t.Error()
		}
		var resp http.Response
		resp.StatusCode = 200
		resp.Body = ioutil.NopCloser(strings.NewReader(`["a","b"]`))
		return &resp, nil
	})
	c := Client{url.URL{}, d}
	subs, err := c.Subjects()
	if err != nil {
		t.Error()
	}
	if len(subs) != 2 || subs[0] != "a" || subs[1] != "b" {
		t.Error()
	}
}
