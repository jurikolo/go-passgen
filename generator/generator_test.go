package generator

import (
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Length:    16,
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
				Symbols:   true,
				Count:     5,
			},
			wantErr: false,
		},
		{
			name: "length too small",
			config: Config{
				Length:    7,
				Uppercase: true,
			},
			wantErr: true,
		},
		{
			name: "length too large",
			config: Config{
				Length:    129,
				Uppercase: true,
			},
			wantErr: true,
		},
		{
			name: "count too small",
			config: Config{
				Length:    12,
				Uppercase: true,
				Count:     0,
			},
			wantErr: true,
		},
		{
			name: "count too large",
			config: Config{
				Length:    12,
				Uppercase: true,
				Count:     101,
			},
			wantErr: true,
		},
		{
			name: "no character sets",
			config: Config{
				Length: 12,
				Count:  1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	config := Config{
		Length:    20,
		Uppercase: true,
		Lowercase: true,
		Digits:    true,
		Symbols:   true,
		Count:     1,
	}

	password, err := Generate(config)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(password) != config.Length {
		t.Errorf("Expected password length %d, got %d", config.Length, len(password))
	}

	// Check that password contains at least one character from each set
	// (probabilistic, but with all sets enabled it's extremely likely)
	hasUpper, hasLower, hasDigit, hasSymbol := false, false, false, false
	for _, ch := range password {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		default:
			// Check if it's a symbol from our set
			for _, s := range symbols {
				if ch == s {
					hasSymbol = true
					break
				}
			}
		}
	}

	if !hasUpper {
		t.Error("Password missing uppercase letters")
	}
	if !hasLower {
		t.Error("Password missing lowercase letters")
	}
	if !hasDigit {
		t.Error("Password missing digits")
	}
	if !hasSymbol {
		t.Error("Password missing symbols")
	}
}

func TestGenerateMultiple(t *testing.T) {
	config := Config{
		Length:    12,
		Uppercase: true,
		Lowercase: true,
		Digits:    false,
		Symbols:   false,
		Count:     5,
	}

	passwords, err := GenerateMultiple(config)
	if err != nil {
		t.Fatalf("GenerateMultiple failed: %v", err)
	}

	if len(passwords) != config.Count {
		t.Errorf("Expected %d passwords, got %d", config.Count, len(passwords))
	}

	// Ensure all passwords are unique (they should be, but with crypto/rand collisions are extremely unlikely)
	seen := make(map[string]bool)
	for _, pwd := range passwords {
		if len(pwd) != config.Length {
			t.Errorf("Password length mismatch: expected %d, got %d", config.Length, len(pwd))
		}
		if seen[pwd] {
			t.Errorf("Duplicate password generated: %s", pwd)
		}
		seen[pwd] = true
	}
}

func TestGenerateWithOnlyDigits(t *testing.T) {
	config := Config{
		Length:    10,
		Uppercase: false,
		Lowercase: false,
		Digits:    true,
		Symbols:   false,
		Count:     1,
	}

	password, err := Generate(config)
	if err != nil {
		t.Fatalf("Generate with only digits failed: %v", err)
	}

	for _, ch := range password {
		if ch < '0' || ch > '9' {
			t.Errorf("Password contains non-digit character: %c", ch)
		}
	}
}

func TestGenerateEmptyPool(t *testing.T) {
	config := Config{
		Length:    10,
		Uppercase: false,
		Lowercase: false,
		Digits:    false,
		Symbols:   false,
		Count:     1,
	}

	_, err := Generate(config)
	if err == nil {
		t.Error("Expected error for empty character pool")
	}
}
