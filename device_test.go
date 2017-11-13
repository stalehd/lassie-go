package lassie

import (
	"net/http"
	"testing"
)

func TestDevice(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	app, err := client.CreateApplication(Application{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteApplication(app.EUI)

	dev, err := client.CreateDevice(app.EUI, Device{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteDevice(app.EUI, dev.EUI)

	tagKey := "test key"
	tagValue := "test value"
	dev.Tags[tagKey] = tagValue
	dev, err = client.UpdateDevice(app.EUI, dev)
	if err != nil {
		t.Fatal(err)
	}
	if len(dev.Tags) != 1 || dev.Tags[tagKey] != tagValue {
		t.Fatal("unexpected tags:", dev.Tags)
	}

	devs, err := client.Devices(app.EUI)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range devs {
		if d.EUI == dev.EUI {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("device %v not found in %v", dev, devs)
	}

	if _, err := client.Device(app.EUI, dev.EUI); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteDevice(app.EUI, dev.EUI); err != nil {
		t.Fatal(err)
	}
	err = client.DeleteDevice(app.EUI, dev.EUI)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusNotFound {
		t.Fatal(err)
	}
}
