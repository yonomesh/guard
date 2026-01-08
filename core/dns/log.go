package dns

import (
	"context"
	"strings"

	"uni/bridge/common/logging"

	"github.com/miekg/dns"
)

func LogCachedResponse(logging logging.ContextLogger, ctx context.Context, response *dns.Msg, ttl int) {
	if logging == nil || len(response.Question) == 0 {
		return
	}
	domain := (FqdnToDomain(response.Question[0].Name))
	logging.DebugContext(ctx, "cached ", domain, " ", dns.RcodeToString[response.Rcode], " ", ttl)

	// [][]dns.RR{} is an anonymous slice, its type is []dns.RR
	for _, recordList := range [][]dns.RR{response.Answer, response.Ns, response.Extra} {
		for _, record := range recordList {
			logging.InfoContext(ctx, "cache ", dns.Type(record.Header().Rrtype).String(), " ", FormatQuestion(record.String()))
		}
	}
}

func LogExchangeResponse(logging logging.ContextLogger, ctx context.Context, response *dns.Msg, ttl uint32) {
	if logging == nil || len(response.Question) == 0 {
		return
	}

	domain := (FqdnToDomain(response.Question[0].Name))
	logging.DebugContext(ctx, "exchanged ", domain, " ", dns.RcodeToString[response.Rcode], " ", ttl)
	for _, recordList := range [][]dns.RR{response.Answer, response.Ns, response.Extra} {
		for _, record := range recordList {
			logging.InfoContext(ctx, "exchange ", dns.Type(record.Header().Rrtype).String(), " ", FormatQuestion(record.String()))
		}
	}

}

func LogRejectResponse(logging logging.ContextLogger, ctx context.Context, response *dns.Msg) {
	if logging == nil || len(response.Question) == 0 {
		return
	}

	logging.DebugContext(ctx, "rejected dns")
	for _, recordList := range [][]dns.RR{response.Answer, response.Ns, response.Extra} {
		for _, record := range recordList {
			logging.InfoContext(ctx, "reject ", dns.Type(record.Header().Rrtype).String(), " ", FormatQuestion(record.String()))
		}
	}

}

// fqdnToDomain convert an Fully Qualified Domain Name string to normal domain string
//
// Example:
//
//	fmt.Println(FqdnToDomain("example.com.")) // output example.com
//	fmt.Println(FqdnToDomain("example.com")) // output example.com
//	fmt.Println(FqdnToDomain(".")) // ""
func FqdnToDomain(fqdn string) string {
	if dns.IsFqdn(fqdn) {
		return fqdn[:len(fqdn)-1]
	}
	return fqdn
}

// FormatQuestion formats question string
func FormatQuestion(str string) string {
	for strings.HasPrefix(str, ";") {
		str = str[1:]
	}
	str = strings.ReplaceAll(str, "\t", " ")
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, ";; ", " ")
	str = strings.ReplaceAll(str, "; ", " ")
	for strings.Contains(str, "  ") {
		str = strings.ReplaceAll(str, "  ", " ")
	}
	return strings.TrimSpace(str)
}
