package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func TestProcessQueue(t *testing.T) {
	logger, _ = zap.NewProduction()
	// Mock Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	// Mock webhook server
	webhookServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result ScanResult
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &result)

		if result.Image != "registry.suse.com/bci/bci-busybox:latest" {
			t.Errorf("Unexpected image name: %s", result.Image)
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer webhookServer.Close()

	// Set environment variables
	os.Setenv("WEBHOOK_URL", webhookServer.URL)
	os.Setenv("REDIS_HOST", mr.Host())
	os.Setenv("REDIS_PORT", mr.Port())

	// Initialize Redis client and push a mock entry
	rdb = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	_, err = rdb.RPush(ctx, "topush", "registry.suse.com/bci/bci-busybox:latest|./test_reports/registry.suse.com_bci_bci-busybox_latest.json").Result()
	if err != nil {
		t.Fatal(err)
	}

	// Call the function to process the queue
	processQueue(webhookServer.URL)

	// Allow some time for the HTTP request to complete
	time.Sleep(time.Second)
}
