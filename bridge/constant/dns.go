package constant

const (
	DNSDefaultTTL = 300
)

const (
	DNSUnrealStrategyDefault = iota
	DNSUnrealStrategyPreferIPv4
	DNSUnrealStrategyPreferIPv6
	DNSUnrealStrategyOnlyIPv4
	DNSUnrealStrategyOnlyIPv6
)

const (
	DNSForwarderTypeStandard = "standard"
	DNSForwarderTypeUDP      = "udp"
	DNSForwarderTypeTCP      = "tcp"
	DNSForwarderTypeTLS      = "tls"
	DNSForwarderTypeDOT      = "dot"
	DNSForwarderTypeQUIC     = "quic"
	DNSForwarderTypeHTTP     = "http"
	DNSForwarderTypeDOH      = "doh"
	DNSForwarderTypeHTTP2    = "h2"
	DNSForwarderTypeHTTP3    = "h3"
)

const (
	DNSProvideCloudflare = "cloudflare"
	DNSProvideGoogle     = "google"
	DNSProvideAlibaba    = "alibaba"
)

const (
	DNSActionForward  = "forward"
	DNSActionServer   = "server"
	DNSActionDrop     = "drop"
	DNSActionNotFound = "not-found"
	DNSActionFinal    = "final"
)
