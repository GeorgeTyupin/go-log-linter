package rules

import "unicode"

var allowedPunctuation = map[rune]bool{
	' ': true, '-': true, '_': true, '/': true,
	',': true, ':': true,
	'(': true, ')': true, '[': true, ']': true,
	'\'': true,
}

func CheckNoSpecialChars(msg string) (string, bool) {
	for _, r := range msg {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}

		if allowedPunctuation[r] {
			continue
		}

		if isEmoji(r) {
			return "log message must not contain emojis", true
		}

		return "log message must not contain special characters", true
	}

	return "", false
}

func isEmoji(r rune) bool {
	return r > 0x00FF && (
		(r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x1F300 && r <= 0x1F5FF) ||
		(r >= 0x1F680 && r <= 0x1F6FF) ||
		(r >= 0x1F1E0 && r <= 0x1F1FF) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF) ||
		(r >= 0xFE00 && r <= 0xFE0F) ||
		(r >= 0x1F900 && r <= 0x1F9FF) ||
		(r >= 0x1FA00 && r <= 0x1FA6F) ||
		(r >= 0x1FA70 && r <= 0x1FAFF))
}
