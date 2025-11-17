//go:build !unix && !windows

package redishoard

// LockMemory is a no-op on unsupported platforms.
// Memory locking is not available on this platform.
func LockMemory(b []byte) error {
	return nil
}

// UnlockMemory is a no-op on unsupported platforms.
func UnlockMemory(b []byte) error {
	return nil
}
