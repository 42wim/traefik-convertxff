package convertxff

import (
	"context"
	"net/http"
	"net/netip"
	"strings"
)

const (
	XFF = "X-Forwarded-For"
)

type Config struct{}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// ConvertXFF holds the necessary components of a Traefik plugin
type ConvertXFF struct {
	next http.Handler
	name string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &ConvertXFF{
		next: next,
		name: name,
	}, nil
}

// ServeHTTP removes the brackets from IPv6 addresses and rewrites ipv6 mapped ipv4 to real ipv4
func (u *ConvertXFF) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(XFF)

	req.Header.Del(XFF)

	header = strings.ReplaceAll(header, "[", "")
	header = strings.ReplaceAll(header, "]", "")

	var ips []string

	if strings.Contains(header, ",") {
		for _, entry := range strings.Split(header, ",") {
			entry = strings.TrimSpace(entry)

			if !strings.HasPrefix(entry, "::ffff") {
				ips = append(ips, entry)
				continue
			}

			parsed, err := netip.ParseAddr(entry)
			if err != nil {
				ips = append(ips, entry)
				continue
			}

			ips = append(ips, parsed.Unmap().String())
		}
	} else {
		ips = append(ips, header)
	}

	req.Header.Add(XFF, strings.Join(ips, ","))
	u.next.ServeHTTP(rw, req)
}
