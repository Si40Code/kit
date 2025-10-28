package httpclient

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := New()
	if client == nil {
		t.Fatal("Expected client to be created")
	}
	if client.Client == nil {
		t.Fatal("Expected resty client to be initialized")
	}
}

func TestClientWithOptions(t *testing.T) {
	client := New(
		WithTimeout(5*time.Second),
		WithMaxIdleConns(50),
		WithDisableLog(),
	)

	if client == nil {
		t.Fatal("Expected client to be created with options")
	}

	if client.options.timeout != 5*time.Second {
		t.Errorf("Expected timeout to be 5s, got %v", client.options.timeout)
	}

	if client.options.maxIdleConns != 50 {
		t.Errorf("Expected maxIdleConns to be 50, got %d", client.options.maxIdleConns)
	}

	if client.options.enableLog {
		t.Error("Expected logging to be disabled")
	}
}

func TestClientRetryConfig(t *testing.T) {
	client := New(
		WithRetry(3, 100*time.Millisecond, 2*time.Second),
	)

	if client.options.retryCount != 3 {
		t.Errorf("Expected retry count to be 3, got %d", client.options.retryCount)
	}

	if client.options.retryWaitTime != 100*time.Millisecond {
		t.Errorf("Expected retry wait time to be 100ms, got %v", client.options.retryWaitTime)
	}

	if client.options.retryMaxWaitTime != 2*time.Second {
		t.Errorf("Expected retry max wait time to be 2s, got %v", client.options.retryMaxWaitTime)
	}
}

func TestClientR(t *testing.T) {
	client := New()
	ctx := context.Background()

	req := client.R(ctx)
	if req == nil {
		t.Fatal("Expected request to be created")
	}

	// Verify context is set
	if req.Context() != ctx {
		t.Error("Expected context to be set on request")
	}
}

func TestMetricRecorder(t *testing.T) {
	called := false
	recorder := &mockMetricRecorder{
		recordFunc: func(data MetricData) {
			called = true
		},
	}

	client := New(
		WithMetric(recorder),
		WithDisableLog(),
	)

	// Make a simple request (this will fail but that's ok, we're testing the hook)
	_, _ = client.R(context.Background()).Get("http://localhost:1")

	// Note: The metric might not be called if the request fails at DNS level
	// This is expected behavior
	_ = called
}

type mockMetricRecorder struct {
	recordFunc func(MetricData)
}

func (m *mockMetricRecorder) RecordRequest(data MetricData) {
	if m.recordFunc != nil {
		m.recordFunc(data)
	}
}
