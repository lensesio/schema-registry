package schemaregistry

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

var DefaultUrl = "http://localhost:8081"

const (
	SubjectNotFound = 40401
	SchemaNotFound  = 40403
)

type Schema struct {
	Schema string `json:"schema"`
}

type SubjectSchema struct {
	VersionedSchema
	Subject string `json:"subject"`
	Id      int    `json:"id"`
}

type NameSchema struct {
	VersionedSchema
	Name string `json:"name"`
}

type VersionedSchema struct {
	Schema
	Version int `json:"version"`
}

type ConfluentError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func (ce ConfluentError) Error() string {
	return fmt.Sprintf("%s (%d)", ce.Message, ce.ErrorCode)
}

type Client struct {
	url url.URL
}

func confluentErrorToError(resp *http.Response) error {
	var ce ConfluentError
	if err := json.NewDecoder(resp.Body).Decode(&ce); err != nil {
		return err
	}
	return ce
}

func (c *Client) do(m, p string, in interface{}, out interface{}) error {
	u := c.url
	u.Path = path.Join(u.Path, p)
	log.Printf("%s %s\n", m, u.String())
	var rdp io.Reader
	if in != nil {
		var wr *io.PipeWriter
		rdp, wr = io.Pipe()
		go func() {
			wr.CloseWithError(json.NewEncoder(wr).Encode(in))
		}()
	}
	req, err := http.NewRequest(m, u.String(), rdp)
	req.Header.Add("Accept", "application/vnd.schemaregistry.v1+json, application/vnd.schemaregistry+json, application/json")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return confluentErrorToError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) Subjects() ([]string, error) {
	var ss []string
	err := c.do("GET", "subjects", nil, &ss)
	return ss, err
}

func (c *Client) Versions(sub string) ([]int, error) {
	var vs []int
	err := c.do("GET", fmt.Sprintf("subjects/%s/versions", sub), nil, &vs)
	return vs, err
}

func (c *Client) RegisterNewSchema(sub, schema string) (int, error) {
	var resp struct {
		Id int `json:"id"`
	}
	//	s := Schema{schema}
	err := c.do("POST", fmt.Sprintf("/subjects/%s/versions", sub), Schema{schema}, &resp)
	return resp.Id, err
}

func (c *Client) IsRegistered(sub, schema string) (bool, SubjectSchema, error) {
	var fs SubjectSchema
	err := c.do("POST", fmt.Sprintf("/subjects/%s", sub), Schema{schema}, &fs)
	// schema not found?
	if ce, confluentErr := err.(ConfluentError); confluentErr && ce.ErrorCode == SchemaNotFound {
		return false, fs, nil
	}
	// error?
	if err != nil {
		return false, fs, err
	}
	// so we have a schema then
	return true, fs, nil
}

func (c *Client) GetSchemaById(id int) (Schema, error) {
	var s Schema
	err := c.do("GET", fmt.Sprintf("/schemas/ids/%d", id), nil, &s)
	return s, err
}

func (c *Client) GetSchemaBySubjectVersion(sub string, ver int) (NameSchema, error) {
	var s NameSchema
	err := c.do("GET", fmt.Sprintf("/subjects/%s/versions/%d", sub, ver), nil, &s)
	return s, err
}

func (c *Client) GetLatestSchema(sub string) (NameSchema, error) {
	var s NameSchema
	err := c.do("GET", fmt.Sprintf("/subjects/%s/versions/latest", sub), nil, &s)
	return s, err
}

func NewClient(baseurl string) (*Client, error) {
	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}
	return &Client{*u}, nil
}
