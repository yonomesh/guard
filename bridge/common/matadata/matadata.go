package matadata

import (
	"guard/bridge/tools"
	"net"
	"net/netip"
)

type SocksAddr struct {
	AddrPort netip.AddrPort // IP & Port
	FQDN     string         // Fully Qualified Domain Name
}

func AddrFromIP(ip net.IP) netip.Addr {
	addr, _ := netip.AddrFromSlice(ip)
	return addr
}

// ParseAddr parses the given IP address string into a netip.Addr.
//
// Example:
//
//	addr1 := ParseAddr("192.168.1.1")
//	fmt.Println(addr1) // Output: 192.168.1.1
//
//	addr2 := ParseAddr("[2001:db8::1]")
//	fmt.Println(addr2) // Output: 2001:db8::1
//
//	addr3 := ParseAddr("::ffff:127.0.0.1")
//	fmt.Println(addr3) // Output: ::ffff:127.0.0.1
func ParseAddr(addr string) netip.Addr {
	return tools.Must1(netip.ParseAddr(unwrapIPv6Address(addr)))
}

// unwrapIPv6Address removes surrounding square brackets from an IPv6 address string.
//
// For example:
//
//	unwrapIPv6Address("[2001:db8::1]") → "2001:db8::1"
//	unwrapIPv6Address("[::1]")         → "::1"
func unwrapIPv6Address(addr string) string {
	if len(addr) > 2 && addr[0] == '[' && addr[len(addr)-1] == ']' {
		return addr[1 : len(addr)-1]
	}
	return addr
}
