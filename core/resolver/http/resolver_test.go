package http_test

import (
	"errors"
	"github.com/csocquet/route53-dyndns-agent/core"
	"github.com/csocquet/route53-dyndns-agent/core/resolver/http"
	"github.com/csocquet/route53-dyndns-agent/core/resolver/http/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net"
	http2 "net/http"
	"strings"
	"testing"
	"testing/iotest"
)

type ResolverTestSuite struct {
	suite.Suite
	httpClient *mock.HttpClient
	resolver   core.Resolver
}

func (suite *ResolverTestSuite) SetupTest() {
	suite.httpClient = new(mock.HttpClient)
}

func (suite *ResolverTestSuite) TestResolve() {
	type TestCase struct {
		Name             string
		Opts             []http.ResolverOpt
		ClientUrl        string
		ClientResponse   *http2.Response
		ClientError      error
		ExpectedIp       net.IP
		ExpectedErrorStr string
	}

	testCases := []TestCase{
		{
			Name:      "Simple 127.0.0.1 response",
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("127.0.0.1")),
			},
			ExpectedIp:       net.IPv4(127, 0, 0, 1),
			ExpectedErrorStr: "",
		},
		{
			Name: "Simple 127.0.0.1 response with different URL",
			Opts: []http.ResolverOpt{
				http.WithUrl("https://www.my-custom-url.com"),
			},
			ClientUrl: "https://www.my-custom-url.com",
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("127.0.0.1")),
			},
			ExpectedIp:       net.IPv4(127, 0, 0, 1),
			ExpectedErrorStr: "",
		},
		{
			Name:      "Simple 192.168.0.1 response",
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("192.168.0.1")),
			},
			ExpectedIp:       net.IPv4(192, 168, 0, 1),
			ExpectedErrorStr: "",
		},
		{
			Name:             "Http client error",
			ClientUrl:        http.DEFAULT_URL,
			ClientResponse:   nil,
			ClientError:      errors.New("test http client error"),
			ExpectedIp:       net.IPv4zero,
			ExpectedErrorStr: "test http client error",
		},
		{
			Name:      "403 response",
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusForbidden,
				Status:     http2.StatusText(http2.StatusForbidden),
				Body:       ioutil.NopCloser(strings.NewReader("")),
			},
			ExpectedIp:       net.IPv4zero,
			ExpectedErrorStr: "invalid http client response: 403 - Forbidden",
		},
		{
			Name:      "404 response",
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusNotFound,
				Status:     http2.StatusText(http2.StatusNotFound),
				Body:       ioutil.NopCloser(strings.NewReader("")),
			},
			ExpectedIp:       net.IPv4zero,
			ExpectedErrorStr: "invalid http client response: 404 - Not Found",
		},
		{
			Name:      "500 response",
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusInternalServerError,
				Status:     http2.StatusText(http2.StatusInternalServerError),
				Body:       ioutil.NopCloser(strings.NewReader("")),
			},
			ExpectedIp:       net.IPv4zero,
			ExpectedErrorStr: "invalid http client response: 500 - Internal Server Error",
		},
		{
			Name:      "503 response",
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusServiceUnavailable,
				Status:     http2.StatusText(http2.StatusServiceUnavailable),
				Body:       ioutil.NopCloser(strings.NewReader("")),
			},
			ExpectedIp:       net.IPv4zero,
			ExpectedErrorStr: "invalid http client response: 503 - Service Unavailable",
		},
		{
			Name:      "Body reader error",
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(iotest.ErrReader(errors.New("test body reader error"))),
			},
			ClientError:      nil,
			ExpectedIp:       net.IPv4zero,
			ExpectedErrorStr: "test body reader error",
		},
		{
			Name:      "Parse invalid IP",
			ClientUrl: http.DEFAULT_URL,
			Opts: []http.ResolverOpt{
				http.WithRegexp(`(?P<ip>.+)`),
			},
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("an_invalid_ipv4")),
			},
			ExpectedIp:       net.IPv4zero,
			ExpectedErrorStr: "failed to parse IPv4 `an_invalid_ipv4`",
		},
		{
			Name:      "Parse IPv6",
			ClientUrl: http.DEFAULT_URL,
			Opts: []http.ResolverOpt{
				http.WithRegexp(`(?P<ip>.+)`),
			},
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("2345:0425:2CA1:0000:0000:0567:5673:23b5")),
			},
			ExpectedIp:       net.IPv4zero,
			ExpectedErrorStr: "failed to parse IPv4 `2345:0425:2CA1:0000:0000:0567:5673:23b5`",
		},
		{
			Name:      "Html 127.0.0.1 response",
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My IP</title>
</head>
<body>
    <div>
        <p>Mon IP: <span>127.0.0.1</span></p>
    </div>
</body>
</html>`)),
			},
			ExpectedIp:       net.IPv4(127, 0, 0, 1),
			ExpectedErrorStr: "",
		},
		{
			Name: "Custom regexp with ip part groups",
			Opts: []http.ResolverOpt{
				http.WithRegexp(`Second IP: (?P<ip1>\d{1,3})\.(?P<ip2>\d{1,3})\.(?P<ip3>\d{1,3}).(?P<ip4>\d{1,3})`),
			},
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("First IP: 127.0.0.1\nSecond IP: 10.20.0.1")),
			},
			ExpectedIp:       net.IPv4(10, 20, 0, 1),
			ExpectedErrorStr: "",
		},
		{
			Name: "Custom regexp with ip group",
			Opts: []http.ResolverOpt{
				http.WithRegexp(`Second IP: (?P<ip>\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3})`),
			},
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("First IP: 127.0.0.1\nSecond IP: 10.20.0.1")),
			},
			ExpectedIp:       net.IPv4(10, 20, 0, 1),
			ExpectedErrorStr: "",
		},
		{
			Name: "Custom regexp without ip group",
			Opts: []http.ResolverOpt{
				http.WithRegexp(`\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}`),
			},
			ClientUrl: http.DEFAULT_URL,
			ClientResponse: &http2.Response{
				StatusCode: http2.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("10.20.0.1")),
			},
			ExpectedIp:       net.IPv4(10, 20, 0, 1),
			ExpectedErrorStr: "",
		},
	}

	for _, testCase := range testCases {

		suite.SetupTest()
		suite.Run(testCase.Name, func() {
			opts := append(testCase.Opts, http.WithHttpClient(suite.httpClient))
			resolver, _ := http.NewHttpResolver(opts...)
			suite.httpClient.On("Get", testCase.ClientUrl).Return(testCase.ClientResponse, testCase.ClientError).Once()

			ip, err := resolver.Resolve()

			if testCase.ExpectedErrorStr != "" {
				suite.Error(err)
				suite.ErrorContains(err, testCase.ExpectedErrorStr)
			} else {
				suite.NoError(err)
			}
			suite.Equal(testCase.ExpectedIp.String(), ip.String())
			suite.httpClient.AssertExpectations(suite.T())

		})
	}
}

func TestResolverTestSuite(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}
