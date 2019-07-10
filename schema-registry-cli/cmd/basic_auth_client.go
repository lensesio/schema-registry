package cmd

import (
	"net/http"
)

type BasicAuthCredential struct {
	Username string
	Password string
	transport http.RoundTripper
}

func (c BasicAuthCredential) GetClient() *http.Client {
	if c.transport == nil {
		c.transport = http.DefaultTransport
	}

	return &http.Client{ Transport: c }
}

func (c BasicAuthCredential) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(c.Username, c.Password)

	return c.transport.RoundTrip(r)
}
