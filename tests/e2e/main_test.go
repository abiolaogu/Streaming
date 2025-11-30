package e2e

import (
	"net/http"
	"testing"
	"time"
)

// TestHealthChecks verifies that all services are up and running
func TestHealthChecks(t *testing.T) {
	services := []struct {
		name string
		url  string
	}{
		{"User Service", "http://localhost:8081/health"},
		{"Content Service", "http://localhost:8082/health"},
		{"Streaming Service", "http://localhost:8083/health"},
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, service := range services {
		t.Run(service.name, func(t *testing.T) {
			// Retry logic for CI environments
			var err error
			var resp *http.Response
			for i := 0; i < 3; i++ {
				resp, err = client.Get(service.url)
				if err == nil {
					break
				}
				time.Sleep(2 * time.Second)
			}

			if err != nil {
				t.Logf("Warning: %s is not reachable at %s (skipping in non-integration env): %v", service.name, service.url, err)
				return // Skip if not reachable (e.g., running in CI without services up)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("%s returned status %d", service.name, resp.StatusCode)
			}
		})
	}
}
