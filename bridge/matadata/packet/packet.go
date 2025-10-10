package packet

type Network uint8

const (
	TCP Network = iota
	UDP
)

type IPVersion uint8

const (
	IPv4 IPVersion = iota
	IPv6
)

type Protocol string

const (
	Dns        = "dns"
	Http       = "http"
	TLS        = "tls"
	Dtls       = "dtls"
	SSH        = "ssh"
	Quic       = "quic"
	Stun       = "stun"
	RDP        = "rdp"
	NTP        = "ntp"
	ICMP       = "icmp"
	Bittorrent = "bittorrent"

	// Encrypted Traffic Inspector
	// MITM
	TLS2  = "tls1.2"
	TLS3  = "tls1.3"
	Http2 = "http2"
	Http3 = "http3"

	// TODO
	FTP       = "ftp"
	WebSocket = "ws"
	SMB       = "smb"
	Mqtt      = "mqtt"
)
