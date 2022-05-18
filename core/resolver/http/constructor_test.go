package http

import (
	"errors"
	"fmt"
	"github.com/csocquet/route53-dyndns-agent/core"
	mock2 "github.com/csocquet/route53-dyndns-agent/core/resolver/http/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type OptMock struct {
	mock.Mock
}

func (m *OptMock) Apply(r *resolver) error {
	args := m.Called(r)
	return args.Error(0)
}

type ConstructorTestSuite struct {
	suite.Suite
}

func (suite *ConstructorTestSuite) TestNewResolver() {
	r, err := NewHttpResolver()
	suite.NoError(err)
	suite.NotNil(r)
	suite.Implements((*core.Resolver)(nil), r)

	suite.Equal(http.DefaultClient, r.(*resolver).httpClient)
	suite.Equal(DEFAULT_URL, r.(*resolver).url)
	suite.Equal(DEFAULT_REGEXP, r.(*resolver).regexp.String())

}

func (suite *ConstructorTestSuite) TestNewResolverOptions() {
	opt1 := new(OptMock)
	opt2 := new(OptMock)

	opt1.On("Apply", mock.Anything).Return(nil).Once()
	opt2.On("Apply", mock.Anything).Return(nil).Once()

	_, err := NewHttpResolver(opt1.Apply, opt2.Apply)
	suite.NoError(err)
	opt1.AssertExpectations(suite.T())
	opt2.AssertExpectations(suite.T())
}

func (suite *ConstructorTestSuite) TestNewResolverOptionsError() {
	expectedErr := errors.New("test option error")

	opt1 := new(OptMock)
	opt2 := new(OptMock)

	opt1.On("Apply", mock.Anything).Return(expectedErr).Once()
	opt2.On("Apply", mock.Anything).Return(nil)

	r, err := NewHttpResolver(opt1.Apply, opt2.Apply)
	suite.Nil(r)
	suite.Error(err)
	suite.ErrorContains(err, expectedErr.Error())

	opt1.AssertExpectations(suite.T())
	opt2.AssertNotCalled(suite.T(), "Apply", mock.Anything)
}

func (suite *ConstructorTestSuite) TestNewResolverWithHttpClient() {
	c := new(mock2.HttpClient)

	r, err := NewHttpResolver(WithHttpClient(c))
	suite.NoError(err)
	suite.NotNil(r)
	suite.Implements((*core.Resolver)(nil), r)

	suite.Equal(c, r.(*resolver).httpClient)
}

func (suite *ConstructorTestSuite) TestNewResolverWithUrl() {
	url := "https://my-custom-url.com"
	r, err := NewHttpResolver(WithUrl(url))
	suite.NoError(err)
	suite.NotNil(r)
	suite.Implements((*core.Resolver)(nil), r)

	suite.Equal(url, r.(*resolver).url)
}

func (suite *ConstructorTestSuite) TestNewResolverWithInvalidUrl() {
	url := "an_invalid_url"
	r, err := NewHttpResolver(WithUrl(url))
	suite.Nil(r)
	suite.Error(err)
	suite.ErrorContains(err, fmt.Sprintf("invalid url `%s`", url))
}

func (suite *ConstructorTestSuite) TestWithRegex() {
	expr := `[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`

	r, err := NewHttpResolver(WithRegexp(expr))

	suite.NoError(err)
	suite.NotNil(r)
	suite.Implements((*core.Resolver)(nil), r)
	suite.Equal(expr, r.(*resolver).regexp.String())
}

func (suite *ConstructorTestSuite) TestWithInvalidRegex() {
	expr := "[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\"

	r, err := NewHttpResolver(WithRegexp(expr))

	suite.Error(err)
	suite.ErrorContains(err, fmt.Sprintf("invalid regexp `%s`", expr))
	suite.Nil(r)
}

func TestConstructorTestSuite(t *testing.T) {
	suite.Run(t, new(ConstructorTestSuite))
}
