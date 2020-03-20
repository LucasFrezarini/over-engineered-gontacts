// +build integration

package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/LucasFrezarini/go-contacts/contacts"
	"github.com/LucasFrezarini/go-contacts/container"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../.env.test")
	fmt.Println(os.Getenv("MYSQL_HOST"))
	m.Run()
}

func TestGetAllContacts(t *testing.T) {
	app, err := container.InitializeServer()
	if err != nil {
		t.Fatalf("GET /contacts/: error while initializing server from container: %v", err)
	}

	srv := httptest.NewServer(app.Router.BuildMux())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/contacts")
	if err != nil {
		t.Fatalf("GET /contacts/: http.Get returned an error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /contacts/: resp.StatusCode == %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if expected, got := "application/json", resp.Header.Get("Content-Type"); got != expected {
		t.Errorf("GET /contacts/: resp.Headers['Content-Type'] == %s, want %s", got, expected)
	}

	var expectedResponse = struct {
		Contacts []*contacts.Contact `json:"contacts"`
	}{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("GET /contacts/: error while readping response body: %v", err)
	}

	err = json.Unmarshal(b, &expectedResponse)
	if err != nil {
		t.Fatalf("GET /contacts/: error while unmarshaling response: %v", err)
	}

	if expected, got := 2, len(expectedResponse.Contacts); expected != got {
		t.Errorf("GET /contacts/: len(response.Contacts) == %d, want %d", got, expected)
	}
}
