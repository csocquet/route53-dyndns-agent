package mock

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type HttpClient struct {
	mock.Mock
}

func (m *HttpClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}
