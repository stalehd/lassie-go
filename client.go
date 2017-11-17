/*
Package lassie provides a client for the REST API for Telenor LoRa.

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

// Client is a client for the Congress service.
type Client struct {
	addr   string
	token  string
	client http.Client
}

// New creates a new client with the default configuration. The default
// configuration can be specified in a configuration file or through
// environment variables.
func New() (*Client, error) {
	return NewWithAddr(addressTokenFromConfig(ConfigFile))
}

// NewWithAddr creates a new client with the specified address and token.
func NewWithAddr(addr, token string) (*Client, error) {
	c := &Client{
		addr:  addr,
		token: token,
	}
	return c, c.ping()
}

// Address returns the client's address.
func (c *Client) Address() string {
	return c.addr
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

// ClientError is the errors returned by the Client.
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
