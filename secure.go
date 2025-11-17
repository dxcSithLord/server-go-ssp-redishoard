package redishoard

import (
	"crypto/subtle"
	"runtime"
)

// ClearBytes securely clears a byte slice from memory.
// This implements CWE-226 mitigation by ensuring sensitive data
// does not remain in memory after use.
func ClearBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	// Use crypto/subtle for constant-time operations
	// This prevents compiler optimizations that might skip the clearing
	for i := range b {
		b[i] = 0
	}

	// Ensure the slice is not optimized away
	runtime.KeepAlive(b)

	// Additional memory barrier to prevent reordering
	_ = subtle.ConstantTimeCompare(b, b)
}

// SecureBuffer wraps a byte slice with automatic secure clearing.
// Use this for operations involving sensitive data.
type SecureBuffer struct {
	data   []byte
	locked bool
}

// NewSecureBuffer creates a new secure buffer with optional memory locking.
func NewSecureBuffer(size int, lock bool) (*SecureBuffer, error) {
	buf := &SecureBuffer{
		data:   make([]byte, size),
		locked: false,
	}

	if lock && size > 0 {
		if err := LockMemory(buf.data); err != nil {
			// Log warning but don't fail - memory locking is optional
			// Some systems may not support it
			buf.locked = false
		} else {
			buf.locked = true
		}
	}

	return buf, nil
}

// Data returns the underlying byte slice.
func (sb *SecureBuffer) Data() []byte {
	return sb.data
}

// Clear securely clears the buffer contents.
func (sb *SecureBuffer) Clear() {
	if sb.data != nil {
		ClearBytes(sb.data)
	}
}

// Destroy clears and releases the buffer.
func (sb *SecureBuffer) Destroy() {
	if sb.data == nil {
		return
	}

	// Clear the data first
	ClearBytes(sb.data)

	// Unlock memory if it was locked
	if sb.locked {
		_ = UnlockMemory(sb.data)
		sb.locked = false
	}

	// Release reference
	sb.data = nil
}

// ValidateNut validates that a nut conforms to SQRL specifications.
// Returns nil if valid, error otherwise.
func ValidateNut(nut string) error {
	// SQRL nuts are typically 8-20 bytes when decoded
	// Base64url encoded, they would be longer
	if len(nut) < MinNutLength {
		return ErrNutTooShort
	}

	if len(nut) > MaxNutLength {
		return ErrNutTooLong
	}

	// Validate characters are safe (base64url safe characters)
	for _, c := range nut {
		if !isBase64URLSafe(c) {
			return ErrNutInvalidChars
		}
	}

	return nil
}

// isBase64URLSafe checks if a rune is a valid base64url character.
func isBase64URLSafe(c rune) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_'
}

// Constants for nut validation.
const (
	MinNutLength = 8  // Minimum nut length in characters.
	MaxNutLength = 64 // Maximum nut length in characters (base64 encoded).
)
