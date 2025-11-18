package redishoard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	ssp "github.com/sqrldev/server-go-ssp"
)

// Hoard implements an ssp.Hoard using redis as
// a backing store with secure memory handling.
type Hoard struct {
	client redis.UniversalClient
}

// NewHoard creates a redis backed Hoard.
func NewHoard(client redis.UniversalClient) *Hoard {
	return &Hoard{
		client: client,
	}
}

// Get implements ssp.Hoard with secure memory handling.
// Retrieves the HoardCache for the given nut from Redis.
func (h *Hoard) Get(nut ssp.Nut) (*ssp.HoardCache, error) {
	// Validate nut input to prevent injection attacks
	if err := ValidateNut(string(nut)); err != nil {
		return nil, fmt.Errorf("invalid nut: %w", err)
	}

	ctx := context.Background()
	data, err := h.client.Get(ctx, string(nut)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ssp.ErrNotFound
		}
		return nil, fmt.Errorf("redis nut lookup failed: %w", err)
	}

	// SECURITY: Ensure sensitive data is cleared from memory after use
	defer ClearBytes(data)

	return h.fromBytes(data)
}

// GetAndDelete implements ssp.Hoard with atomic operations and secure memory handling.
// Atomically retrieves and deletes the HoardCache for the given nut.
// This ensures the nut can only be used once, preventing replay attacks.
func (h *Hoard) GetAndDelete(nut ssp.Nut) (*ssp.HoardCache, error) {
	// Validate nut input
	if err := ValidateNut(string(nut)); err != nil {
		return nil, fmt.Errorf("invalid nut: %w", err)
	}

	ctx := context.Background()

	// Use transaction to ensure atomic get-and-delete
	ret, err := h.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Get(ctx, string(nut))
		pipe.Del(ctx, string(nut))
		return nil
	})

	if err != nil && !errors.Is(err, redis.Nil) {
		// Check if any command in the pipeline failed
		for _, cmd := range ret {
			if cmdErr := cmd.Err(); cmdErr != nil && !errors.Is(cmdErr, redis.Nil) {
				return nil, fmt.Errorf("redis transaction failed: %w", cmdErr)
			}
		}
	}

	// Check for empty pipeline result
	if len(ret) == 0 {
		return nil, ssp.ErrNotFound
	}

	// Process the Get command result
	stringCmd, ok := ret[0].(*redis.StringCmd)
	if !ok {
		return nil, fmt.Errorf("unexpected redis command type in pipeline")
	}

	// Check for not found error
	if cmdErr := stringCmd.Err(); cmdErr != nil {
		if errors.Is(cmdErr, redis.Nil) {
			return nil, ssp.ErrNotFound
		}
		return nil, fmt.Errorf("redis nut lookup failed: %w", cmdErr)
	}

	data, err := stringCmd.Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis HoardCache read failed: %w", err)
	}

	// SECURITY: Clear sensitive data from memory after use
	defer ClearBytes(data)

	// NOTE: Removed sensitive data logging (CWE-200 fix)
	// Previously: log.Printf("data: %v", string(data))

	return h.fromBytes(data)
}

// fromBytes deserializes JSON data into a HoardCache object.
// Internal helper function with secure memory considerations.
func (h *Hoard) fromBytes(data []byte) (*ssp.HoardCache, error) {
	hoardCache := &ssp.HoardCache{}
	err := json.Unmarshal(data, hoardCache)
	if err != nil {
		// SECURITY: Don't include nut value in error message (could be sensitive)
		return nil, fmt.Errorf("failed to decode HoardCache: %w", err)
	}
	return hoardCache, nil
}

// Save implements ssp.Hoard with secure memory handling.
// Stores the HoardCache in Redis with the specified TTL.
func (h *Hoard) Save(nut ssp.Nut, value *ssp.HoardCache, expiration time.Duration) error {
	// Validate nut input
	if err := ValidateNut(string(nut)); err != nil {
		return fmt.Errorf("invalid nut: %w", err)
	}

	// Validate cache object
	if value == nil {
		return fmt.Errorf("cannot save nil HoardCache")
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed json encoding HoardCache: %w", err)
	}

	// SECURITY: Clear serialized sensitive data from memory after use
	defer ClearBytes(jsonBytes)

	ctx := context.Background()
	return h.client.Set(ctx, string(nut), jsonBytes, expiration).Err()
}
