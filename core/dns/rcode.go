package dns

import F "uni/bridge/tools/format"

type RCodeError uint16

const (
	RCodeSuccess        RCodeError = iota // NoError
	RCodeFormatError                      // FormErr
	RCodeServerFailure                    // ServFail
	RCodeNameError                        // NXDomain
	RCodeNotImplemented                   // NotImp
	RCodeRefused                          // Refused
)

func (e RCodeError) Error() string {
	switch e {
	case RCodeSuccess:
		return "success"
	case RCodeFormatError:
		return "format error"
	case RCodeServerFailure:
		return "server failure"
	case RCodeNameError:
		return "name error"
	case RCodeNotImplemented:
		return "not implemented"
	case RCodeRefused:
		return "refused"
	default:
		return F.ToString("unknown error: ", uint16(e))
	}
}
