package format

import (
	"guard/bridge/tools"
	"strconv"
)

type Stringer interface {
	String() string
}

// ToString converts multiple input parameters into a single string.
// This function supports various types of parameters and converts them accordingly:
// - If the parameter is nil, it is converted to the string "nil"
// - It supports string, boolean, integer types, unsigned integer types, error types, and types implementing the Stringer interface
// - If an unknown type is encountered, the function will panic.
//
// Parameters:
//   - msgs: One or more parameters of any type.
//
// Returns:
//   - A string concatenation of the string representations of all parameters.
//
// Example usage:
//
//	fmt.Println(ToString(1, "hello", true, nil)) // Output: "1hellotrue(nil)"
//	fmt.Println(ToString(123, 3.14, "world")) // Output: "1233.14world"
func ToString(msgs ...any) string {
	var output string
	for _, rawMsg := range msgs {
		if rawMsg == nil {
			output += "nil"
			continue
		}
		switch msg := rawMsg.(type) {
		case string:
			output += msg
		case bool:
			if msg {
				output += "true"
			} else {
				output += "false"
			}
		case uint:
			output += strconv.FormatUint(uint64(msg), 10)
		case uint8:
			output += strconv.FormatUint(uint64(msg), 10)
		case uint16:
			output += strconv.FormatUint(uint64(msg), 10)
		case uint32:
			output += strconv.FormatUint(uint64(msg), 10)
		case uint64:
			output += strconv.FormatUint(msg, 10)
		case int:
			output += strconv.FormatInt(int64(msg), 10)
		case int8:
			output += strconv.FormatInt(int64(msg), 10)
		case int16:
			output += strconv.FormatInt(int64(msg), 10)
		case int32:
			output += strconv.FormatInt(int64(msg), 10)
		case int64:
			output += strconv.FormatInt(msg, 10)
		case uintptr:
			output += strconv.FormatUint(uint64(msg), 10)
		case error:
			output += msg.Error()
		case Stringer:
			output += msg.String()
		default:
			panic("unknown value")
		}
	}
	return output
}

func ToString0[T any](message T) string {
	return ToString(message)
}

func MapToString[T any](arr []T) []string {
	return tools.Map(arr, ToString0[T])
}

// func MapToString2(slice []any) []string {
// 	return common.Map(slice, ToString)
// }

func Seconds(seconds float64) string {
	seconds100 := int(seconds * 100)
	return ToString(seconds100/100, ".", seconds100%100, seconds100%10)
}
