package core

import "net"

type Agent interface {
	Run() error
}

func NewAgent(r Resolver, h Handler) (Agent, error) {
	return &agentImpl{
		resolver:  r,
		handler:   h,
		currentIP: net.IPv4zero,
	}, nil
}

type agentImpl struct {
	resolver  Resolver
	handler   Handler
	currentIP net.IP
}

func (a *agentImpl) Run() error {
	var (
		err   error
		newIP net.IP
		oldIP = a.currentIP
	)

	if newIP, err = a.resolver.Resolve(); err != nil {
		return err
	}

	if newIP.Equal(oldIP) {
		return nil
	}

	if err = a.handler.Handle(NewChangedIpEvent(newIP, oldIP)); err != nil {
		return err
	}

	a.currentIP = newIP

	return nil
}
