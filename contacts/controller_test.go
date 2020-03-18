package contacts

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

var contactsList = []*Contact{
	&Contact{1, "Inosuke", "Hashibira"},
	&Contact{2, "Gonpachiro", "Kamaboko"},
}

type MockedContactsRepository struct{}
type MockedResponseWriter struct {
	headers http.Header
	status  int
	body    []byte
}

func (m *MockedResponseWriter) Header() http.Header {
	return m.headers
}

func (m *MockedResponseWriter) Write(body []byte) (int, error) {
	m.body = body
	return len(body), nil
}

func (m *MockedResponseWriter) WriteHeader(status int) {
	m.status = status
}

func NewMockedResponseWriter() *MockedResponseWriter {
	return &MockedResponseWriter{
		headers: make(http.Header),
	}
}

func (m *MockedContactsRepository) FindAll() ([]*Contact, error) {
	return contactsList, nil
}

func TestControllerFindAll(t *testing.T) {
	wr := NewMockedResponseWriter()
	controller := ProvideContactsController(
		&MockedContactsRepository{},
		zap.NewNop(),
	)

	controller.FindAll(wr, &http.Request{})

	var response = struct {
		Contacts []*Contact `json:"contacts"`
	}{}

	err := json.Unmarshal(wr.body, &response)

	if err != nil {
		t.Errorf("JSON unmarshal returned an error while parsing response body: %v", err)
	}

	if expected, got := "application/json", wr.Header().Get("Content-Type"); got != expected {
		t.Errorf("FindAll response header 'Content-Type' = %s, want %s", got, expected)
	}

	if expected := 200; wr.status != expected {
		t.Errorf("FindAll wrote respose status %d, want %d", expected, wr.status)
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
