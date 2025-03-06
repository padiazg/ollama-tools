package client

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CheckClientFn func(t *testing.T, c *RestClient)
type CheckHttpClientFn func(t *testing.T, c HTTPClient)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

var (
	checkClient     = func(fns ...CheckClientFn) []CheckClientFn { return fns }
	checkHttpClient = func(fns ...CheckHttpClientFn) []CheckHttpClientFn { return fns }
)

func TestRestClient_New(t *testing.T) {
	var (
		// checkUrlSet = func() CheckClientFn {
		// 	return func(t *testing.T, c *RestClient) {
		// 		t.Helper()
		// 		assert.NotEmpty(t, c.url, "url cannot be empty")
		// 	}
		// }

		checJsonUnmarshal = func() CheckClientFn {
			return func(t *testing.T, c *RestClient) {
				t.Helper()
				assert.NotNil(t, c.jsonUnmarshal, "jsonUnmarshal cannot not be empty")
			}
		}

		checHttpNewRequest = func() CheckClientFn {
			return func(t *testing.T, c *RestClient) {
				t.Helper()
				assert.NotNil(t, c.httpNewRequest, "httpNewRequest cannot not be empty")
			}
		}

		tests = []struct {
			name   string
			url    string
			checks []CheckClientFn
		}{
			// {
			// 	name:   "empty-url",
			// 	checks: checkClient(checkUrlSet()),
			// },
			{
				name: "full-setup",
				checks: checkClient(
					// checkUrlSet(),
					checJsonUnmarshal(),
					checHttpNewRequest(),
				),
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			n := New(&ResClientConfig{})
			for _, c := range tt.checks {
				c(t, n)
			}
		})
	}
}

func TestRestClient_getClient(t *testing.T) {
	var (
		checkRegularClient = func() CheckHttpClientFn {
			return func(t *testing.T, c HTTPClient) {
				t.Helper()
				assert.IsTypef(t, &http.Client{}, c, "checkRegularClient = %T, expected http.Client", c)
			}
		}

		checkMockClient = func() CheckHttpClientFn {
			return func(t *testing.T, c HTTPClient) {
				t.Helper()
				assert.IsTypef(t, &mockHTTPClient{}, c, "checkRegularClient = %T, expected mockHTTPClient", c)
			}
		}

		tests = []struct {
			name         string
			before       func(o *RestClient)
			ollamaClient *RestClient
			checks       []CheckHttpClientFn
		}{
			{
				name:         "assign-on-empty-client",
				ollamaClient: New(&ResClientConfig{}),
				checks: checkHttpClient(
					checkRegularClient(),
				),
			},
			{
				name:         "regular-client",
				ollamaClient: New(&ResClientConfig{}),
				before: func(o *RestClient) {
					o.client = &http.Client{}
				},
				checks: checkHttpClient(
					checkRegularClient(),
				),
			},
			{
				name:         "mock-client",
				ollamaClient: New(&ResClientConfig{}),
				before: func(o *RestClient) {
					o.client = &mockHTTPClient{
						DoFunc: func(req *http.Request) (*http.Response, error) {
							return nil, errors.New("test http new request error")
						},
					}
				},
				checks: checkHttpClient(
					checkMockClient(),
				),
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(tt.ollamaClient)
			}
			httpClient := tt.ollamaClient.getClient()

			for _, c := range tt.checks {
				c(t, httpClient)
			}
		})
	}
}
