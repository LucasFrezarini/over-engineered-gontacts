package contacts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var contactsList = []*Contact{
	&Contact{1, "Inosuke", "Hashibira"},
	&Contact{2, "Gonpachiro", "Kamaboko"},
}

type MockedContactsRepository struct{}

func (m *MockedContactsRepository) FindAll() ([]*Contact, error) {
	return contactsList, nil
}

func TestGetAllContacts(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := ProvideContactsController(
		&MockedContactsRepository{},
		zap.NewNop(),
		e,
	)

	err := controller.FindAll(c)

	if err != nil {
		t.Errorf("controller FindAll() returned an error: %v", err)
	}

	if contentType := rec.Header().Get("Content-Type"); !strings.Contains(contentType, "application/json") {
		t.Errorf("FindAll response header 'Content-Type' doesn't contains 'application/json' (%s)", contentType)
	}

	if expected := 200; rec.Code != expected {
		t.Errorf("FindAll wrote respose status %d, want %d", expected, rec.Code)
	}

	var response = struct {
		Contacts []*Contact `json:"contacts"`
	}{}

	err = json.Unmarshal(rec.Body.Bytes(), &response)

	if err != nil {
		t.Errorf("controller FindAll() error while unmarshaling response body: %v", err)
	}

	if expected, got := 2, len(response.Contacts); got != expected {
		t.Errorf("FindAll wrote %d contacts in response body, want %d", got, expected)
	}

	for i, c := range response.Contacts {
		if expected, got := contactsList[i], c; !reflect.DeepEqual(expected, got) {
			t.Errorf("FindAll response body's contact[%d] = %v, want %v", i, got, expected)
		}
	}
}
