package rules_test

import (
	"testing"

	"github.com/GeorgeTyupin/go-log-linter/internal/analyzer/rules"
)

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{"starts with lowercase", "starting server", false},
		{"starts with uppercase", "Starting server", true},
		{"empty message", "", false},
		{"number at start", "404 not found", false},
		{"uppercase error", "Error connecting", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, violated := rules.CheckLowercase(tt.msg)
			if violated != tt.wantErr {
				t.Errorf("CheckLowercase(%q) violated=%v, want %v", tt.msg, violated, tt.wantErr)
			}
		})
	}
}

func TestCheckEnglishOnly(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{"english only", "starting server", false},
		{"cyrillic", "запуск сервера", true},
		{"mixed", "starting сервер", true},
		{"numbers and symbols", "server started on port 8080", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, violated := rules.CheckEnglishOnly(tt.msg)
			if violated != tt.wantErr {
				t.Errorf("CheckEnglishOnly(%q) violated=%v, want %v", tt.msg, violated, tt.wantErr)
			}
		})
	}
}

func TestCheckNoSpecialChars(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{"clean message", "server started", false},
		{"with exclamation", "server started!", true},
		{"with emoji", "server started 🚀", true},
		{"with multiple exclamations", "connection failed!!!", true},
		{"with ellipsis", "warning: something went wrong...", true},
		{"with hyphen allowed", "request-id not found", false},
		{"with colon allowed", "status: ok", false},
		{"with slash allowed", "path /api/v1 called", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, violated := rules.CheckNoSpecialChars(tt.msg)
			if violated != tt.wantErr {
				t.Errorf("CheckNoSpecialChars(%q) violated=%v, want %v", tt.msg, violated, tt.wantErr)
			}
		})
	}
}

func TestCheckNoSensitiveData(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{"safe message", "user authenticated successfully", false},
		{"contains password", "user password: 123", true},
		{"contains token", "token: abc123", true},
		{"contains api_key", "api_key=secret", true},
		{"contains secret", "client secret exposed", true},
		{"uppercase keyword", "User PASSWORD reset", true},
		{"safe with user", "user logged in", false},
		{"safe api message", "api request completed", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, violated := rules.CheckNoSensitiveData(tt.msg)
			if violated != tt.wantErr {
				t.Errorf("CheckNoSensitiveData(%q) violated=%v, want %v", tt.msg, violated, tt.wantErr)
			}
		})
	}
}
