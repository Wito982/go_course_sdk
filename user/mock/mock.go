package mock

import (
	"errors"

	"github.com/Wito982/gocourse_domain/domain"
)

type UserSDKMock struct {
	GetMock func(id string) (*domain.User, error)
}

func (m *UserSDKMock) Get(id string) (*domain.User, error) {
	if m.GetMock == nil {
		return nil, errors.New("not implemented")
	}
	return m.GetMock(id)
}
