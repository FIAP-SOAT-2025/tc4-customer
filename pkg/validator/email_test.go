package validator

import "testing"

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid email", "test@example.com", true},
		{"Valid email with subdomain", "user@mail.example.com", true},
		{"Valid email with plus", "user+tag@example.com", true},
		{"Invalid - no @", "testexample.com", false},
		{"Invalid - no domain", "test@", false},
		{"Invalid - no local part", "@example.com", false},
		{"Invalid - multiple @", "test@@example.com", false},
		{"Invalid - spaces", "test @example.com", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%s) = %v; want %v", tt.email, result, tt.expected)
			}
		})
	}
}
