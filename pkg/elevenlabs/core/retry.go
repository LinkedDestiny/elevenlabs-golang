package core

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

// RetryConfig configures retry behavior for HTTP requests
type RetryConfig struct {
	MaxAttempts       int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
	JitterFactor      float64
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:       3,
		InitialDelay:      500 * time.Millisecond,
		MaxDelay:          30 * time.Second,
		BackoffMultiplier: 2.0,
		JitterFactor:      0.25,
	}
}

// ShouldRetry determines if a request should be retried based on the status code
func ShouldRetry(statusCode int) bool {
	retryable400s := []int{429, 408, 409}

	// Retry on 5xx server errors
	if statusCode >= 500 {
		return true
	}

	// Retry on specific 4xx errors
	for _, code := range retryable400s {
		if statusCode == code {
			return true
		}
	}

	return false
}

// CalculateDelay calculates the delay before retrying a request
func CalculateDelay(attempt int, retryAfter string) time.Duration {
	// If server provides Retry-After header, respect it (with reasonable limits)
	if retryAfter != "" {
		if seconds, err := strconv.Atoi(retryAfter); err == nil {
			delay := time.Duration(seconds) * time.Second
			if delay <= 30*time.Second { // Cap at 30 seconds
				return delay
			}
		}
	}

	// Use exponential backoff with jitter
	config := DefaultRetryConfig()

	// Calculate exponential backoff
	backoffDelay := float64(config.InitialDelay) * math.Pow(config.BackoffMultiplier, float64(attempt))

	// Cap at maximum delay
	if backoffDelay > float64(config.MaxDelay) {
		backoffDelay = float64(config.MaxDelay)
	}

	// Add jitter to avoid thundering herd
	jitter := backoffDelay * config.JitterFactor * (rand.Float64()*2 - 1) // Random between -jitterFactor and +jitterFactor
	finalDelay := backoffDelay + jitter

	// Ensure delay is non-negative
	if finalDelay < 0 {
		finalDelay = 0
	}

	return time.Duration(finalDelay)
}
