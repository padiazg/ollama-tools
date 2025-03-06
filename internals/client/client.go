package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ResClientConfig struct {
	OnDecode func(io.ReadCloser, interface{}) error
}

type RestClient struct {
	client         HTTPClient
	jsonUnmarshal  func(data []byte, v any) error
	httpNewRequest func(method string, url string, body io.Reader) (*http.Request, error)
	onDecode       func(io.ReadCloser, interface{}) error
}

func New(cfg *ResClientConfig) *RestClient {
	return (&RestClient{}).New(cfg)
}

func (c *RestClient) New(cfg *ResClientConfig) *RestClient {

	if cfg == nil {
		cfg = &ResClientConfig{}
	}

	if cfg.OnDecode == nil {
		cfg.OnDecode = func(rc io.ReadCloser, v interface{}) error {
			if err := json.NewDecoder(rc).Decode(v); err != nil {
				return err
			}
			return nil
		}
	}

	c = &RestClient{
		jsonUnmarshal:  json.Unmarshal,
		httpNewRequest: http.NewRequest,
		onDecode:       cfg.OnDecode,
	}

	return c
}

func (c *RestClient) getClient() HTTPClient {
	if c.client == nil {
		c.client = &http.Client{Timeout: time.Duration(1) * time.Second}
	}

	return c.client
}

func (c *RestClient) Request(req *http.Request, v interface{}) error {
	var (
		client = c.getClient()
		res    *http.Response
		err    error
	)
	req.Header.Add("Accept", "application/josn")
	if res, err = client.Do(req); err != nil {
		return fmt.Errorf("request: calling api %v\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Get status code: %d\n", res.StatusCode)
	}

	if err = c.onDecode(res.Body, v); err != nil {
		return err
	}

	return nil
}
