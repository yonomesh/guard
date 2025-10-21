package constant

const (
	IPv4 uint8 = iota
	IPv6
)

// Application Protocol
const (
	ProtoDNS        = "dns"
	ProtoHTTP       = "http"
	ProtoTLS        = "tls"
	ProtoDtls       = "dtls"
	ProtoSSH        = "ssh"
	ProtoQUIC       = "quic"
	ProtoSTUN       = "stun"
	ProtoRDP        = "rdp"
	ProtoNTP        = "ntp"
	ProtoBittorrent = "bittorrent"

	// Encrypted Traffic Inspector
	// MITM
	ProtoTLS2  = "tls1.2"
	ProtoTLS3  = "tls1.3"
	ProtoHTTP2 = "http2"
	ProtoHTTP3 = "http3"

	// TODO
	ProtoFTP       = "ftp"
	ProtoWebSocket = "ws"
	ProtoSMB       = "smb"
	ProtoMQTT      = "mqtt"
)

const (
	ClientFingerprintUnknown  = "unknown"
	ClientFingerprintChromium = "chromium"
	ClientFingerprintSafari   = "safari"
	ClientFingerprintFirefox  = "firefox"
	ClientFingerprintQuicGO   = "quic-go"
)
