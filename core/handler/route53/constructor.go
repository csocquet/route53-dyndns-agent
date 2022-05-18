package route53

import "github.com/csocquet/route53-dyndns-agent/core"

type OptionFn func(*handler) error

func NewRoute53Handler(opts ...OptionFn) (core.Handler, error) {
	h := &handler{}

	for _, opt := range opts {
		_ = opt(h)
	}

	return h, nil
}

func WithApiClient(cli ApiClient) OptionFn {
	return func(h *handler) error {
		h.client = cli
		return nil
	}
}
