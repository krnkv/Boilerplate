package tracing_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krnkv/Boilerplate/internal/config"
	"github.com/krnkv/Boilerplate/internal/logger"
	"github.com/krnkv/Boilerplate/internal/observability/tracing"
)

// TestNewTracerService_Success ensures tracer initialization works with valid config.
func TestNewTracerService_Success(t *testing.T) {
	ctx := context.Background()
	log := logger.NewZerologLogger("info", io.Discard)

	cfg := &config.Tracing{
		ServiceName:  "test-service",
		CollectorURL: "localhost:4318", // dummy collector
	}

	svc, err := tracing.NewTracerService(ctx, &tracing.Opts{
		Config: cfg,
		Logger: log,
	})
	require.NoError(t, err)
	require.NotNil(t, svc)

	assert.Equal(t, cfg.ServiceName, svc.Config.ServiceName)
	assert.Equal(t, cfg.CollectorURL, svc.Config.CollectorURL)

	// Ensure shutdown does not panic or error out
	err = svc.Shutdown(ctx)
	assert.NoError(t, err)
}
