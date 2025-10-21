package constant

import "time"

const (
	TCPKeepAliveInitial        = 10 * time.Minute
	TCPKeepAliveInterval       = 75 * time.Second
	TCPConnectTimeout          = 5 * time.Second
	TCPTimeout                 = 15 * time.Second
	ReadPayloadTimeout         = 300 * time.Millisecond
	DNSTimeout                 = 10 * time.Second
	UDPTimeout                 = 5 * time.Minute
	DefaultURLTestInterval     = 3 * time.Minute
	DefaultURLTestIdleTimeout  = 30 * time.Minute
	StartTimeout               = 10 * time.Second
	StopTimeout                = 5 * time.Second
	FatalStopTimeout           = 10 * time.Second
	FakeIPMetadataSaveInterval = 10 * time.Second
	TLSFragmentFallbackDelay   = 500 * time.Millisecond
)

var PortProtocols = map[uint16]string{
	53:   ProtoDNS,
	123:  ProtoNTP,
	3478: ProtoSTUN,
	443:  ProtoQUIC,
}

var ProtocolTimeouts = map[string]time.Duration{
	ProtoDNS:  10 * time.Second,
	ProtoNTP:  10 * time.Second,
	ProtoSTUN: 10 * time.Second,
	ProtoQUIC: 30 * time.Second,
	ProtoDtls: 30 * time.Second,
}
