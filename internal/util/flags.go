package util

import (
	"github.com/clevyr/yampl/cmd"
	"regexp"
	"strings"
)

func FixStringToStringNewlines(s []string) []string {
	var replaceNext bool
	for i, arg := range s {
		switch {
		case arg == "--":
			return s
		case hasValueFlag(arg) || replaceNext:
			replaceNext = false
			if strings.ContainsRune(arg, '=') {
				if strings.ContainsRune(arg, '\n') {
					arg = regexp.MustCompile("\n+").ReplaceAllString(arg, ",")
					arg = strings.Trim(arg, ",")
					arg = strings.Replace(arg, "=,", "=", 1)
					s[i] = arg
				}
			} else {
				replaceNext = true
			}
		}
	}
	return s
}

func hasValueFlag(s string) bool {
	return s == "-"+cmd.ValueFlagShort ||
		s == "--"+cmd.ValueFlag ||
		strings.HasPrefix(s, "-"+cmd.ValueFlagShort+"=") ||
		strings.HasPrefix(s, "--"+cmd.ValueFlag+"=")
}
