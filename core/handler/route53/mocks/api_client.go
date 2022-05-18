package mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/stretchr/testify/mock"
)

type ApiClientMock struct {
	mock.Mock
}

func (m *ApiClientMock) ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(options *route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.ChangeResourceRecordSetsOutput), args.Error(1)
}
