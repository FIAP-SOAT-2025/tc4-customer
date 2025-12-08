package validator

import "testing"

func TestCleanCPF(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"With dots and dash", "123.456.789-00", "12345678900"},
		{"Only numbers", "12345678900", "12345678900"},
		{"With spaces", "123 456 789 00", "12345678900"},
		{"Mixed", "123.456 789-00", "12345678900"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanCPF(tt.input)
			if result != tt.expected {
				t.Errorf("CleanCPF(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"Valid CPF 1", "11144477735", true},
		{"Valid CPF 2", "12345678909", true},
		{"Valid CPF formatted", "111.444.777-35", true},
		{"Invalid - too short", "123456789", false},
		{"Invalid - too long", "123456789012", false},
		{"Invalid - all same digits", "11111111111", false},
		{"Invalid - all zeros", "00000000000", false},
		{"Invalid - wrong check digit", "12345678901", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCPF(tt.cpf)
			if result != tt.expected {
				t.Errorf("IsValidCPF(%s) = %v; want %v", tt.cpf, result, tt.expected)
			}
		})
	}
}
