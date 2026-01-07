package user

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Wito982/gocourse_domain/domain"
	c "github.com/Wito982/golang-restclient/rest"
)

type (
	DataResponse struct {
		Message string      `json:"message"`
		Code    int         `json:"code"`
		Data    interface{} `json:"data"`
		Meta    interface{} `json:"meta"`
	}

	Transport interface {
		Get(id string) (*domain.User, error)
	}

	clientHTTP struct {
		client *c.RequestBuilder
	}
)

func NewHttpClient(baseUrl, token string) Transport {
	header := http.Header{}

	if token != "" {
		header.Set("Authorization", token)
	}

	customClient := &c.RequestBuilder{
		BaseURL: baseUrl,
		Timeout: 5000 * time.Millisecond,
		Headers: header,
	}

	return &clientHTTP{
		client: customClient,
	}
}

func (c *clientHTTP) Get(id string) (*domain.User, error) {
	dataResponse := &DataResponse{Data: &domain.User{}}

	u := url.URL{}
	u.Path += fmt.Sprintf("/users/%s", id)
	reps := c.client.Get(u.String())
	if reps.Err != nil {
		return nil, reps.Err
	}

	if err := reps.FillUp(&dataResponse); err != nil {
		return nil, err
	}

	if reps.StatusCode == 404 {
		return nil, ErrNotFound{fmt.Sprintf("%s", dataResponse.Message)}
	}

	if reps.StatusCode > 299 {
		return nil, fmt.Errorf("%s", dataResponse.Message)
	}

	return dataResponse.Data.(*domain.User), nil

}
