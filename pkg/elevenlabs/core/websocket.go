package core

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketError represents a WebSocket-specific error
type WebSocketError struct {
	Message string
}

func (e *WebSocketError) Error() string {
	return e.Message
}

// WebSocketClient handles WebSocket connections to the ElevenLabs API
type WebSocketClient struct {
	conn      *websocket.Conn
	apiKey    string
	baseURL   string
	connected bool
	mu        sync.RWMutex
}

// NewWebSocketClient creates a new WebSocket client
func NewWebSocketClient(baseURL, apiKey string) *WebSocketClient {
	return &WebSocketClient{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// Connect establishes a WebSocket connection to the specified endpoint
func (w *WebSocketClient) Connect(ctx context.Context, endpoint string, headers map[string]string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.connected {
		return nil // Already connected
	}

	// Prepare headers
	header := http.Header{}
	if w.apiKey != "" {
		header.Set(APIKeyHeader, w.apiKey)
	}

	// Add custom headers
	for key, value := range headers {
		header.Set(key, value)
	}

	// Create WebSocket connection
	dialer := websocket.Dialer{
		HandshakeTimeout: DefaultRetryConfig().MaxDelay,
	}

	url := w.baseURL + "/" + endpoint
	conn, _, err := dialer.DialContext(ctx, url, header)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	w.conn = conn
	w.connected = true

	return nil
}

// Send sends data through the WebSocket connection
func (w *WebSocketClient) Send(data interface{}) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if !w.connected || w.conn == nil {
		return &WebSocketError{Message: "WebSocket not connected"}
	}

	return w.conn.WriteJSON(data)
}

// SendText sends text data through the WebSocket connection
func (w *WebSocketClient) SendText(data string) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if !w.connected || w.conn == nil {
		return &WebSocketError{Message: "WebSocket not connected"}
	}

	return w.conn.WriteMessage(websocket.TextMessage, []byte(data))
}

// SendBinary sends binary data through the WebSocket connection
func (w *WebSocketClient) SendBinary(data []byte) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if !w.connected || w.conn == nil {
		return &WebSocketError{Message: "WebSocket not connected"}
	}

	return w.conn.WriteMessage(websocket.BinaryMessage, data)
}

// Receive receives data from the WebSocket connection
func (w *WebSocketClient) Receive() ([]byte, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if !w.connected || w.conn == nil {
		return nil, &WebSocketError{Message: "WebSocket not connected"}
	}

	_, message, err := w.conn.ReadMessage()
	return message, err
}

// ReceiveJSON receives JSON data from the WebSocket connection
func (w *WebSocketClient) ReceiveJSON(v interface{}) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if !w.connected || w.conn == nil {
		return &WebSocketError{Message: "WebSocket not connected"}
	}

	return w.conn.ReadJSON(v)
}

// Close closes the WebSocket connection
func (w *WebSocketClient) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.connected || w.conn == nil {
		return nil
	}

	err := w.conn.Close()
	w.connected = false
	w.conn = nil

	return err
}

// IsConnected returns whether the WebSocket is connected
func (w *WebSocketClient) IsConnected() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.connected
}

// SetCloseHandler sets a handler for connection close events
func (w *WebSocketClient) SetCloseHandler(handler func(code int, text string) error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.conn != nil {
		w.conn.SetCloseHandler(handler)
	}
}

// SetPingHandler sets a handler for ping messages
func (w *WebSocketClient) SetPingHandler(handler func(appData string) error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.conn != nil {
		w.conn.SetPingHandler(handler)
	}
}

// SetPongHandler sets a handler for pong messages
func (w *WebSocketClient) SetPongHandler(handler func(appData string) error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.conn != nil {
		w.conn.SetPongHandler(handler)
	}
}
