package route53

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/csocquet/route53-dyndns-agent/core"
)

type ApiClient interface {
	ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(options *route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)
}

type handler struct {
	client ApiClient
}

func (h *handler) Handle(e *core.ChangedIpEvent) error {

	_, _ = h.client.ChangeResourceRecordSets(context.TODO(), &route53.ChangeResourceRecordSetsInput{}, nil)

	return nil
}
