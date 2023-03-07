package util

import (
	"bufio"
	"github.com/clevyr/yampl/cmd"
	"strings"
)

func FixStringToStringNewlines(s []string) []string {
	var prevValueFlag string
	result := make([]string, 0, len(s))
	for i, arg := range s {
		switch {
		case arg == "--":
			if prevValueFlag != "" {
				result = append(result, prevValueFlag)
			}
			result = append(result, s[i:]...)
			return result
		case prevValueFlag != "":
			result = append(result, fixArgNewlines(prevValueFlag+"="+arg)...)
			prevValueFlag = ""
		case hasValueFlag(arg):
			if strings.ContainsRune(arg, '=') {
				result = append(result, fixArgNewlines(arg)...)
			} else {
				prevValueFlag = arg
			}
		default:
			result = append(result, arg)
		}
	}
	return result
}

func hasValueFlag(s string) bool {
	return s == "-"+cmd.ValueFlagShort ||
		s == "--"+cmd.ValueFlag ||
		strings.HasPrefix(s, "-"+cmd.ValueFlagShort+"=") ||
		strings.HasPrefix(s, "--"+cmd.ValueFlag+"=")
}

func fixArgNewlines(arg string) []string {
	if strings.ContainsRune(arg, '\n') {
		prefix, arg, found := strings.Cut(arg, "=")
		if !found {
			return []string{prefix}
		}

		result := make([]string, 0, 2)
		s := bufio.NewScanner(strings.NewReader(arg))
		for s.Scan() {
			if len(s.Bytes()) > 0 {
				result = append(result, prefix+"="+s.Text())
			}
		}
		return result
	} else {
		return []string{arg}
	}
}
