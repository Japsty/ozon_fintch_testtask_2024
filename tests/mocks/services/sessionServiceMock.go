package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type SessionServiceMock struct {
	mock.Mock
}

func NewSessionServiceMock() *SessionServiceMock {
	return &SessionServiceMock{}
}

func (ss *SessionServiceMock) GetTokenData(r *http.Request) (string, error) {
	args := ss.Called(r)
	return args.Get(0).(string), args.Error(1)
}
