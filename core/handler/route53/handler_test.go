package route53

import (
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/csocquet/route53-dyndns-agent/core"
	"github.com/csocquet/route53-dyndns-agent/core/handler/route53/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

type HandlerTestSuite struct {
	suite.Suite
	handler   core.Handler
	apiClient *mocks.ApiClientMock
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.apiClient = new(mocks.ApiClientMock)
	suite.handler, _ = NewRoute53Handler(WithApiClient(suite.apiClient))
}

func (suite *HandlerTestSuite) TestHandle() {
	suite.apiClient.On(
		"ChangeResourceRecordSets",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(new(route53.ChangeResourceRecordSetsOutput), nil)

	err := suite.handler.Handle(new(core.ChangedIpEvent))
	suite.NoError(err)

	suite.apiClient.AssertExpectations(suite.T())
}
