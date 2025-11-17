package redishoard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	ssp "github.com/sqrldev/server-go-ssp"
)

const testStateValue = "boom!"

// Redis must be installed locally on the default port to run these tests.
func TestSave(t *testing.T) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{})
	defer client.Close()

	h := NewHoard(client)

	// Use valid base64url safe nut
	var nut ssp.Nut = "test-nut-save-1234"
	hoardCache := &ssp.HoardCache{
		State: testStateValue,
	}
	err := h.Save(nut, hoardCache, time.Second*10)
	if err != nil {
		t.Fatalf("Failed saving hoard cache: %v", err)
	}

	val, err := h.Get(nut)
	if err != nil {
		t.Fatalf("Failed getting from Hoard: %v", err)
	}

	if val.State != testStateValue {
		t.Fatalf("Wrong value from Hoard: %#v", val)
	}

	// Cleanup
	client.Del(context.Background(), string(nut))
}

func TestGetAndDelete(t *testing.T) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{})
	defer client.Close()

	h := NewHoard(client)

	// Use valid base64url safe nut
	var nut ssp.Nut = "test-get-and-delete"
	hoardCache := &ssp.HoardCache{
		State: testStateValue,
	}
	err := h.Save(nut, hoardCache, time.Second*10)
	if err != nil {
		t.Fatalf("Failed saving hoard cache: %v", err)
	}

	val, err := h.Get(nut)
	if err != nil {
		t.Fatalf("Failed Get from Hoard: %v", err)
	}

	if val.State != testStateValue {
		t.Fatalf("Wrong value from Hoard: %#v", val)
	}

	val, err = h.GetAndDelete(nut)
	if err != nil {
		t.Fatalf("Failed GetAndDelete from Hoard: %v", err)
	}

	if val.State != testStateValue {
		t.Fatalf("Wrong value from Hoard: %#v", val)
	}

	val, err = h.GetAndDelete(nut)
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("Should have been deleted but wasn't: %v", err)
	}

	if val != nil {
		t.Fatalf("Wrong value from Hoard: %#v", val)
	}
}

func TestNutValidation(t *testing.T) {
	tests := []struct {
		name    string
		nut     string
		wantErr bool
	}{
		{
			name:    "valid nut",
			nut:     "valid-nut-12345",
			wantErr: false,
		},
		{
			name:    "nut too short",
			nut:     "short",
			wantErr: true,
		},
		{
			name:    "nut at minimum length",
			nut:     "12345678",
			wantErr: false, // exactly 8 characters, all valid
		},
		{
			name:    "nut with invalid characters",
			nut:     "nut@with#invalid!",
			wantErr: true,
		},
		{
			name:    "valid base64url nut",
			nut:     "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNut(tt.nut)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNut(%q) error = %v, wantErr %v", tt.nut, err, tt.wantErr)
			}
		})
	}
}

func TestSaveNilCache(t *testing.T) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{})
	defer client.Close()

	h := NewHoard(client)
	var nut ssp.Nut = "test-nil-cache"

	err := h.Save(nut, nil, time.Second)
	if err == nil {
		t.Fatal("Expected error when saving nil cache, got nil")
	}
}

func TestGetNonExistentNut(t *testing.T) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{})
	defer client.Close()

	h := NewHoard(client)
	var nut ssp.Nut = "nonexistent-nut-xyz"

	_, err := h.Get(nut)
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("Expected ErrNotFound, got: %v", err)
	}
}

func TestInvalidNutCharacters(t *testing.T) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{})
	defer client.Close()

	h := NewHoard(client)

	// Test with invalid characters
	invalidNuts := []ssp.Nut{
		"nut@with!special",
		"nut with spaces",
		"nut/with/slashes",
		"nut+with+plus",
	}

	for _, nut := range invalidNuts {
		_, err := h.Get(nut)
		if err == nil {
			t.Errorf("Expected error for invalid nut %q, got nil", nut)
		}
	}
}
