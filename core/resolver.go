package core

import "net"

type Resolver interface {
	Resolve() (net.IP, error)
}
