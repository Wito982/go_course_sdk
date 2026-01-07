package mock

import (
	"errors"

	"github.com/Wito982/gocourse_domain/domain"
)

type CourseSDKMock struct {
	GetMock func(id string) (*domain.Course, error)
}

func (m *CourseSDKMock) Get(id string) (*domain.Course, error) {
	if m.GetMock == nil {
		return nil, errors.New("not implemented")
	}
	return m.GetMock(id)
}
