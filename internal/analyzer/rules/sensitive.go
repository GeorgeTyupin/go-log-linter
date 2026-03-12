package rules

import (
	"fmt"
	"strings"
	"unicode"
)

var exactKeywords = []string{
	"password", "passwd", "pass",
	"secret", "token",
	"credential", "credentials",
	"authorization",
}

var substringKeywords = []string{
	"api_key", "apikey", "api-key",
	"private_key", "privatekey",
	"access_key", "accesskey",
	"client_secret",
}

func CheckNoSensitiveData(msg string) (string, bool) {
	lower := strings.ToLower(msg)

	for _, keyword := range substringKeywords {
		if strings.Contains(lower, keyword) {
			return fmt.Sprintf("log message may contain sensitive data (keyword: %q)", keyword), true
		}
	}

	words := strings.FieldsFunc(lower, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for _, word := range words {
		for _, keyword := range exactKeywords {
			if word == keyword {
				return fmt.Sprintf("log message may contain sensitive data (keyword: %q)", keyword), true
			}
		}
	}

	return "", false
}
