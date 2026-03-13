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

// CheckNoSensitiveData проверяет сообщение на наличие чувствительных данных.
func CheckNoSensitiveData(msg string, extraPatterns []string) (string, bool) {
	lower := strings.ToLower(msg)

	// Проверяем встроенные substring-паттерны
	for _, keyword := range substringKeywords {
		if strings.Contains(lower, keyword) {
			return fmt.Sprintf("log message may contain sensitive data (keyword: %q)", keyword), true
		}
	}

	// Проверяем кастомные паттерны как подстроки
	for _, pattern := range extraPatterns {
		p := strings.ToLower(strings.TrimSpace(pattern))
		if p != "" && strings.Contains(lower, p) {
			return fmt.Sprintf("log message may contain sensitive data (custom pattern: %q)", p), true
		}
	}

	// Разбиваем на слова и ищем точные совпадения со встроенными exact-словами
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
