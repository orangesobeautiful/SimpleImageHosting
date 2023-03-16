package utils

import "strings"

// IsTrueVal 返回设置的值是否为真
func IsTrueVal(val string) bool {
	return val == "true" || strings.EqualFold(val, "true")
}
