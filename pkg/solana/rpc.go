package solana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// VersionResponse represents the response from the getVersion RPC method
type VersionResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		FeatureSet int64  `json:"feature-set"`
		SolanaCore string `json:"solana-core"`
	} `json:"result"`
	ID int `json:"id"`
}

// Version represents the Solana version information
type Version struct {
	SolanaCore    string `json:"solana_core"`
	FeatureSet    int64  `json:"feature_set"`
	FetchedAt     string `json:"fetched_at"`
	FetchedAtUnix int64  `json:"fetched_at_unix"`
}

// GetVersion fetches the Solana version from the RPC endpoint
func GetVersion(rpcUrl string) (*Version, error) {
	// Create the request payload
	payload := []byte(`{"jsonrpc":"2.0","id":1,"method":"getVersion"}`)

	// Create the HTTP request
	req, err := http.NewRequest("POST", rpcUrl, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response
	var versionResp VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Get current time
	now := time.Now()

	// Return the version information
	return &Version{
		SolanaCore:    versionResp.Result.SolanaCore,
		FeatureSet:    versionResp.Result.FeatureSet,
		FetchedAt:     now.Format(time.RFC3339),
		FetchedAtUnix: now.Unix(),
	}, nil
}
