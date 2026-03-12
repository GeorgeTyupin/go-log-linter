package rules

import "unicode"

func CheckEnglishOnly(msg string) (string, bool) {
	for _, r := range msg {
		if r > unicode.MaxASCII {
			return "log message must be in English only (non-ASCII characters found)", true
		}
	}
	return "", false
}
