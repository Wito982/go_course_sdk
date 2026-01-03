package course

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Wito982/gocourse_domain/domain"
	c "github.com/mercadolibre/golang-restclient/rest"
)

type (
	DataResponse struct {
		Message string      `json:"message"`
		Code    int         `json:"code"`
		Data    interface{} `json:"data"`
		Meta    interface{} `json:"meta"`
	}

	Transport interface {
		Get(id string) (*domain.Course, error)
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

func (c *clientHTTP) Get(id string) (*domain.Course, error) {
	dataResponse := &DataResponse{Data: &domain.Course{}}

	u := url.URL{}
	u.Path += fmt.Sprintf("/courses/%s", id)
	reps := c.client.Get(u.String())
	if reps.Err != nil {
		return nil, reps.Err
	}

	if reps.StatusCode == 404 {
		return nil, ErrNotFound{fmt.Sprintf("%s", reps)}
	}

	if reps.StatusCode > 299 {
		return nil, fmt.Errorf("%s", reps)
	}

	if err := reps.FillUp(&dataResponse); err != nil {
		return nil, err
	}

	return dataResponse.Data.(*domain.Course), nil
}
