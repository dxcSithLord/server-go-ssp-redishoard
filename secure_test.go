package redishoard

import (
	"bytes"
	"errors"
	"testing"
)

func TestClearBytes(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "sensitive password",
			data: []byte("super-secret-password-123!@#"),
		},
		{
			name: "cryptographic key",
			data: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		},
		{
			name: "JSON with tokens",
			data: []byte(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9","secret":"confidential"}`),
		},
		{
			name: "empty slice",
			data: []byte{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Keep a copy of original length
			originalLen := len(tc.data)

			// Clear the bytes
			ClearBytes(tc.data)

			// Verify all bytes are zeroed
			for i, b := range tc.data {
				if b != 0 {
					t.Errorf("byte %d not cleared: got %d, want 0", i, b)
				}
			}

			// Verify length unchanged
			if len(tc.data) != originalLen {
				t.Errorf("slice length changed: got %d, want %d", len(tc.data), originalLen)
			}
		})
	}
}

func TestSecureBuffer(t *testing.T) {
	t.Run("create and destroy buffer", func(t *testing.T) {
		buf, err := NewSecureBuffer(32, false)
		if err != nil {
			t.Fatalf("failed to create secure buffer: %v", err)
		}

		// Write sensitive data
		sensitiveData := []byte("sensitive-authentication-token!")
		copy(buf.Data(), sensitiveData)

		// Verify data is stored
		if !bytes.Equal(buf.Data()[:len(sensitiveData)], sensitiveData) {
			t.Error("data not properly stored in buffer")
		}

		// Destroy the buffer
		buf.Destroy()

		// Buffer should be nil after destroy
		if buf.Data() != nil {
			t.Error("buffer not nil after destroy")
		}
	})

	t.Run("clear buffer contents", func(t *testing.T) {
		buf, err := NewSecureBuffer(16, false)
		if err != nil {
			t.Fatalf("failed to create secure buffer: %v", err)
		}

		// Fill with data
		for i := range buf.Data() {
			buf.Data()[i] = byte(i + 1)
		}

		// Clear
		buf.Clear()

		// Verify cleared
		for i, b := range buf.Data() {
			if b != 0 {
				t.Errorf("byte %d not cleared: got %d", i, b)
			}
		}

		buf.Destroy()
	})

	t.Run("zero-size buffer", func(t *testing.T) {
		buf, err := NewSecureBuffer(0, false)
		if err != nil {
			t.Fatalf("failed to create zero-size buffer: %v", err)
		}

		if len(buf.Data()) != 0 {
			t.Errorf("expected empty buffer, got length %d", len(buf.Data()))
		}

		buf.Destroy()
	})
}

func TestValidateNut(t *testing.T) {
	tests := []struct {
		name        string
		nut         string
		expectError error
	}{
		{
			name:        "valid nut - alphanumeric",
			nut:         "abcdefghij123456",
			expectError: nil,
		},
		{
			name:        "valid nut - with dash and underscore",
			nut:         "test-nut_value_123",
			expectError: nil,
		},
		{
			name:        "valid nut - uppercase",
			nut:         "ABCDEFGHIJKLMNOP",
			expectError: nil,
		},
		{
			name:        "too short - 7 chars",
			nut:         "abcdefg",
			expectError: ErrNutTooShort,
		},
		{
			name:        "minimum length - 8 chars",
			nut:         "abcdefgh",
			expectError: nil,
		},
		{
			name:        "invalid char - space",
			nut:         "test nut value",
			expectError: ErrNutInvalidChars,
		},
		{
			name:        "invalid char - at symbol",
			nut:         "test@nut#value",
			expectError: ErrNutInvalidChars,
		},
		{
			name:        "invalid char - slash",
			nut:         "test/nut/value",
			expectError: ErrNutInvalidChars,
		},
		{
			name:        "invalid char - plus",
			nut:         "test+nut+value",
			expectError: ErrNutInvalidChars,
		},
		{
			name:        "invalid char - equals",
			nut:         "testnutvalue==",
			expectError: ErrNutInvalidChars,
		},
		{
			name:        "maximum length - 64 chars",
			nut:         "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_",
			expectError: nil,
		},
		{
			name:        "too long - 65 chars",
			nut:         "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_X",
			expectError: ErrNutTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNut(tt.nut)
			if tt.expectError == nil {
				if err != nil {
					t.Errorf("ValidateNut(%q) = %v, want nil", tt.nut, err)
				}
			} else if !errors.Is(err, tt.expectError) {
				t.Errorf("ValidateNut(%q) = %v, want %v", tt.nut, err, tt.expectError)
			}
		})
	}
}

func TestIsBase64URLSafe(t *testing.T) {
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	invalidChars := "@#$%^&*()+=!~`[]{}|\\:;\"'<>,./? "

	for _, c := range validChars {
		if !isBase64URLSafe(c) {
			t.Errorf("isBase64URLSafe(%q) = false, want true", c)
		}
	}

	for _, c := range invalidChars {
		if isBase64URLSafe(c) {
			t.Errorf("isBase64URLSafe(%q) = true, want false", c)
		}
	}
}

// BenchmarkClearBytes measures performance of secure clearing.
func BenchmarkClearBytes(b *testing.B) {
	sizes := []int{16, 64, 256, 1024, 4096}

	for _, size := range sizes {
		b.Run(string(rune(size)), func(b *testing.B) {
			data := make([]byte, size)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				ClearBytes(data)
			}
		})
	}
}

// BenchmarkValidateNut measures validation performance.
func BenchmarkValidateNut(b *testing.B) {
	nut := "test-nut-value-1234567890"
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ValidateNut(nut)
	}
}
