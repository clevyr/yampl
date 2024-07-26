package util

import "strings"

func UpperFirst(s string) string {
	switch len(s) {
	case 0:
		return s
	case 1:
		return strings.ToUpper(s)
	default:
		return strings.ToUpper(s[0:1]) + s[1:]
	}
}

func LowerFirst(s string) string {
	switch len(s) {
	case 0:
		return s
	case 1:
		return strings.ToLower(s)
	default:
		return strings.ToLower(s[0:1]) + s[1:]
	}
}
