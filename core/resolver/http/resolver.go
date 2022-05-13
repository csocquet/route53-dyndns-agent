package http

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
)

type HttpClient interface {
	Get(string) (*http.Response, error)
}

type resolver struct {
	httpClient HttpClient
	url        string
	regexp     *regexp.Regexp
}

func (r resolver) Resolve() (net.IP, error) {
	raw, err := r.callHttp()
	if err != nil {
		return net.IPv4zero, err
	}

	rawIp := r.extractRawIp(raw)

	return r.parseIPv4(rawIp)
}

func (r resolver) callHttp() ([]byte, error) {
	resp, err := r.httpClient.Get(r.url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"invalid http client response: %d - %s",
			resp.StatusCode,
			resp.Status,
		)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (r resolver) extractRawIp(data []byte) []byte {
	if rawIp, ok := r.extractRawIpFromWithIpPartsGroup(data); ok {
		return rawIp
	}

	if rawIp, ok := r.extractRawIpFromWithIpGroup(data); ok {
		return rawIp
	}

	return r.extractRawIpFromMatch(data)
}

func (r resolver) extractRawIpFromWithIpPartsGroup(data []byte) ([]byte, bool) {
	idx1 := r.regexp.SubexpIndex("ip1")
	idx2 := r.regexp.SubexpIndex("ip2")
	idx3 := r.regexp.SubexpIndex("ip3")
	idx4 := r.regexp.SubexpIndex("ip4")

	if idx1 == -1 || idx2 == -1 || idx3 == -1 || idx4 == -1 {
		return nil, false
	}

	matches := r.regexp.FindSubmatch(data)
	if len(matches) < idx1 || len(matches) < idx2 || len(matches) < idx3 || len(matches) < idx4 {
		return nil, true
	}

	rawIp := []byte(fmt.Sprintf(
		"%s.%s.%s.%s",
		matches[idx1],
		matches[idx2],
		matches[idx3],
		matches[idx4],
	))

	return rawIp, true
}

func (r resolver) extractRawIpFromWithIpGroup(data []byte) ([]byte, bool) {
	idx := r.regexp.SubexpIndex("ip")
	if idx == -1 {
		return nil, false
	}

	matches := r.regexp.FindSubmatch(data)
	if len(matches) < idx {
		return nil, true
	}

	return matches[idx], true
}

func (r resolver) extractRawIpFromMatch(data []byte) []byte {
	return r.regexp.Find(data)
}

func (r resolver) parseIPv4(raw []byte) (net.IP, error) {
	ip := net.ParseIP(string(raw))
	if ip == nil || ip.To4() == nil {
		return net.IPv4zero, fmt.Errorf("failed to parse IPv4 `%s`", raw)
	}

	return ip, nil
}
