package text_to_speech

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LinkedDestiny/elevenlabs-golang/pkg/elevenlabs/core"
)

// Client handles text-to-speech operations
type Client struct {
	httpClient *core.HTTPClient
}

// NewClient creates a new text-to-speech client
func NewClient(httpClient *core.HTTPClient) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

// Convert converts text to speech and returns audio bytes
func (c *Client) Convert(ctx context.Context, req ConvertRequest) ([]byte, error) {
	// Build the request path
	path := fmt.Sprintf("v1/text-to-speech/%s", req.VoiceID)

	// Prepare request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Set headers
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "audio/mpeg",
	}

	// Make the request
	resp, err := c.httpClient.Request(ctx, "POST", path, bytes.NewReader(requestBody), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	// Read the audio data
	audioData := make([]byte, 0)
	buffer := make([]byte, 8192)

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			audioData = append(audioData, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	return audioData, nil
}

// ConvertWithTimestamps converts text to speech with timing information
func (c *Client) ConvertWithTimestamps(ctx context.Context, req ConvertRequest) (*TimestampResponse, error) {
	// Build the request path
	path := fmt.Sprintf("v1/text-to-speech/%s/with-timestamps", req.VoiceID)

	// Prepare request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Set headers
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// Make the request
	resp, err := c.httpClient.Request(ctx, "POST", path, bytes.NewReader(requestBody), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	// Parse the response
	var result TimestampResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Stream converts text to speech and returns a channel of audio chunks
func (c *Client) Stream(ctx context.Context, req StreamRequest) (<-chan []byte, error) {
	// Build the request path
	path := fmt.Sprintf("v1/text-to-speech/%s/stream", req.VoiceID)

	// Prepare request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Set headers
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "audio/mpeg",
	}

	// Make the streaming request
	resp, err := c.httpClient.Stream(ctx, "POST", path, bytes.NewReader(requestBody), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, parseAPIError(resp)
	}

	// Return streaming channel
	return core.StreamResponse(resp, 8192), nil
}

// StreamWithTimestamps converts text to speech with timing information in streaming mode
func (c *Client) StreamWithTimestamps(ctx context.Context, req StreamRequest) (<-chan TimestampChunk, error) {
	// Build the request path
	path := fmt.Sprintf("v1/text-to-speech/%s/stream-with-timestamps", req.VoiceID)

	// Prepare request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Set headers
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// Make the streaming request
	resp, err := c.httpClient.Stream(ctx, "POST", path, bytes.NewReader(requestBody), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, parseAPIError(resp)
	}

	// Create channel for timestamp chunks
	ch := make(chan TimestampChunk)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		// Stream JSON lines
		lines := core.StreamLines(ctx, resp)
		for chunk := range lines {
			if chunk.Err != nil {
				return
			}

			// Parse each line as a TimestampChunk
			var timestampChunk TimestampChunk
			if err := json.Unmarshal(chunk.Data, &timestampChunk); err == nil {
				ch <- timestampChunk
			}
		}
	}()

	return ch, nil
}

// ConvertRealtime performs real-time text-to-speech conversion via WebSocket
func (c *Client) ConvertRealtime(ctx context.Context, req RealtimeRequest) (<-chan []byte, error) {
	// Create WebSocket client
	wsClient := core.NewWebSocketClient(c.httpClient.GetWebSocketURL(), c.httpClient.GetAPIKey())

	// Build the WebSocket endpoint path
	path := fmt.Sprintf("v1/text-to-speech/%s/stream-input", req.VoiceID)

	// Connect to WebSocket
	if err := wsClient.Connect(ctx, path, nil); err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Send initial configuration
	initMessage := map[string]interface{}{
		"model_id":       req.ModelID,
		"output_format":  req.OutputFormat,
		"voice_settings": req.VoiceSettings,
	}

	if err := wsClient.Send(initMessage); err != nil {
		wsClient.Close()
		return nil, fmt.Errorf("failed to send initial configuration: %w", err)
	}

	// Create audio output channel
	audioCh := make(chan []byte)

	// Start goroutines for sending text and receiving audio
	go func() {
		defer wsClient.Close()

		// Send text chunks
		for text := range req.TextStream {
			textMessage := map[string]interface{}{
				"text": text,
			}
			if err := wsClient.Send(textMessage); err != nil {
				return
			}
		}

		// Send end of stream
		endMessage := map[string]interface{}{
			"text": "",
		}
		wsClient.Send(endMessage)
	}()

	go func() {
		defer close(audioCh)

		// Receive audio chunks
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := wsClient.Receive()
				if err != nil {
					return
				}

				// Parse the message to extract audio data
				var message map[string]interface{}
				if err := json.Unmarshal(data, &message); err == nil {
					if audioData, ok := message["audio"].(string); ok {
						// Audio is base64 encoded, you might want to decode it
						audioCh <- []byte(audioData)
					}
				}
			}
		}
	}()

	return audioCh, nil
}

// parseAPIError parses an HTTP error response
func parseAPIError(resp *http.Response) error {
	// This is a simplified error parser - you might want to implement
	// more sophisticated error parsing based on the actual API responses
	return fmt.Errorf("API error: %s", resp.Status)
}
