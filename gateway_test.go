package lassie

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestGateway(t *testing.T) {
	client, err := NewWithAddr(*addr, *token)
	if err != nil {
		t.Fatal(err)
	}

	randomEUI := func() string {
		var eui [8]byte
		rand.Read(eui[:])
		return fmt.Sprintf("%02x-%02x-%02x-%02x-%02x-%02x-%02x-%02x", eui[0], eui[1], eui[2], eui[3], eui[4], eui[5], eui[6], eui[7])
	}

	gw, err := client.CreateGateway(Gateway{EUI: randomEUI(), IP: "127.0.1.1", StrictIP: true})
	if err != nil {
		t.Fatal(err)
	}
	gw.Tags["name"] = "The name of the gateway"

	updated, err := client.UpdateGateway(gw)
	if err != nil {
		t.Fatal(err)
	}
	updatedGw, err := client.Gateway(gw.EUI)
	if err != nil {
		t.Fatal(err)
	}
	if updatedGw.Tags["name"] != updated.Tags["name"] {
		t.Fatalf("Tags for gateway isn't updated: %v != %v", updated.Tags, updatedGw.Tags)
	}

	if err := client.DeleteGateway(updatedGw.EUI); err != nil {
		t.Fatal(err)
	}
}
