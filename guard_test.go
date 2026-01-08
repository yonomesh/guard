package uni

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["娜扎"] = 60
	for k, v := range scoreMap {
		fmt.Println(k, v)
	}

	for k := range scoreMap {
		fmt.Println(k)
	}
}
