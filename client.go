/*
Lassie provides a client for the REST API for Telenor LoRa.

All Create* and Update* methods return the created and updated entity, respectively, which may include setting fields that were not set.
*/
package lassie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	addr   string
	token  string
	client http.Client
}

func New() (*Client, error) {
	return NewWithAddr(addressTokenFromConfig(ConfigFile))
}

func NewWithAddr(addr, token string) (*Client, error) {
	c := &Client{
		addr:  addr,
		token: token,
	}
	return c, c.ping()
}

func (c *Client) ping() error {
	return c.get("/", nil)
}

func (c *Client) create(path string, x interface{}) error {
	return c.request(http.MethodPost, path, x, x)
}

func (c *Client) update(path string, x interface{}) error {
	return c.request(http.MethodPut, path, x, x)
}

func (c *Client) get(path string, x interface{}) error {
	return c.request(http.MethodGet, path, nil, x)
}

func (c *Client) delete(path string) error {
	return c.request(http.MethodDelete, path, nil, nil)
}

func (c *Client) request(method, path string, input, output interface{}) error {
	body := new(bytes.Buffer)
	if input != nil {
		if err := json.NewEncoder(body).Encode(input); err != nil {
			return err
		}
	}
	req, err := http.NewRequest(method, c.addr+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return newClientError(resp)
	}
	if output != nil {
		return json.NewDecoder(resp.Body).Decode(output)
	}
	return nil
}

type ClientError struct {
	HTTPStatusCode int
	Message        string
}

func newClientError(resp *http.Response) ClientError {
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ClientError{resp.StatusCode, err.Error()}
	}
	return ClientError{resp.StatusCode, string(buf)}
}

func (e ClientError) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.HTTPStatusCode), e.Message)
}
