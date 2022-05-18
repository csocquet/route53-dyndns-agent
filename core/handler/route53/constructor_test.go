package route53

import (
	"github.com/csocquet/route53-dyndns-agent/core"
	"github.com/csocquet/route53-dyndns-agent/core/handler/route53/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestConstructorTestSuite(t *testing.T) {
	suite.Run(t, new(ConstructorTestSuite))
}

type ConstructorTestSuite struct {
	suite.Suite
}

func (suite *ConstructorTestSuite) TestNewRoute53Handler() {
	h, err := NewRoute53Handler()
	suite.NoError(err)
	suite.NotNil(h)
	suite.Implements((*core.Handler)(nil), h)
	suite.IsType((*handler)(nil), h)
}

func (suite *ConstructorTestSuite) TestNewRoute53HandlerWithApiClient() {
	cli := new(mocks.ApiClientMock)

	h, err := NewRoute53Handler(WithApiClient(cli))
	suite.NoError(err)
	suite.NotNil(h)
	suite.Implements((*core.Handler)(nil), h)
	suite.IsType((*handler)(nil), h)

	suite.Equal(cli, h.(*handler).client)
}
