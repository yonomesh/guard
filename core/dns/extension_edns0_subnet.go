package dns

import (
	"net/netip"

	"github.com/miekg/dns"
)

type edns0SubnetTransportWrapper struct {
	Transport
	clientSubnet netip.Prefix
}

func SetClientSubnet(msg *dns.Msg, clientSubnet netip.Prefix, override bool) *dns.Msg {
	return setClientSubnet(msg, clientSubnet, override, true)
}

func setClientSubnet(msg *dns.Msg, clientSubnet netip.Prefix, override bool, clone bool) *dns.Msg {
	var (
		optRecord    *dns.OPT
		subnetOption *dns.EDNS0_SUBNET
	)
findExists:
	for _, record := range msg.Extra {
		var ok bool
		if optRecord, ok = record.(*dns.OPT); ok {
			for _, option := range optRecord.Option {
				subnetOption, ok = option.(*dns.EDNS0_SUBNET)
				if ok {
					if !override {
						return msg
					}
					break findExists
				}
			}
		}
	}

	if optRecord == nil {
		exMsg := *msg
		msg = &exMsg
		optRecord = &dns.OPT{
			Hdr: dns.RR_Header{
				Name:   ".",
				Rrtype: dns.TypeOPT,
			},
		}

		msg.Extra = append(msg.Extra, optRecord)
	} else if clone {
		return setClientSubnet(msg.Copy(), clientSubnet, override, false)
	}

	if subnetOption == nil {
		subnetOption = new(dns.EDNS0_SUBNET)
		subnetOption.Code = dns.EDNS0SUBNET
		optRecord.Option = append(optRecord.Option, subnetOption)
	}

	if clientSubnet.Addr().Is4() {
		subnetOption.Family = 1
	} else {
		subnetOption.Family = 2
	}

	subnetOption.SourceNetmask = uint8(clientSubnet.Bits())
	subnetOption.Address = clientSubnet.Addr().AsSlice()
	return msg
}
