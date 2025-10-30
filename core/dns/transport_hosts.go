package dns

import (
	"net/netip"

	"github.com/miekg/dns"
)

// FixedResponse creates a DNS response message with either A or AAAA records based on the provided addresses.
//
// It takes the following parameters:
//
// - id: The ID to be set in the DNS header
//
// - question: The DNS question being answered
//
// - addrs: A slice of IP addresses (IPv4 or IPv6) to be included in the response
//
// - timeToLive: The TTL (Time to Live) value for the DNS records
//
// It returns a pointer to a dns.Msg containing the response, including A or AAAA records.
func FixedResponse(id uint16, question dns.Question, addrs []netip.Addr, timeToLive uint32) *dns.Msg {
	response := dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:       id,
			Rcode:    dns.RcodeSuccess,
			Response: true,
		},
		Question: []dns.Question{question},
	}

	for _, addr := range addrs {
		if addr.Is4() {
			response.Answer = append(response.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    timeToLive,
				},
				A: addr.AsSlice(),
			})
		} else {
			response.Answer = append(response.Answer, &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    timeToLive,
				},
				AAAA: addr.AsSlice(),
			})
		}
	}

	return &response
}
