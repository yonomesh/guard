package network

import (
	"context"
	"net"
)

type Dialer interface {
}

type DefaultDialer struct {
	net.Dialer
	net.ListenConfig
}

func (d *DefaultDialer) DialContext(ctx context.Context, network string) {

}
