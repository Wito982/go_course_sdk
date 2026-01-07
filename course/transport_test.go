package course_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	courseSDK "github.com/Wito982/go_course_sdk/course"
	"github.com/Wito982/gocourse_domain/domain"
	c "github.com/Wito982/golang-restclient/rest"
)

var header = http.Header{}
var sdk courseSDK.Transport

func TestMain(m *testing.M) {
	header = http.Header{}
	header.Set("Content-Type", "application/json")
	c.StartMockupServer()
	sdk = courseSDK.NewHttpClient("base-url", "")
	os.Exit(m.Run())
}

func TestGet_Response404Error(t *testing.T) {
	expectedError := courseSDK.ErrNotFound{Message: "course 1 not found"}

	c.FlushMockups()
	err := c.AddMockups(&c.Mock{
		URL:          "base-url/courses/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 404,
		RespBody:     fmt.Sprintf(`{"message": "%s","status": 404}`, expectedError.Error()),
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	course, err := sdk.Get("1")
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error %v", err)
	}

	if course != nil {
		t.Errorf("expected nil but got course %v", course)
	}

}

func TestGet_Response500Error(t *testing.T) {
	expectedError := errors.New("internal server error")

	err := c.AddMockups(&c.Mock{
		URL:          "base-url/courses/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 500,
		RespBody:     fmt.Sprintf(`{"message": "%s","status": 500}`, "internal server error"),
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	course, err := sdk.Get("1")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected error, got nil")
	}

	if course != nil {
		t.Errorf("expected nil but got course %v", course)
	}
}

func TestGet_ResponseMarshallError(t *testing.T) {
	expectedError := errors.New("unexpected end of JSON input")

	err := c.AddMockups(&c.Mock{
		URL:          "base-url/courses/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespBody:     `{`,
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	course, err := sdk.Get("1")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}

	if course != nil {
		t.Errorf("expected nil but got course %v", course)
	}

}

func TestGet_ClientError(t *testing.T) {
	expectedError := "client error"

	err := c.AddMockups(&c.Mock{
		URL:          "base-url/courses/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 400,
		RespBody:     fmt.Sprintf(`{"message": "%s","status": 400}`, expectedError),
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	course, err := sdk.Get("1")
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}

	if course != nil {
		t.Errorf("expected nil but got course %v", course)
	}
}

func TestGet_ResponseSuccess(t *testing.T) {
	expectedCourse := &domain.Course{
		ID:   "1",
		Name: "course 1",
	}

	expectedCourseJSON, err := json.Marshal(expectedCourse)
	if err != nil {
		t.Errorf("failed to marshal expected course: %v", err)
	}

	err = c.AddMockups(&c.Mock{
		URL:          "base-url/courses/1",
		RespHeaders:  header,
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: 200,
		RespBody:     fmt.Sprintf(`{"message": "success","status": 200, "data": %s}`, expectedCourseJSON),
	})

	if err != nil {
		t.Errorf("failed to add mockup: %v", err)
	}

	course, err := sdk.Get("1")
	if err != nil {
		t.Errorf("expected nil but got error %v", err)
	}

	if course == nil {
		t.Errorf("expected course but got nil")
	}

	if course.ID != expectedCourse.ID || course.Name != expectedCourse.Name {
		t.Errorf("expected course %v, got %v", expectedCourse, course)
	}
}
