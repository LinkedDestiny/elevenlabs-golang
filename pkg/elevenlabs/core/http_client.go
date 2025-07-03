package core

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

// Environment represents an ElevenLabs API environment
type Environment struct {
	BaseURL      string
	WebSocketURL string
}

// HTTPClient handles all HTTP communication with the ElevenLabs API
type HTTPClient struct {
	httpClient  *http.Client
	baseURL     string
	wsURL       string
	apiKey      string
	userAgent   string
	timeout     time.Duration
	retryConfig RetryConfig
}

// Config represents HTTP client configuration
type Config struct {
	APIKey      string
	Environment Environment
	Timeout     time.Duration
	HTTPClient  *http.Client
	UserAgent   string
	RetryConfig RetryConfig
}

// NewHTTPClient creates a new HTTP client with the specified configuration
func NewHTTPClient(config Config) *HTTPClient {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return &HTTPClient{
		httpClient:  config.HTTPClient,
		baseURL:     config.Environment.BaseURL,
		wsURL:       config.Environment.WebSocketURL,
		apiKey:      config.APIKey,
		userAgent:   config.UserAgent,
		timeout:     config.Timeout,
		retryConfig: config.RetryConfig,
	}
}

// GetWebSocketURL returns the WebSocket base URL
func (c *HTTPClient) GetWebSocketURL() string {
	return c.wsURL
}

// GetAPIKey returns the API key
func (c *HTTPClient) GetAPIKey() string {
	return c.apiKey
}

// Request makes an HTTP request with the specified method, path, body and headers
func (c *HTTPClient) Request(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + "/" + strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Add authentication headers
	AddAuthHeaders(req, c.apiKey)

	// Add user agent
	req.Header.Set("User-Agent", c.userAgent)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make request with retry logic
	return c.RequestWithRetry(ctx, req)
}

// RequestWithRetry executes the HTTP request with retry logic
func (c *HTTPClient) RequestWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retryConfig.MaxAttempts; attempt++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < c.retryConfig.MaxAttempts {
				delay := CalculateDelay(attempt, "")
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
					continue
				}
			}
			continue
		}

		// Check if we should retry based on status code
		if ShouldRetry(resp.StatusCode) && attempt < c.retryConfig.MaxAttempts {
			resp.Body.Close()
			retryAfter := resp.Header.Get("Retry-After")
			delay := CalculateDelay(attempt, retryAfter)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
				continue
			}
		}

		return resp, nil
	}

	return nil, lastErr
}

// Stream makes a streaming HTTP request
func (c *HTTPClient) Stream(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + "/" + strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Add authentication headers
	AddAuthHeaders(req, c.apiKey)

	// Add user agent
	req.Header.Set("User-Agent", c.userAgent)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// For streaming, we don't want to retry as it could duplicate data
	return c.httpClient.Do(req)
}
