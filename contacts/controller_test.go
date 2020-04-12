package contacts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LucasFrezarini/go-contacts/contacts/email"
	"github.com/LucasFrezarini/go-contacts/contacts/phone"
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

func TestCreateContactSuccess(t *testing.T) {
	var testCases = []struct {
		body map[string]interface{}
	}{
		{
			body: map[string]interface{}{
				"first_name": "Giyu",
				"last_name":  "Tomioka",
			},
		},
		{
			body: map[string]interface{}{
				"first_name": "Kyojuro",
				"last_name":  "Rengoku",
				"emails":     []string{},
			},
		},
		{
			body: map[string]interface{}{
				"first_name": "Zenitsu",
				"last_name":  "Agatsuma",
				"emails":     []string{"zenitsu.agatsuma@gmail.com"},
				"phones": []map[string]string{
					map[string]string{
						"type":   "home",
						"number": "551122223333",
					},
					map[string]string{
						"type":   "mobile",
						"number": "5511944445555",
					},
					map[string]string{
						"type":   "fax",
						"number": "5511666677777",
					},
					map[string]string{
						"type":   "work",
						"number": "551188889999",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		body := tc.body
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
			ProvideContactMockedService(),
			&MockedContactsRepository{},
			zap.NewNop(),
			e,
		)

		err = controller.Create(c)
		if err != nil {
			t.Errorf("controller Create() returned an error: %v\nRequest body:\t%v", err, body)
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

		var expectedEmails []email.Email
		var expectedPhones []phone.Phone

		if emails, exists := body["emails"]; exists {
			for _, e := range emails.([]string) {
				expectedEmails = append(expectedEmails, email.Email{
					Address: e,
				})
			}
		}

		if phones, exists := body["phones"]; exists {
			for _, p := range phones.([]map[string]string) {
				expectedPhones = append(expectedPhones, phone.Phone{
					Type:   p["type"],
					Number: p["number"],
				})
			}
		}

		expectedBody := Contact{
			ID:        1,
			FirstName: body["first_name"].(string),
			LastName:  body["last_name"].(string),
			Emails:    expectedEmails,
			Phones:    expectedPhones,
		}

		if expected, got := expectedBody.FirstName, responseBody.FirstName; expected != got {
			t.Errorf("Create() response body first_name returned %q, want %q", got, expected)
		}

		if expected, got := expectedBody.LastName, responseBody.LastName; expected != got {
			t.Errorf("Create() response body last_name returned %q, want %q", got, expected)
		}

		if expected, got := len(expectedBody.Emails), len(responseBody.Emails); expected != got {
			t.Errorf("Create() response body returned %d emails, want %d", got, expected)
		}

		if expected, got := len(expectedBody.Phones), len(responseBody.Phones); expected != got {
			t.Errorf("Create() response body returned %d phones, want %d", got, expected)
		}
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
