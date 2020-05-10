package schemaregistry

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

const testHost = "testhost:1337"
const testURL = "http://" + testHost
const testUsername = "test"
const testPassword = "password"

type D func(req *http.Request) (*http.Response, error)

func (d D) Do(req *http.Request) (*http.Response, error) {
	return d(req)
}

// verifies the http.Request, creates an http.Response
func dummyHTTPHandler(t *testing.T, authEnabled bool, method, path string, status int, reqBody, respBody interface{}) D {
	d := D(func(req *http.Request) (*http.Response, error) {
		if authEnabled {
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				t.Error("basic auth is enabled but Authorization header is not set")
			}
			basicAuth := basicAuth{testUsername, testPassword}
			if authHeader != basicAuth.String() {
				t.Errorf("basic auth header is enabled but credentials do not match: %s:%s", testUsername, testPassword)
			}
		}
		if method != "" && req.Method != method {
			t.Errorf("method is wrong, expected `%s`, got `%s`", method, req.Method)
		}
		if req.URL.Host != testHost {
			t.Errorf("expected host `%s`, got `%s`", testHost, req.URL.Host)
		}
		if path != "" && req.URL.Path != path {
			t.Errorf("path is wrong, expected `%s`, got `%s`", path, req.URL.Path)
		}
		if reqBody != nil {
			expbs, err := json.Marshal(reqBody)
			if err != nil {
				t.Error(err)
			}
			bs, err := ioutil.ReadAll(req.Body)
			mustEqual(t, strings.Trim(string(bs), "\r\n"), strings.Trim(string(expbs), "\r\n"))
		}
		var resp http.Response
		resp.Header = http.Header{contentTypeHeaderKey: []string{contentTypeJSON}}
		resp.StatusCode = status
		if respBody != nil {
			bs, err := json.Marshal(respBody)
			if err != nil {
				t.Error(err)
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
		}
		return &resp, nil
	})
	return d
}

func httpSuccess(t *testing.T, authEnabled bool, method, path string, reqBody, respBody interface{}) *Client {
	client := &Client{testURL, dummyHTTPHandler(t, authEnabled, method, path, 200, reqBody, respBody), nil}
	if authEnabled {
		client.SetBasicAuth(testUsername, testPassword)
	}
	return client
}

func httpError(t *testing.T, authEnabled bool, status, errCode int, errMsg string) *Client {
	client := &Client{testURL, dummyHTTPHandler(t, authEnabled, "", "", status, nil, ResourceError{ErrorCode: errCode, Message: errMsg}), nil}
	if authEnabled {
		client.SetBasicAuth(testUsername, testPassword)
	}
	return client
}

func mustEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected `%#v`, got `%#v`", expected, actual)
	}
}

var params = []struct {
	name        string
	authEnabled bool
}{
	{"basic auth enabled", true},
	{"basic auth disabled", false},
}

func TestSubjects(t *testing.T) {
	for _, tt := range params {
		t.Run(tt.name, func(t *testing.T) {
			subsIn := []string{"rollulus", "hello-subject"}
			c := httpSuccess(t, tt.authEnabled, "GET", "/subjects", nil, subsIn)
			subs, err := c.Subjects()
			if err != nil {
				t.Error()
			}
			mustEqual(t, subs, subsIn)
		})
	}
}

func TestVersions(t *testing.T) {
	for _, tt := range params {
		t.Run(tt.name, func(t *testing.T) {
			versIn := []int{1, 2, 3}
			c := httpSuccess(t, tt.authEnabled, "GET", "/subjects/mysubject/versions", nil, versIn)
			vers, err := c.Versions("mysubject")
			if err != nil {
				t.Error()
			}
			mustEqual(t, vers, versIn)
		})
	}
}

func TestIsRegistered_yes(t *testing.T) {
	s := `{"x":"y"}`
	ss := schemaOnlyJSON{s}
	sIn := Schema{s, "mysubject", 4, 7}
	for _, tt := range params {
		t.Run(tt.name, func(t *testing.T) {
			c := httpSuccess(t, tt.authEnabled, "POST", "/subjects/mysubject", ss, sIn)
			isreg, sOut, err := c.IsRegistered("mysubject", s)
			if err != nil {
				t.Error()
			}
			if !isreg {
				t.Error()
			}
			mustEqual(t, sOut, sIn)
		})
	}
}

func TestIsRegistered_not(t *testing.T) {
	for _, tt := range params {
		t.Run(tt.name, func(t *testing.T) {
			c := httpError(t, tt.authEnabled, 404, schemaNotFoundCode, "too bad")
			isreg, _, err := c.IsRegistered("mysubject", "{}")
			if err != nil {
				t.Fatal(err)
			}
			if isreg {
				t.Fatalf("is registered: %v", err)
			}
		})
	}
}
