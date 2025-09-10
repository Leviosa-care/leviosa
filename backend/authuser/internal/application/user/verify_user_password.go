package user

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// CompareSecureHashAndValue compares a secure hash with a value
func CompareSecureHashAndValue(ctx context.Context, value string, hashValue string) (bool, error) {
	v := []byte(value)

	// Parse the stored hash to extract parameters, salt, and hash
	parts := strings.Split(hashValue, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return false, fmt.Errorf("invalid hash format")
	}

	// Parse version
	versionPart := parts[2]
	if !strings.HasPrefix(versionPart, "v=") {
		return false, fmt.Errorf("invalid version format")
	}
	version, err := strconv.Atoi(versionPart[2:])
	if err != nil {
		return false, fmt.Errorf("invalid version number: %w", err)
	}
	if version != argon2.Version {
		return false, fmt.Errorf("unsupported Argon2 version")
	}

	// Parse parameters (m=memory,t=iterations,p=parallelism)
	paramsPart := parts[3]
	paramPairs := strings.Split(paramsPart, ",")
	if len(paramPairs) != 3 {
		return false, fmt.Errorf("invalid parameters format")
	}

	var memory, iterations uint32
	var parallelism uint8

	for _, pair := range paramPairs {
		keyValue := strings.Split(pair, "=")
		if len(keyValue) != 2 {
			return false, fmt.Errorf("invalid parameter format")
		}
		value, err := strconv.ParseUint(keyValue[1], 10, 32)
		if err != nil {
			return false, fmt.Errorf("invalid parameter value: %w", err)
		}
		switch keyValue[0] {
		case "m":
			memory = uint32(value)
		case "t":
			iterations = uint32(value)
		case "p":
			parallelism = uint8(value)
		default:
			return false, fmt.Errorf("unknown parameter: %s", keyValue[0])
		}
	}

	// Decode salt and stored hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Combine value with pepper
	// peppered := append(v, h.pepper[:]...)
	peppered := append(v, []byte("testpepper123456testpepper123456")[:]...)

	// Generate hash using the extracted salt and parameters
	computedHash := argon2.IDKey(
		peppered,
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(storedHash)),
	)

	// Compare hashes
	if len(computedHash) != len(storedHash) {
		return false, nil
	}
	for i := range len(computedHash) {
		if computedHash[i] != storedHash[i] {
			return false, nil
		}
	}
	return true, nil
}
