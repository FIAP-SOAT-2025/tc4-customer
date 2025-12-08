package validator

import (
	"regexp"
)

// CleanCPF removes all non-numeric characters from CPF
func CleanCPF(cpf string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(cpf, "")
}

// IsValidCPF validates a Brazilian CPF
func IsValidCPF(cpf string) bool {
	cpf = CleanCPF(cpf)

	// CPF must have exactly 11 digits
	if len(cpf) != 11 {
		return false
	}

	// Check if all digits are the same (invalid CPFs)
	allSame := true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	// Validate first check digit
	sum := 0
	for i := 0; i < 9; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	firstCheckDigit := 0
	if remainder >= 2 {
		firstCheckDigit = 11 - remainder
	}
	if int(cpf[9]-'0') != firstCheckDigit {
		return false
	}

	// Validate second check digit
	sum = 0
	for i := 0; i < 10; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	secondCheckDigit := 0
	if remainder >= 2 {
		secondCheckDigit = 11 - remainder
	}
	if int(cpf[10]-'0') != secondCheckDigit {
		return false
	}

	return true
}
