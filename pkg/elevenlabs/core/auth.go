package core

import "net/http"

// APIKeyHeader is the header name for the ElevenLabs API key
const APIKeyHeader = "xi-api-key"

// AddAuthHeaders adds authentication headers to an HTTP request
func AddAuthHeaders(req *http.Request, apiKey string) {
	if apiKey != "" {
		req.Header.Set(APIKeyHeader, apiKey)
	}
}
