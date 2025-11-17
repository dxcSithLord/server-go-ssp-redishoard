//go:build windows

package redishoard

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// LockMemory prevents memory from being swapped to disk (Windows).
// This helps prevent sensitive data from being written to the page file.
// May require elevated privileges.
func LockMemory(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return windows.VirtualLock(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
}

// UnlockMemory unlocks previously locked memory (Windows).
func UnlockMemory(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return windows.VirtualUnlock(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
}
