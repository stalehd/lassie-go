package lassie

import (
	"net/http"
	"testing"
)

func TestApplication(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	app, err := client.CreateApplication(Application{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteApplication(app.EUI)

	tagKey := "test key"
	tagValue := "test value"
	app.Tags[tagKey] = tagValue
	app, err = client.UpdateApplication(app)
	if err != nil {
		t.Fatal(err)
	}
	if len(app.Tags) != 1 || app.Tags[tagKey] != tagValue {
		t.Fatal("unexpected tags:", app.Tags)
	}

	apps, err := client.Applications()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, a := range apps {
		if a.EUI == app.EUI {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("application %v not found in %v", app, apps)
	}

	if _, err := client.Application(app.EUI); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteApplication(app.EUI); err != nil {
		t.Fatal(err)
	}
	err = client.DeleteApplication(app.EUI)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusNotFound {
		t.Fatal(err)
	}
}
