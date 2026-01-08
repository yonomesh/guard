package bridge

import (
	"net/netip"
	"time"

	M "uni/bridge/common/matadata"
)

type IngressContext struct {
	Ingress     string
	IngressType string
	// 这里还不清楚到底是什么
	Network     string
	IPVersion   uint8
	Source      M.SocksAddr
	Destination M.SocksAddr

	// DPI
	DPIContext        any
	Proto             string
	Domain            string
	DPIError          error
	ClientFingerprint string
	SSHClientName     string
	// eti
	ETI any

	// cache
	LastIngress              string
	OriginDestination        M.SocksAddr
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
