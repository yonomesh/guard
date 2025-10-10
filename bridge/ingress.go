package bridge

import (
	"guard/bridge/matadata/ingress"
	"guard/bridge/matadata/packet"
	"net/netip"
	"time"
)

type IngressContext struct {
	Ingress     string
	IngressType ingress.Type
	IPVersion   packet.IPVersion
	Network     packet.Network
	Source      netip.AddrPort
	Destination netip.AddrPort

	// DPI
	DPIContext        any
	Protocol          packet.Protocol
	Domain            string
	DPIError          error
	ClientFingerprint string
	SSHClientName     string
	// eti
	ETI any

	// cache
	LastIngress              string
	OriginDestination        netip.AddrPort
	RouteOriginalDestination netip.Addr

	UDPDisableDomainUnmapping bool
	UDPConnect                bool
	UDPTimeout                time.Duration
	TLSFragment               bool
	TLSFragmentFallbackDelay  time.Duration
	TLSRecordFragment         bool

	DAddrs          []netip.Addr
	SourceGeoIPCode string
	GeoIPCode       string

	// DNS
	QueryType uint16
	Unreal    bool

	// rule cache
	IPCidrMatchSAddr  bool
	IPCidrAcceptEmpty bool

	SAddrMatch                   bool
	SPortMatch                   bool
	DAddrMatch                   bool
	DPortMatch                   bool
	IgnoreDestinationIPCIDRMatch bool

	// unknow
	DidMatch bool

	// no
	// ProcessInfo          *process.Info
	// NetworkStrategy     *C.NetworkStrategy
	// NetworkType         []C.InterfaceType
	// FallbackNetworkType []C.InterfaceType
	// FallbackDelay       time.Duration
	// // Deprecated: to be removed
	// // nolint:staticcheck
	// InboundOptions option.InboundOptions

	// // Deprecated: implement in rule action
	// InboundDetour string
}

type IngressContext2 struct {

	// cache

}
