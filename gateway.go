package lassie

import "fmt"

// Gateway represents a gateway
type Gateway struct {
	EUI       string            `json:"gatewayEUI"`
	IP        string            `json:"ip"`
	StrictIP  bool              `json:"strictIP"`
	Latitude  float64           `json:"latitude"`
	Longitude float64           `json:"longitude"`
	Altitude  float64           `json:"altitude"`
	Tags      map[string]string `json:"tags"`
}

// CreateGateway creates a new gateway
func (c *Client) CreateGateway(gw Gateway) (Gateway, error) {
	err := c.create("/gateways", &gw)
	return gw, err
}

// Gateway retrieves a gateway
func (c *Client) Gateway(eui string) (Gateway, error) {
	var gw Gateway
	err := c.get(fmt.Sprintf("/gateways/%s", eui), &gw)
	return gw, err
}

// DeleteGateway removes a gateway
func (c *Client) DeleteGateway(eui string) error {
	return c.delete(fmt.Sprintf("/gateways/%s", eui))
}

// UpdateGateway updates the gateway
func (c *Client) UpdateGateway(gw Gateway) (Gateway, error) {
	err := c.update(fmt.Sprintf("/gateways/%s", gw.EUI), &gw)
	return gw, err
}
