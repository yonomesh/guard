package unreal

type DNSUnrealStrategy = uint8

const (
	DNSUnrealStrategyAsIs DNSUnrealStrategy = iota
	DNSUnrealStrategyPreferIPv4
	DNSUnrealStrategyPreferIPv6
	DNSUnrealStrategyOnlyIPv4
	DNSUnrealStrategyOnlyIPv6
)
