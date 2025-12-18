package dns

import (
	"context"
	"net/netip"
	"net/url"

	E "guard/bridge/common/errors"
	"guard/bridge/common/logging"
	N "guard/bridge/common/network"
	"guard/core/dns/server/unreal"

	"github.com/yonomesh/dns"
)

type Transport interface {
	Name() string
	Start() error
	Reset()
	Close()
	Raw() bool
	Exchange(ctx context.Context, message *dns.Msg) (*dns.Msg, error)
	Lookup(ctx context.Context, domain string, strategy unreal.DNSUnrealStrategy) ([]netip.Addr, error)
}

type TransportOptions struct {
	Name         string
	Context      context.Context
	Dialer       N.Dialer
	Address      string
	ClientSubnet netip.Prefix // ? is for edns0_subnet 我不确定
	Logger       logging.ContextLogger
}

type TransportConstructor = func(options TransportOptions) (Transport, error)

// transport map
// 不确定
// address----TransportConstructor
var transports map[string]TransportConstructor

func RegisterTransport(schemes []string, constructor TransportConstructor) {
	if transports == nil {
		transports = make(map[string]TransportConstructor)
	}

	for _, scheme := range schemes {
		transports[scheme] = constructor
	}
}

func CreateTransport(options TransportOptions) (Transport, error) {
	constructor := transports[options.Address]
	if constructor == nil {
		serverURL, _ := url.Parse(options.Address)
		var scheme string
		if serverURL != nil {
			scheme = serverURL.Scheme
		}
		constructor = transports[scheme] // WTF
	}

	if constructor == nil {
		return nil, E.New("unknown DNS server format: " + options.Address)
	}
	options.Context = contextWithTransportName(options.Context, options.Name)
	transport, err := constructor(options)
	if err != nil {
		return nil, err
	}
	if options.ClientSubnet.IsValid() {
		transport = &edns0SubnetTransportWrapper{transport, options.ClientSubnet}
	}

	return transport, nil
}
