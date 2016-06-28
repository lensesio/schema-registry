// Package schemaregistry provides a client for Confluent's Kafka Schema Registry REST API.
package schemaregistry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// DefaultUrl is the address where a local schema registry listens by default.
var DefaultUrl = "http://localhost:8081"

// These numbers are used by the schema registry to communicate errors.
const (
	subjectNotFound = 40401
	schemaNotFound  = 40403
)

// The Schema type is an object produced by the schema registry.
type Schema struct {
	Schema  string `json:"schema"`  // The actual AVRO schema
	Subject string `json:"subject"` // Subject where the schema is registered for
	Version int    `json:"version"` // Version within this subject
	Id      int    `json:"id"`      // Registry's unique id
}

type simpleSchema struct {
	Schema string `json:"schema"`
}

// A ConfluentError is an error as communicated by the schema registry.
// Some day this type might be exposed so that callers can do type assertions on it.
type confluentError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// Error makes confluentError implement the error interface.
func (ce confluentError) Error() string {
	return fmt.Sprintf("%s (%d)", ce.Message, ce.ErrorCode)
}

type httpDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// A Client is a client for the schema registry.
type Client struct {
	url    url.URL
	client httpDoer
}

func parseSchemaRegistryError(resp *http.Response) error {
	var ce confluentError
	if err := json.NewDecoder(resp.Body).Decode(&ce); err != nil {
		return err
	}
	return ce
}

// do performs http requests and json (de)serialization.
func (c *Client) do(method, urlPath string, in interface{}, out interface{}) error {
	u := c.url
	u.Path = path.Join(u.Path, urlPath)
	var rdp io.Reader
	if in != nil {
		var wr *io.PipeWriter
		rdp, wr = io.Pipe()
		go func() {
			wr.CloseWithError(json.NewEncoder(wr).Encode(in))
		}()
	}
	req, err := http.NewRequest(method, u.String(), rdp)
	req.Header.Add("Accept", "application/vnd.schemaregistry.v1+json, application/vnd.schemaregistry+json, application/json")
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseSchemaRegistryError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// Subjects returns all registered subjects.
func (c *Client) Subjects() (subjects []string, err error) {
	err = c.do("GET", "subjects", nil, &subjects)
	return
}

// Versions returns all schema version numbers registered for this subject.
func (c *Client) Versions(subject string) (versions []int, err error) {
	err = c.do("GET", fmt.Sprintf("subjects/%s/versions", subject), nil, &versions)
	return
}

// RegisterNewSchema registers the given schema for this subject.
func (c *Client) RegisterNewSchema(subject, schema string) (int, error) {
	var resp struct {
		Id int `json:"id"`
	}
	err := c.do("POST", fmt.Sprintf("/subjects/%s/versions", subject), simpleSchema{schema}, &resp)
	return resp.Id, err
}

// IsRegistered tells if the given schema is registred for this subject.
func (c *Client) IsRegistered(subject, schema string) (bool, Schema, error) {
	var fs Schema
	err := c.do("POST", fmt.Sprintf("/subjects/%s", subject), simpleSchema{schema}, &fs)
	// schema not found?
	if ce, confluentErr := err.(confluentError); confluentErr && ce.ErrorCode == schemaNotFound {
		return false, fs, nil
	}
	// error?
	if err != nil {
		return false, fs, err
	}
	// so we have a schema then
	return true, fs, nil
}

// GetSchemaById returns the schema for some id.
// The schema registry only provides the schema itself, not the id, subject or version.
func (c *Client) GetSchemaById(id int) (string, error) {
	var s Schema
	err := c.do("GET", fmt.Sprintf("/schemas/ids/%d", id), nil, &s)
	return s.Schema, err
}

// GetSchemaBySubject returns the schema for a particular subject and version.
func (c *Client) GetSchemaBySubject(subject string, ver int) (s Schema, err error) {
	err = c.do("GET", fmt.Sprintf("/subjects/%s/versions/%d", subject, ver), nil, &s)
	return
}

// GetLatestSchema returns the latest version of the subject's schema.
func (c *Client) GetLatestSchema(subject string) (s Schema, err error) {
	err = c.do("GET", fmt.Sprintf("/subjects/%s/versions/latest", subject), nil, &s)
	return
}

// NewClient returns a new Client that connects to baseurl.
func NewClient(baseurl string) (*Client, error) {
	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}
	return &Client{*u, http.DefaultClient}, nil
}
