package dns

import (
	"context"
	"net/netip"
	"strings"
	"time"

	E "guard/bridge/common/errors"
	"guard/bridge/common/logging"
	M "guard/bridge/common/matadata"
	"guard/bridge/common/task"
	"guard/bridge/tools"
	"guard/bridge/tools/freelru"
	"guard/bridge/tools/maphash"
	"guard/core/dns/server/unreal"

	"github.com/yonomesh/dns"
)

const (
	DefaultTTL     int           = 600
	DefaultTimeout time.Duration = 10 * time.Second
)

var (
	ErrNotRawSupport          = E.New("now raw query support by current transport")
	ErrNotCached              = E.New("not cached")
	ErrResponseRejected       = E.New("response rejected")
	ErrResponseRejectedCached = E.Extend(ErrResponseRejected, "cached")
)

// DNS Resolver
type Resolver struct {
	timeout        time.Duration
	disableCache   bool
	sharedCache    bool
	disableExpire  bool
	rdrc           RDRCStore
	initRDRCStore  func() RDRCStore
	logger         logging.ContextLogger
	cache          freelru.Cache[dns.Question, *dns.Msg]
	transportCache freelru.Cache[transportCacheKey, *dns.Msg]
}

type ResolverOptions struct {
	Timeout       time.Duration
	DisableCache  bool
	SharedCache   bool
	DisableExpire bool
	CacheCapacity uint32
	RDRC          func() RDRCStore
	Logger        logging.ContextLogger
}

type QueryOptions struct {
	UnrealStrategy unreal.DNSUnrealStrategy
	DisableCache   bool
	TTL            *TTLOptions
	ClientSubnet   netip.Prefix
}

type TTLOptions struct {
	Type    string
	Max     uint32
	Min     uint32
	ReWrite uint32
}

type transportCacheKey struct {
	dns.Question
	transportName string
}

// Rejected DNS Response Cache
// Response Data Rejection Cache
// Recursive DNS Result Cache ?
// 我不确定，大概率是拒绝
type RDRCStore interface {
	Load(transportName string, QName string, QType uint16) (rejected bool)
	Save(transportName string, QName string, QType uint16) error
	SaveAsync(transportName string, QName string, QType uint16, logger logging.Logger)
}

func NewResolver(options ResolverOptions) *Resolver {
	resolver := &Resolver{
		timeout:       options.Timeout,
		disableCache:  options.DisableCache,
		sharedCache:   options.SharedCache,
		disableExpire: options.DisableExpire,
		initRDRCStore: options.RDRC,
		logger:        options.Logger,
	}

	if resolver.timeout == 0 {
		resolver.timeout = DefaultTimeout
	}

	cacheCapacity := options.CacheCapacity
	if cacheCapacity < 1024 {
		cacheCapacity = 1024
	}
	if resolver.sharedCache {
		resolver.cache = tools.Must1(freelru.NewSharded[dns.Question, *dns.Msg](cacheCapacity, maphash.NewHasher[dns.Question]().Hash32))
	} else {
		resolver.transportCache = tools.Must1(freelru.NewSharded[transportCacheKey, *dns.Msg](cacheCapacity, maphash.NewHasher[transportCacheKey]().Hash32))
	}

	return resolver
}

func (r *Resolver) Start() {
	if r.initRDRCStore != nil {
		r.rdrc = r.initRDRCStore()
	}
}

// 核心查询函数
func (r *Resolver) Exchange(ctx context.Context, transport Transport, msg *dns.Msg, options QueryOptions) (*dns.Msg, error) {
	return r.ExchangeWithResponseCheck(ctx, transport, msg, options, nil)
}

// 这个函数在处理 DNS 查询时，先尝试从缓存中获取结果，如果缓存没有命中，则发送实际的查询请求并进行响应验证，最后对响应进行一些处理（如TTL调整、日志记录等）。
func (r *Resolver) ExchangeWithResponseCheck(ctx context.Context, transport Transport, msg *dns.Msg, options QueryOptions, responseChecker func(responseAddrs []netip.Addr) bool) (*dns.Msg, error) {
	if len(msg.Question) == 0 {
		if r.logger != nil {
			r.logger.WarnContext(ctx, "question is 0")
		}
		responseMsg := dns.Msg{
			MsgHdr: dns.MsgHdr{
				Id:       msg.Id,
				Response: true,
				Rcode:    dns.RcodeFormatError,
			},
			Question: msg.Question,
		}
		return &responseMsg, nil
	}

	question := msg.Question[0]
	if options.ClientSubnet.IsValid() {
		msg = SetClientSubnet(msg, options.ClientSubnet, true)
	}

	isSimpleRequest := len(msg.Question) == 1 &&
		len(msg.Ns) == 0 &&
		len(msg.Extra) == 0 &&
		!options.ClientSubnet.IsValid()

	disableCache := !isSimpleRequest || r.disableCache || options.DisableCache
	if !disableCache {
		response, ttl := r.loadResponse(question, transport)
		if response != nil {
			LogCachedResponse(r.logger, ctx, response, ttl)
			response.Id = msg.Id
			return response, nil
		}
	}

	if question.Qtype == dns.TypeA && options.UnrealStrategy == unreal.DNSUnrealStrategyOnlyIPv6 || question.Qtype == dns.TypeAAAA && options.UnrealStrategy == unreal.DNSUnrealStrategyOnlyIPv4 {
		responseMsg := dns.Msg{
			MsgHdr: dns.MsgHdr{
				Id:       msg.Id,
				Response: true,
				Rcode:    dns.RcodeSuccess,
			},
			Question: []dns.Question{question},
		}
		if r.logger != nil {
			r.logger.DebugContext(ctx, "strategy rejected")
		}

		return &responseMsg, nil
	}

	if !transport.Raw() {
		if question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA {
			return r.exchangeToLookup(ctx, transport, msg, question, options, responseChecker)
		}
		return nil, ErrNotRawSupport
	}

	msgID := msg.Id
	contextTransport, ok := transportNameFromContext(ctx)
	if ok && transport.Name() == contextTransport {
		return nil, E.New("DNS query loopback in transport[", contextTransport, "]")
	}

	ctx = contextWithTransportName(ctx, transport.Name())
	if responseChecker != nil && r.rdrc != nil {
		rejected := r.rdrc.Load(transport.Name(), question.Name, question.Qtype)
		if rejected {
			return nil, ErrResponseRejectedCached
		}
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)

	response, err := transport.Exchange(ctx, msg) // DNS Query
	cancel()
	if err != nil {
		return nil, err
	}

	if responseChecker != nil {
		var rejected bool
		if !(response.Rcode == dns.RcodeSuccess || response.Rcode == dns.RcodeNameError) {
			rejected = true
		} else {
			rejected = !responseChecker(MsgToAddrs(response))
		}
		if rejected {
			if r.rdrc != nil {
				r.rdrc.SaveAsync(transport.Name(), question.Name, question.Qtype, r.logger)
			}
			LogRejectResponse(r.logger, ctx, response)
			return response, ErrResponseRejected
		}
	}

	if question.Qtype == dns.TypeHTTPS {
		if options.UnrealStrategy == unreal.DNSUnrealStrategyOnlyIPv4 || options.UnrealStrategy == unreal.DNSUnrealStrategyOnlyIPv6 {
			for _, answer := range response.Answer {
				https, ok := answer.(*dns.HTTPS)
				if !ok {
					continue
				}

				content := https.SVCB
				content.Value = tools.Filter(content.Value, func(it dns.SVCBKeyValue) bool {
					if options.UnrealStrategy == unreal.DNSUnrealStrategyOnlyIPv4 {
						return it.Key() != dns.SVCB_IPV6HINT
					} else {
						return it.Key() != dns.SVCB_IPV4HINT
					}
				})
				https.SVCB = content
			}

		}
	}

	var timeToLive uint32
	for _, recordList := range [][]dns.RR{response.Answer, response.Ns, response.Extra} {
		for _, record := range recordList {
			if timeToLive == 0 || record.Header().Ttl > 0 && record.Header().Ttl < timeToLive {
				timeToLive = record.Header().Ttl
			}
		}
	}

	if options.TTL != nil {
		timeToLive = options.TTL.ReWrite
	}

	for _, recordList := range [][]dns.RR{response.Answer, response.Ns, response.Extra} {
		for _, record := range recordList {
			record.Header().Ttl = timeToLive
		}
	}

	response.Id = msgID
	if !disableCache {
		r.storeCache(transport, question, response, timeToLive)
	}
	LogExchangeResponse(r.logger, ctx, response, timeToLive)
	return response, err
}

func (r *Resolver) exchangeToLookup(ctx context.Context, transport Transport, msg *dns.Msg, question dns.Question, options QueryOptions, responseChecker func(responseAddrs []netip.Addr) bool) (*dns.Msg, error) {
	domain := question.Name
	if question.Qtype == dns.TypeA {
		options.UnrealStrategy = unreal.DNSUnrealStrategyOnlyIPv4
	} else {
		options.UnrealStrategy = unreal.DNSUnrealStrategyOnlyIPv6
	}

	result, err := r.LookupWithResponseCheck(ctx, transport, domain, options, responseChecker)
	if err != nil {
		return nil, wrapError(err)
	}

	var timeToLive uint32
	if options.TTL != nil {
		timeToLive = options.TTL.ReWrite
	} else {
		timeToLive = uint32(DefaultTTL)
	}

	response := FixedResponse(msg.Id, question, result, timeToLive)
	LogExchangeResponse(r.logger, ctx, response, timeToLive)

	return &dns.Msg{}, nil
}

func (r *Resolver) Lookup(ctx context.Context, transport Transport, domain string, options QueryOptions) ([]netip.Addr, error) {
	return r.LookupWithResponseCheck(ctx, transport, domain, options, nil)
}

func (r *Resolver) LookupWithResponseCheck(ctx context.Context, transport Transport, domain string, options QueryOptions, responseChecker func(responseAddrs []netip.Addr) bool) ([]netip.Addr, error) {
	if dns.IsFqdn(domain) {
		domain = domain[:len(domain)-1]
	}

	dnsName := dns.Fqdn(domain)

	// transport Raw
	if transport.Raw() {
		switch options.UnrealStrategy {
		case unreal.DNSUnrealStrategyOnlyIPv4:
			return r.lookupToExchange(ctx, transport, dnsName, dns.TypeA, options, responseChecker)
		case unreal.DNSUnrealStrategyOnlyIPv6:
			return r.lookupToExchange(ctx, transport, dnsName, dns.TypeAAAA, options, responseChecker)
		}

		var response4 []netip.Addr
		var response6 []netip.Addr

		var group task.Group
		group.Append("exchange4", func(ctx context.Context) error {
			response, err := r.lookupToExchange(ctx, transport, dnsName, dns.TypeA, options, responseChecker)
			if err != nil {
				return err
			}
			response4 = response
			return nil
		})

		group.Append("exchange6", func(ctx context.Context) error {
			response, err := r.lookupToExchange(ctx, transport, dnsName, dns.TypeAAAA, options, responseChecker)
			if err != nil {
				return err
			}
			response6 = response
			return nil
		})

		err := group.Run(ctx)
		if len(response4) == 0 && len(response6) == 0 {
			return nil, err
		}
		return sortAddrs(response4, response6, options.UnrealStrategy), nil
	}

	// disableCache
	disableCache := r.disableCache || options.DisableCache
	if !disableCache {
		if options.UnrealStrategy == unreal.DNSUnrealStrategyOnlyIPv4 {
			response, err := r.questionCache(dns.Question{
				Name:   dnsName,
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			}, transport)
			if err != ErrNotCached {
				return response, err
			}
		} else if options.UnrealStrategy == unreal.DNSUnrealStrategyOnlyIPv6 {
			response, err := r.questionCache(dns.Question{
				Name:   dnsName,
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassINET,
			}, transport)
			if err != ErrNotCached {
				return response, err
			}
		} else {
			response4, _ := r.questionCache(dns.Question{
				Name:   dnsName,
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			}, transport)
			response6, _ := r.questionCache(dns.Question{
				Name:   dnsName,
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassINET,
			}, transport)
			if len(response4) > 0 || len(response6) > 0 {
				return sortAddrs(response4, response6, options.UnrealStrategy), nil
			}
		}
	}

	// responseChecker
	if responseChecker != nil && r.rdrc != nil {
		var rejected bool
		if options.UnrealStrategy != unreal.DNSUnrealStrategyOnlyIPv6 {
			rejected = r.rdrc.Load(transport.Name(), dnsName, dns.TypeA)
		}
		if !rejected && options.UnrealStrategy != unreal.DNSUnrealStrategyOnlyIPv4 {
			rejected = r.rdrc.Load(transport.Name(), dnsName, dns.TypeAAAA)
		}
		if rejected {
			return nil, ErrResponseRejectedCached
		}
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	response, err := transport.Lookup(ctx, domain, options.UnrealStrategy)
	cancel()
	if err != nil {
		return nil, wrapError(err)
	}
	if responseChecker != nil && !responseChecker(response) {
		if r.rdrc != nil {
			if tools.Any(response, func(addr netip.Addr) bool {
				return addr.Is4()
			}) {
				r.rdrc.SaveAsync(transport.Name(), dnsName, dns.TypeA, r.logger)
			}

			if tools.Any(response, func(addr netip.Addr) bool {
				return addr.Is6()
			}) {
				r.rdrc.SaveAsync(transport.Name(), dnsName, dns.TypeAAAA, r.logger)
			}
		}
		LogRejectResponse(r.logger, ctx, FixedResponse(0, dns.Question{}, response, uint32(DefaultTTL)))
		return response, ErrResponseRejected
	}

	var rCode int
	header := dns.MsgHdr{
		Response: true,
		Rcode:    rCode,
	}
	if !disableCache {
		var timeToLive uint32
		if options.TTL != nil {
			timeToLive = options.TTL.ReWrite
		} else {
			timeToLive = uint32(DefaultTTL)
		}
		if options.UnrealStrategy != unreal.DNSUnrealStrategyOnlyIPv6 {
			question4 := dns.Question{
				Name:   dnsName,
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			}
			response4 := tools.Filter(response, func(addr netip.Addr) bool {
				return addr.Is4() || addr.Is4In6()
			})

			msg4 := &dns.Msg{
				MsgHdr:   header,
				Question: []dns.Question{question4},
			}

			if len(response4) > 0 {
				for _, addr := range response4 {
					msg4.Answer = append(msg4.Answer, &dns.A{
						Hdr: dns.RR_Header{
							Name:   question4.Name,
							Rrtype: dns.TypeA,
							Class:  dns.ClassINET,
							Ttl:    timeToLive,
						},
						A: addr.AsSlice(),
					})
				}
			}
			r.storeCache(transport, question4, msg4, timeToLive)
		}
		if options.UnrealStrategy != unreal.DNSUnrealStrategyOnlyIPv4 {
			question6 := dns.Question{
				Name:   dnsName,
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassINET,
			}
			response6 := tools.Filter(response, func(addr netip.Addr) bool {
				return addr.Is6() && !addr.Is4In6()
			})
			msg6 := &dns.Msg{
				MsgHdr:   header,
				Question: []dns.Question{question6},
			}
			if len(response6) > 0 {
				for _, addr := range response6 {
					msg6.Answer = append(msg6.Answer, &dns.AAAA{
						Hdr: dns.RR_Header{
							Name:   question6.Name,
							Rrtype: dns.TypeAAAA,
							Class:  dns.ClassINET,
							Ttl:    uint32(DefaultTTL),
						},
						AAAA: addr.AsSlice(),
					})
				}
			}
			r.storeCache(transport, question6, msg6, timeToLive)
		}
	}
	return response, nil
}

// 执行一次 DNS 查询，
// 优先尝试从缓存获取结果；若缓存无效，则通过指定的传输方式发送查询请求，
// 接收响应、校验返回码，并提取响应中的所有 IP 地址。
func (r *Resolver) lookupToExchange(ctx context.Context, transport Transport, name string, qType uint16, options QueryOptions, responseChecker func(responseAddrs []netip.Addr) bool) ([]netip.Addr, error) {
	question := dns.Question{
		Name:   name,
		Qtype:  qType,
		Qclass: dns.ClassINET,
	}

	disableCache := r.disableCache || options.DisableCache
	if !disableCache {
		cacheAddrs, err := r.questionCache(question, transport)
		if err != ErrNotCached {
			return cacheAddrs, err
		}
	}

	msg := dns.Msg{
		MsgHdr: dns.MsgHdr{
			RecursionDesired: true,
		},
		Question: []dns.Question{question},
	}

	var (
		response *dns.Msg
		err      error
	)

	if responseChecker != nil {
		response, err = r.ExchangeWithResponseCheck(ctx, transport, &msg, options, responseChecker)
	} else {
		response, err = r.Exchange(ctx, transport, &msg, options)
	}

	if err != nil {
		return nil, err
	}

	if response.Rcode != dns.RcodeSuccess {
		return nil, RCodeError(response.Rcode)
	}

	return MsgToAddrs(response), nil
}

// 从缓存中读取 DNS 结果，并提取其中的 IP 地址。
func (r *Resolver) questionCache(question dns.Question, transport Transport) ([]netip.Addr, error) {
	response, _ := r.loadResponse(question, transport)
	if response == nil {
		return nil, ErrNotCached
	}

	if response.Rcode != dns.RcodeSuccess {
		return nil, RCodeError(response.Rcode)
	}

	return MsgToAddrs(response), nil
}

// MsgToAddrs extracts all IP addresses from a DNS message response.
func MsgToAddrs(response *dns.Msg) []netip.Addr {
	addrs := make([]netip.Addr, 0, len(response.Answer))
	for _, rawAnswer := range response.Answer {
		switch answer := rawAnswer.(type) {
		case *dns.A:
			addrs = append(addrs, M.AddrFromIP(answer.A))
		case *dns.AAAA:
			addrs = append(addrs, M.AddrFromIP(answer.AAAA))
		case *dns.HTTPS:
			// WTF?
			for _, v := range answer.SVCB.Value {
				if v.Key() == dns.SVCB_IPV4HINT || v.Key() == dns.SVCB_IPV6HINT {
					addrs = append(addrs, tools.Map(strings.Split(v.String(), ","), M.ParseAddr)...)
				}
			}
		}
	}

	return addrs
}

// 从缓存中获取 DNS 响应并更新 TTL
func (r *Resolver) loadResponse(question dns.Question, transport Transport) (*dns.Msg, int) {
	var (
		response *dns.Msg
		ok       bool
	)
	if r.disableExpire {
		if r.sharedCache {
			response, ok = r.cache.Get(question)
		} else {
			response, ok = r.transportCache.Get(transportCacheKey{
				Question:      question,
				transportName: transport.Name(),
			})
		}
		if !ok {
			return nil, 0
		}
		return response.Copy(), 0
	} else {
		var expireAt time.Time
		if r.sharedCache {
			response, expireAt, ok = r.cache.GetWithLifetime(question)
		} else {
			response, expireAt, ok = r.transportCache.GetWithLifetime(transportCacheKey{
				Question:      question,
				transportName: transport.Name(),
			})
		}

		if !ok {
			return nil, 0
		}

		timeNow := time.Now()
		if timeNow.After(expireAt) {
			if r.sharedCache {
				r.cache.Remove(question)
			} else {
				r.transportCache.Remove(transportCacheKey{
					Question:      question,
					transportName: transport.Name(),
				})
			}
			return nil, 0
		}

		var originalTTL int
		// find min ttl
		for _, recordList := range [][]dns.RR{response.Answer, response.Ns, response.Extra} {
			for _, record := range recordList {
				if originalTTL == 0 || record.Header().Ttl > 0 && int(record.Header().Ttl) < originalTTL {
					originalTTL = int(record.Header().Ttl)
				}
			}
		}

		// remainingTTL = expireAt - timeNow
		remainingTT := int(expireAt.Sub(timeNow).Seconds())
		if remainingTT < 0 {
			remainingTT = 0
		}

		response = response.Copy()
		if originalTTL > 0 {
			duration := uint32(originalTTL - remainingTT)
			for _, recordList := range [][]dns.RR{response.Answer, response.Ns, response.Extra} {
				for _, record := range recordList {
					record.Header().Ttl = record.Header().Ttl - duration
				}
			}
		} else {
			for _, recordList := range [][]dns.RR{response.Answer, response.Ns, response.Extra} {
				for _, record := range recordList {
					record.Header().Ttl = uint32(remainingTT)
				}
			}
		}
		return response, remainingTT
	}
}

// storeCache stores a DNS response message in the cache with an optional TTL (Time-To-Live).
func (r *Resolver) storeCache(transport Transport, question dns.Question, msg *dns.Msg, timeToLive uint32) {
	if timeToLive == 0 {
		return
	}

	if r.disableExpire {
		if r.sharedCache {
			r.cache.Add(question, msg)
		} else {
			r.transportCache.Add(transportCacheKey{
				Question:      question,
				transportName: transport.Name(),
			}, msg)
		}
		return
	}

	if r.sharedCache {
		r.cache.AddWithLifetime(question, msg, time.Second*time.Duration(timeToLive))
	} else {
		r.transportCache.AddWithLifetime(transportCacheKey{
			Question:      question,
			transportName: transport.Name(),
		}, msg, time.Second*time.Duration(timeToLive))
	}
}

func (r *Resolver) ClearCache() {
	if r.cache != nil {
		r.cache.Purge()
	}

	if r.transportCache != nil {
		r.transportCache.Purge()
	}
}

func sortAddrs(response4 []netip.Addr, response6 []netip.Addr, strategy unreal.DNSUnrealStrategy) []netip.Addr {
	if strategy == unreal.DNSUnrealStrategyPreferIPv6 {
		return append(response6, response4...)
	} else {
		return append(response4, response6...)
	}
}

type transportKey struct{}

// contextWithTransportName returns a new context derived from ctx
// that carries the given transportName value.
//
// It stores the name under an unexported key type to avoid collisions with other context keys.
func contextWithTransportName(ctx context.Context, transportName string) context.Context {
	return context.WithValue(ctx, transportKey{}, transportName)
}

func transportNameFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(transportKey{}).(string)
	return value, ok
}
