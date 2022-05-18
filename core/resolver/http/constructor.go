package http

import (
	"fmt"
	"github.com/csocquet/route53-dyndns-agent/core"
	"net/http"
	url2 "net/url"
	"regexp"
)

const (
	DEFAULT_URL    = "https://checkip.amazonaws.com"
	DEFAULT_REGEXP = `(?:\d{1,3}\.){3}\d{1,3}`
)

type ResolverOpt func(*resolver) error

func NewHttpResolver(opts ...ResolverOpt) (core.Resolver, error) {
	r := &resolver{
		httpClient: http.DefaultClient,
		url:        DEFAULT_URL,
		regexp:     regexp.MustCompile(DEFAULT_REGEXP),
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func WithHttpClient(c HttpClient) ResolverOpt {
	return func(r *resolver) error {
		r.httpClient = c
		return nil
	}
}

func WithUrl(url string) ResolverOpt {
	return func(r *resolver) error {
		if _, err := url2.ParseRequestURI(url); err != nil {
			return fmt.Errorf("invalid url `%s`: %w", url, err)
		}
		r.url = url
		return nil
	}
}

func WithRegexp(expr string) ResolverOpt {
	return func(r *resolver) error {
		rgx, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("invalid regexp `%s`: %w", expr, err)
		}

		r.regexp = rgx
		return nil
	}
}
