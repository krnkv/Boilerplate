package server_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/krnkv/Boilerplate/internal/config"
	"github.com/krnkv/Boilerplate/internal/logger"
	"github.com/krnkv/Boilerplate/internal/tests/mock"
	"github.com/krnkv/Boilerplate/internal/transports/grpc/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

// TestNewServer ensures a GRPCServer struct is created with correct dependencies.
func TestNewServer(t *testing.T) {
	log := logger.NewZerologLogger("info", io.Discard)
	cfg := &config.GRPCServer{}

	mockDB := new(mock.MockDatabase)
	mockCache := new(mock.MockRedisCache)

	// Create in-memory sqlite gorm.DB
	fakeDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockDB.On("DB").Return(fakeDB)
	mockDB.On("Close").Return(nil)

	srv := server.NewServer(&server.Opts{
		Config:   cfg,
		Logger:   log,
		Database: mockDB,
		Cache:    mockCache,
	})

	require.NotNil(t, srv)
	assert.Equal(t, cfg, srv.Config)
	assert.Equal(t, log, srv.Logger)
	assert.NotNil(t, srv.Server)
}

// TestServeListener verifies server can run on a bufconn listener.
func TestServeListener(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)

	mockDB := new(mock.MockDatabase)
	mockCache := new(mock.MockRedisCache)

	// Create in-memory sqlite gorm.DB
	fakeDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockDB.On("DB").Return(fakeDB)
	mockDB.On("Close").Return(nil)

	srv := server.NewServer(&server.Opts{
		Config:   &config.GRPCServer{},
		Logger:   log,
		Database: mockDB,
		Cache:    mockCache,
	})

	lis := bufconn.Listen(bufSize)

	go func() {
		_ = srv.ServeListener(lis)
	}()
	defer srv.Server.Stop()

	// Give server some time to start
	time.Sleep(100 * time.Millisecond)

	assert.Contains(t, buf.String(), "gRPC server started")
}

// TestServe verifies server can run on a real TCP listener (ephemeral port).
func TestServe(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologLogger("info", &buf)

	mockDB := new(mock.MockDatabase)
	mockCache := new(mock.MockRedisCache)

	// Create in-memory sqlite gorm.DB
	fakeDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	mockDB.On("DB").Return(fakeDB)
	mockDB.On("Close").Return(nil)

	srv := server.NewServer(&server.Opts{
		Config: &config.GRPCServer{
			URL: ":0", // Use :0 to let OS pick a free port
		},
		Logger:   log,
		Database: mockDB,
		Cache:    mockCache,
	})

	go func() {
		_ = srv.Serve()
	}()
	defer srv.Server.Stop()

	// Give server some time to start
	time.Sleep(100 * time.Millisecond)

	assert.Contains(t, buf.String(), "gRPC server started")
}
