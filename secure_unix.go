//go:build unix

package redishoard

import (
	"golang.org/x/sys/unix"
)

// LockMemory prevents memory from being swapped to disk (Unix/Linux/macOS).
// This helps prevent sensitive data from being written to swap space.
// Requires appropriate privileges (CAP_IPC_LOCK on Linux or root).
func LockMemory(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return unix.Mlock(b)
}

// UnlockMemory unlocks previously locked memory (Unix/Linux/macOS).
func UnlockMemory(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return unix.Munlock(b)
}
