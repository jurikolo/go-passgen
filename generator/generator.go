package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	symbols          = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// Config holds password generation configuration
type Config struct {
	Length    int  `json:"length"`
	Uppercase bool `json:"uppercase"`
	Lowercase bool `json:"lowercase"`
	Digits    bool `json:"digits"`
	Symbols   bool `json:"symbols"`
	Count     int  `json:"count"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Length < 8 || c.Length > 128 {
		return errors.New("length must be between 8 and 128")
	}
	if c.Count < 1 || c.Count > 100 {
		return errors.New("count must be between 1 and 100")
	}
	if !c.Uppercase && !c.Lowercase && !c.Digits && !c.Symbols {
		return errors.New("at least one character set must be enabled")
	}
	return nil
}

// Generate creates a single password based on the configuration
func Generate(config Config) (string, error) {
	if err := config.Validate(); err != nil {
		return "", err
	}

	// Build character pool
	var pool []rune
	if config.Lowercase {
		pool = append(pool, []rune(lowercaseLetters)...)
	}
	if config.Uppercase {
		pool = append(pool, []rune(uppercaseLetters)...)
	}
	if config.Digits {
		pool = append(pool, []rune(digits)...)
	}
	if config.Symbols {
		pool = append(pool, []rune(symbols)...)
	}

	if len(pool) == 0 {
		return "", errors.New("character pool is empty")
	}

	// Generate password
	result := make([]rune, config.Length)
	for i := 0; i < config.Length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
		if err != nil {
			return "", err
		}
		result[i] = pool[n.Int64()]
	}

	return string(result), nil
}

// GenerateMultiple creates multiple passwords
func GenerateMultiple(config Config) ([]string, error) {
	passwords := make([]string, config.Count)
	for i := 0; i < config.Count; i++ {
		pwd, err := Generate(config)
		if err != nil {
			return nil, err
		}
		passwords[i] = pwd
	}
	return passwords, nil
}
