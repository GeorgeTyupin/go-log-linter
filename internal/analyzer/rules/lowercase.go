package rules

import "unicode"

func CheckLowercase(msg string) (string, bool) {
	if len(msg) == 0 {
		return "", false
	}

	if unicode.IsUpper([]rune(msg)[0]) {
		return "log message must start with a lowercase letter", true
	}

	return "", false
}
