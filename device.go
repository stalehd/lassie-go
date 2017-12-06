package lassie

import "fmt"

// Device represents a device.
type Device struct {
	EUI                   string            `json:"deviceEUI"`
	Type                  string            `json:"deviceType"`
	DeviceAddress         string            `json:"devAddr"`
	ApplicationKey        string            `json:"appKey"`
	ApplicationSessionKey string            `json:"appSKey"`
	NetworkSessionKey     string            `json:"nwkSKey"`
	FrameCounterUp        uint16            `json:"fCntUp"`
	FrameCounterDown      uint16            `json:"fCntDn"`
	RelaxedCounter        bool              `json:"relaxedCounter"`
	KeyWarning            bool              `json:"keyWarning"`
	Tags                  map[string]string `json:"tags"`
}

// DownstreamMessage represents downstream messages to devices.
type DownstreamMessage struct {
	Port       uint8  `json:"port"`
	Ack        bool   `json:"ack"`
	HexPayload string `json:"data"`
}

// CreateDevice creates a device.
func (c *Client) CreateDevice(appeui string, dev Device) (Device, error) {
	err := c.create(fmt.Sprintf("/applications/%s/devices", appeui), &dev)
	return dev, err
}

// UpdateDevice updates a device.
func (c *Client) UpdateDevice(appeui string, dev Device) (Device, error) {
	err := c.update(fmt.Sprintf("/applications/%s/devices/%s", appeui, dev.EUI), &dev)
	return dev, err
}

// Devices gets all devices.
func (c *Client) Devices(appeui string) ([]Device, error) {
	var devs struct {
		Devs []Device `json:"devices"`
	}
	err := c.get(fmt.Sprintf("/applications/%s/devices", appeui), &devs)
	return devs.Devs, err
}

// Device gets a device.
func (c *Client) Device(appeui, deveui string) (Device, error) {
	var dev Device
	err := c.get(fmt.Sprintf("/applications/%s/devices/%s", appeui, deveui), &dev)
	return dev, err
}

// DeleteDevice deletes a device.
func (c *Client) DeleteDevice(appeui, deveui string) error {
	return c.delete(fmt.Sprintf("/applications/%s/devices/%s", appeui, deveui))
}

// ScheduleMessage schedules a downstream message to a device. Port is in range 1-223.
// If the ack flag is set the message will be re-sent until the device acknowledges it.
// The payload is hex encoded.
func (c *Client) ScheduleMessage(appeui, deveui string, msg DownstreamMessage) error {
	return c.create(fmt.Sprintf("/applications/%s/devices/%s/message", appeui, deveui), &msg)
}

// GetScheduledMessage returns the scheduled message for the device.
func (c *Client) GetScheduledMessage(appeui, deveui string) (DownstreamMessage, error) {
	var msg DownstreamMessage
	err := c.get(fmt.Sprintf("/applications/%s/devices/%s/message", appeui, deveui), &msg)
	return msg, err
}
