package voices

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/LinkedDestiny/elevenlabs-golang/pkg/elevenlabs/core"
)

// Client handles voice operations
type Client struct {
	httpClient *core.HTTPClient
}

// NewClient creates a new voices client
func NewClient(httpClient *core.HTTPClient) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

// GetAll retrieves all available voices
func (c *Client) GetAll(ctx context.Context, opts GetAllOptions) (*VoicesResponse, error) {
	path := "v1/voices"

	// Add query parameters
	params := make([]string, 0)
	if opts.ShowLegacy != nil {
		params = append(params, fmt.Sprintf("show_legacy=%t", *opts.ShowLegacy))
	}
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	// Make the request
	resp, err := c.httpClient.Request(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	// Parse the response
	var result VoicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Get retrieves a specific voice by ID
func (c *Client) Get(ctx context.Context, voiceID string, opts GetOptions) (*Voice, error) {
	path := fmt.Sprintf("v1/voices/%s", voiceID)

	// Add query parameters
	params := make([]string, 0)
	if opts.WithSettings != nil {
		params = append(params, fmt.Sprintf("with_settings=%t", *opts.WithSettings))
	}
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	// Make the request
	resp, err := c.httpClient.Request(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	// Parse the response
	var result Voice
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Delete removes a voice
func (c *Client) Delete(ctx context.Context, voiceID string) error {
	path := fmt.Sprintf("v1/voices/%s", voiceID)

	// Make the request
	resp, err := c.httpClient.Request(ctx, "DELETE", path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}

	return nil
}

// GetSettings retrieves voice settings
func (c *Client) GetSettings(ctx context.Context, voiceID string) (*VoiceSettings, error) {
	path := fmt.Sprintf("v1/voices/%s/settings", voiceID)

	// Make the request
	resp, err := c.httpClient.Request(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	// Parse the response
	var result VoiceSettings
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// EditSettings updates voice settings
func (c *Client) EditSettings(ctx context.Context, voiceID string, settings VoiceSettings) (*VoiceSettings, error) {
	path := fmt.Sprintf("v1/voices/%s/settings/edit", voiceID)

	// Prepare request body
	requestBody, err := json.Marshal(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Set headers
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Make the request
	resp, err := c.httpClient.Request(ctx, "POST", path, strings.NewReader(string(requestBody)), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	// Parse the response
	var result VoiceSettings
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// parseAPIError parses an HTTP error response
func parseAPIError(resp *http.Response) error {
	return fmt.Errorf("API error: %s", resp.Status)
}
