package user_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	userSDK "github.com/Wito982/go_course_sdk/user"
	"github.com/Wito982/gocourse_domain/domain"
	c "github.com/Wito982/golang-restclient/rest"
)

var header = http.Header{}
var sdk userSDK.Transport

func TestMain(m *testing.M) {
	header = http.Header{}
	header.Set("Content-Type", "application/json")
	c.StartMockupServer()
	sdk = userSDK.NewHttpClient("base-url", "")
	os.Exit(m.Run())
}

func TestGet_Response404Error(t *testing.T) {
	expectedError := userSDK.ErrNotFound{Message: "user 1 not found"}

	c.FlushMockups()
	err := c.AddMockups(&c.Mock{
		URL:          "base-url/users/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 404,
		RespBody:     fmt.Sprintf(`{"message": "%s","status": 404}`, expectedError.Error()),
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	user, err := sdk.Get("1")
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error %v", err)
	}

	if user != nil {
		t.Errorf("expected nil but got user %v", user)
	}
}

func TestGet_Response500Error(t *testing.T) {
	expectedError := errors.New("internal server error")

	c.FlushMockups()
	err := c.AddMockups(&c.Mock{
		URL:          "base-url/users/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 500,
		RespBody:     fmt.Sprintf(`{"message": "%s","status": 500}`, "internal server error"),
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	user, err := sdk.Get("1")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected error, got nil")
	}

	if user != nil {
		t.Errorf("expected nil but got user %v", user)
	}
}

func TestGet_ResponseMarshallError(t *testing.T) {
	expectedError := errors.New("unexpected end of JSON input")

	c.FlushMockups()
	err := c.AddMockups(&c.Mock{
		URL:          "base-url/users/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespBody:     `{`,
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	user, err := sdk.Get("1")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}

	if user != nil {
		t.Errorf("expected nil but got user %v", user)
	}
}

func TestGet_ClientError(t *testing.T) {
	expectedError := "client error"

	c.FlushMockups()
	err := c.AddMockups(&c.Mock{
		URL:          "base-url/users/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 400,
		RespBody:     fmt.Sprintf(`{"message": "%s","status": 400}`, expectedError),
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	user, err := sdk.Get("1")
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}

	if user != nil {
		t.Errorf("expected nil but got user %v", user)
	}
}

func TestGet_ResponseSuccess(t *testing.T) {
	expectedUser := &domain.User{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "123456789",
	}

	expectedUserJSON, err := json.Marshal(expectedUser)
	if err != nil {
		t.Errorf("failed to marshal expected user: %v", err)
	}

	c.FlushMockups()
	err = c.AddMockups(&c.Mock{
		URL:          "base-url/users/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespBody:     fmt.Sprintf(`{"message": "success","status": 200, "data": %s}`, expectedUserJSON),
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	user, err := sdk.Get("1")
	if err != nil {
		t.Errorf("expected nil but got error %v", err)
	}

	if user == nil {
		t.Errorf("expected user but got nil")
		return
	}

	if user.ID != expectedUser.ID || user.FirstName != expectedUser.FirstName || user.LastName != expectedUser.LastName {
		t.Errorf("expected user %v, got %v", expectedUser, user)
	}
}
