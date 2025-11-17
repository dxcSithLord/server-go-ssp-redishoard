package redishoard

import "errors"

// Security-related errors
var (
	// ErrNutTooShort indicates the nut is below minimum length
	ErrNutTooShort = errors.New("nut too short: minimum 8 characters")

	// ErrNutTooLong indicates the nut exceeds maximum length
	ErrNutTooLong = errors.New("nut too long: maximum 64 characters")

	// ErrNutInvalidChars indicates the nut contains invalid characters
	ErrNutInvalidChars = errors.New("nut contains invalid characters: must be base64url safe")

	// ErrInvalidHoardCache indicates the cached data is malformed
	ErrInvalidHoardCache = errors.New("invalid hoard cache data")
)
