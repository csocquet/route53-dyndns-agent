package core_test

import (
	"errors"
	"github.com/csocquet/route53-dyndns-agent/core"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net"
	"testing"
)

type ResolverMock struct {
	mock.Mock
}

func (m *ResolverMock) Resolve() (net.IP, error) {
	args := m.Called()
	return args.Get(0).(net.IP), args.Error(1)
}

type HandlerMock struct {
	mock.Mock
}

func (m *HandlerMock) Handle(e *core.ChangedIpEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

type AgentTestSuite struct {
	suite.Suite
	resolver *ResolverMock
	handler  *HandlerMock
	agent    core.Agent
}

func (suite *AgentTestSuite) SetupTest() {
	suite.resolver = new(ResolverMock)
	suite.handler = new(HandlerMock)
	suite.agent, _ = core.NewAgent(suite.resolver, suite.handler)
}

func (suite *AgentTestSuite) TestNewAgent() {
	agent, err := core.NewAgent(suite.resolver, suite.handler)

	suite.NoError(err)
	suite.Implements((*core.Agent)(nil), agent)
}

func (suite *AgentTestSuite) TestRun() {
	ip := net.IPv4(127, 0, 0, 1)
	suite.resolver.On("Resolve").Return(ip, nil)
	suite.handler.On("Handle", core.NewChangedIpEvent(ip, net.IPv4zero)).Return(nil)

	agent, _ := core.NewAgent(suite.resolver, suite.handler)
	err := agent.Run()

	suite.NoError(err)
	suite.resolver.AssertExpectations(suite.T())
	suite.handler.AssertExpectations(suite.T())
}

func (suite *AgentTestSuite) TestRunResolverError() {
	expectedErr := errors.New("test resolver error")
	suite.resolver.On("Resolve").Return(net.IPv4zero, expectedErr)
	suite.handler.On("Handle", mock.Anything).Return(nil)

	err := suite.agent.Run()

	suite.Error(err)
	suite.ErrorContains(err, expectedErr.Error())
	suite.handler.AssertNotCalled(suite.T(), "Handle", mock.Anything)

}

func (suite *AgentTestSuite) TestRunHandlerError() {
	expectedErr := errors.New("test handler error")
	suite.resolver.On("Resolve").Return(net.IPv4(127, 0, 0, 1), nil)
	suite.handler.On("Handle", mock.Anything).Return(expectedErr)

	err := suite.agent.Run()

	suite.Error(err)
	suite.ErrorContains(err, expectedErr.Error())
}

func (suite *AgentTestSuite) TestRunTwice() {
	ip1 := net.IPv4(127, 0, 0, 1)
	ip2 := net.IPv4(192, 168, 0, 1)

	suite.resolver.On("Resolve").Return(ip1, nil).Once()
	suite.resolver.On("Resolve").Return(ip2, nil).Once()
	suite.handler.On("Handle", mock.Anything).Return(nil)

	_ = suite.agent.Run()
	suite.handler.AssertNumberOfCalls(suite.T(), "Handle", 1)
	suite.handler.AssertCalled(suite.T(), "Handle", core.NewChangedIpEvent(ip1, net.IPv4zero))

	_ = suite.agent.Run()
	suite.handler.AssertNumberOfCalls(suite.T(), "Handle", 2)
	suite.handler.AssertCalled(suite.T(), "Handle", core.NewChangedIpEvent(ip2, ip1))
}

func (suite *AgentTestSuite) TestRunTwiceWithSameIp() {
	ip := net.IPv4(127, 0, 0, 1)

	suite.resolver.On("Resolve").Return(ip, nil)
	suite.handler.On("Handle", mock.Anything).Return(nil)

	_ = suite.agent.Run()
	suite.handler.AssertNumberOfCalls(suite.T(), "Handle", 1)
	suite.handler.AssertCalled(suite.T(), "Handle", core.NewChangedIpEvent(ip, net.IPv4zero))

	_ = suite.agent.Run()
	suite.handler.AssertNumberOfCalls(suite.T(), "Handle", 1)
	suite.handler.AssertNotCalled(suite.T(), "Handle", core.NewChangedIpEvent(ip, ip))
}

func (suite *AgentTestSuite) TestRunTwiceAndHandleAfterError() {
	ip := net.IPv4(127, 0, 0, 1)
	expectedErr := errors.New("test error")
	expectedEvent := core.NewChangedIpEvent(ip, net.IPv4zero)

	suite.resolver.On("Resolve").Return(ip, nil)
	suite.handler.On("Handle", expectedEvent).Return(expectedErr).Once()
	suite.handler.On("Handle", expectedEvent).Return(nil).Once()

	err := suite.agent.Run()
	suite.Error(err)
	suite.ErrorContains(err, expectedErr.Error())

	err = suite.agent.Run()
	suite.NoError(err)

	suite.handler.AssertExpectations(suite.T())
}

func TestAgentTestSuite(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}
