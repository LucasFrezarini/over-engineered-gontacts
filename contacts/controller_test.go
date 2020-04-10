package contacts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/LucasFrezarini/go-contacts/server/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetAllContacts(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := ProvideContactsController(
		ProvideContactsService(zap.NewNop(), &MockedContactsRepository{}, &MockedEmailRepository{}, &MockedPhoneRepository{}),
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
		t.Errorf("FindAll wrote respose status %d, want %d", rec.Code, expected)
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
		expected, actual := contactsList[i], c
		assert.Equalf(t, expected, actual, "FindAll response body's contact[%d] != expected", i)
	}
}

func TestCreateContact(t *testing.T) {
	body := map[string]string{
		"first_name": "Zenitsu",
		"last_name":  "Agatsuma",
	}

	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("error while marshaling request body: %v", err)
	}

	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := ProvideContactsController(
		nil,
		&MockedContactsRepository{},
		zap.NewNop(),
		e,
	)

	err = controller.Create(c)
	if err != nil {
		t.Errorf("controller Create() returned an error: %v", err)
	}

	if contentType := rec.Header().Get("Content-Type"); !strings.Contains(contentType, "application/json") {
		t.Errorf("Create() response header 'Content-Type' doesn't contains 'application/json' (%s)", contentType)
	}

	if expected := http.StatusCreated; rec.Code != expected {
		t.Errorf("Create() wrote respose status %d, want %d", rec.Code, expected)
	}

	var responseBody Contact
	err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
	if err != nil {
		t.Errorf("Create() error while unmarshaling response body: %v", err)
	}

	expectedBody := Contact{
		ID:        1,
		FirstName: body["first_name"],
		LastName:  body["last_name"],
	}

	if !reflect.DeepEqual(expectedBody, responseBody) {
		t.Errorf("Create() response body returned %v, want %v", responseBody, expectedBody)
	}
}

func TestCreateContactBadRequest(t *testing.T) {
	var testCases = []struct {
		testName      string
		body          map[string]string
		invalidFields []string
	}{
		{
			"missing_last_name",
			map[string]string{
				"first_name": "Zenitsu",
			},
			[]string{"LastName"},
		},
		{
			"missing_first_name",
			map[string]string{
				"last_name": "Agatsuma",
			},
			[]string{"FirstName"},
		},
		{
			"missing_all_fields",
			map[string]string{},
			[]string{"FirstName", "LastName"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			b, err := json.Marshal(tc.body)
			if err != nil {
				t.Fatalf("error while marshaling request body: %v", err)
			}

			e := echo.New()
			e.Validator = validator.NewCustomValidator()

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			controller := ProvideContactsController(
				nil,
				&MockedContactsRepository{},
				zap.NewNop(),
				e,
			)

			err = controller.Create(c)
			if err == nil {
				t.Error("controller Create() error == nil, want non-nil")
			}

			if contentType := rec.Header().Get("Content-Type"); !strings.Contains(contentType, "application/json") {
				t.Errorf("Create() response header 'Content-Type' doesn't contains 'application/json' (%s)", contentType)
			}

			if expected := http.StatusBadRequest; rec.Code != expected {
				t.Errorf("Create() wrote respose status %d, want %d\n\tRequest body sent: %s", rec.Code, expected, b)
			}

			var expectedResponseBody struct {
				Message string `json:"message"`
			}

			err = json.Unmarshal(rec.Body.Bytes(), &expectedResponseBody)
			if err != nil {
				t.Errorf("Create() error while unmarshaling response body: %v\nOriginal string content: %s", err, rec.Body.String())
			}

			for _, field := range tc.invalidFields {
				if msg := expectedResponseBody.Message; !strings.Contains(msg, field) {
					t.Errorf("Create() bad request message body doesn't contain the missing field information: %s", msg)
				}
			}
		})
	}

}
