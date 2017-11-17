package lassie

import (
	"encoding/hex"
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"
)

// Application represents an application.
type Application struct {
	EUI  string            `json:"applicationEUI"`
	Tags map[string]string `json:"tags"`
}

// CreateApplication creates an application.
func (c *Client) CreateApplication(app Application) (Application, error) {
	err := c.create("/applications", &app)
	return app, err
}

// UpdateApplication updates an application.
func (c *Client) UpdateApplication(app Application) (Application, error) {
	err := c.update("/applications/"+app.EUI, &app)
	return app, err
}

// Applications gets all applications.
func (c *Client) Applications() ([]Application, error) {
	var apps struct {
		Apps []Application `json:"applications"`
	}
	err := c.get("/applications", &apps)
	return apps.Apps, err
}

// Application gets an application.
func (c *Client) Application(eui string) (Application, error) {
	var app Application
	err := c.get("/applications/"+eui, &app)
	return app, err
}

// DeleteApplication deletes an application.
func (c *Client) DeleteApplication(eui string) error {
	return c.delete("/applications/" + eui)
}

// ApplicationStream calls handler in its own goroutine for each device message received.
// It blocks until an error occurs, including EOF when the connection is closed.
func (c *Client) ApplicationStream(appeui string, handler func(DeviceData)) error {
	url, err := url.Parse(c.addr)
	if err != nil {
		return err
	}

	scheme := "wss"
	if url.Scheme == "http" {
		scheme = "ws"
	}

	wscfg, err := websocket.NewConfig(fmt.Sprintf("%s://%s/applications/%s/stream", scheme, url.Host, appeui), "http://example.com")
	if err != nil {
		return err
	}
	wscfg.Header.Set("X-API-Token", c.token)

	ws, err := websocket.DialConfig(wscfg)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		var msg struct {
			Type  string     `json:"type"`
			Error string     `json:"message"`
			Data  DeviceData `json:"data"`
		}
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			return err
		}

		switch msg.Type {
		case "Error":
			return fmt.Errorf(msg.Error)
		case "DeviceData":
			msg.Data.Data, _ = hex.DecodeString(msg.Data.HexData)
			go handler(msg.Data)
		default:
			// Ignore it
		}
	}
}

// DeviceData represents data received from a device.
type DeviceData struct {
	DeviceEUI      string `json:"deviceEUI"`
	DeviceAddress  string `json:"devAddr"`
	GatewayEUI     string `json:"gatewayEUI"`
	ApplicationEUI string `json:"appEUI"`
	Timestamp      int64  `json:"timestamp"`
	HexData        string `json:"data"`
	Data           []byte
	Frequency      float32 `json:"frequency"`
	DataRate       string  `json:"dataRate"`
	RSSI           int32   `json:"rssi"`
	SNR            float32 `json:"snr"`
}
