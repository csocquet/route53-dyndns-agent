package core

import "net"

func NewChangedIpEvent(newIp net.IP, oldIp net.IP) *ChangedIpEvent {
	return &ChangedIpEvent{
		NewIP: newIp,
		OldIP: oldIp,
	}
}

type ChangedIpEvent struct {
	NewIP net.IP
	OldIP net.IP
}

type Handler interface {
	Handle(*ChangedIpEvent) error
}
