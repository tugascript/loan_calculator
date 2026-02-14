package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tugascript/loan_calculator/gateway/internal/config"
	v1 "github.com/tugascript/loan_calculator/gateway/internal/providers/proto/loan_calculator/v1"
)

type GatewayServer struct {
	logger *slog.Logger
	*http.Server
}

func (s *GatewayServer) Start(ctx context.Context) {
	logger := s.logger.With("method", "Start", "port", s.Server.Addr)
	logger.InfoContext(ctx, "Starting gateway server...")

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.ErrorContext(ctx, "Failed to start gateway server", "error", err)
		panic(err)
	}
}

func (s *GatewayServer) Close(ctx context.Context, done chan bool) {
	logger := s.logger.With("method", "Close", "port", s.Server.Addr)

	logger.InfoContext(ctx, "Closing gateway server...")
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.InfoContext(ctx, "Stopping gateway server...")

	if err := s.Shutdown(ctx); err != nil {
		logger.ErrorContext(ctx, "Failed to stop gateway server", "error", err)
		panic(err)
	}

	logger.InfoContext(ctx, "Gateway server stopped")
	done <- true
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	config *config.Config,
) *GatewayServer {
	logger.InfoContext(ctx, "Building gateway server...")

	gwmx := runtime.NewServeMux()

	transportCreds := insecure.NewCredentials()
	if config.ENV() == "production" {
		logger.InfoContext(ctx, "Production environment detected, configuring TLS for calculation service connection")
		tlsCfg := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		if name := config.CalculationTLSServerName(); name != "" {
			tlsCfg.ServerName = name
		}
		if path := config.CalculationTLSCertPath(); path != "" {
			pool, err := loadCACertPool(path)
			if err != nil {
				logger.ErrorContext(ctx, "Failed to load CA certificate", "path", path, "error", err)
				panic(err)
			}
			tlsCfg.RootCAs = pool
			logger.InfoContext(ctx, "Using custom CA for calculation service", "path", path)
		}
		transportCreds = credentials.NewTLS(tlsCfg)
		logger.InfoContext(ctx, "Using TLS for calculation service connection")
	}

	conn, err := grpc.NewClient(
		config.CalculationDomain(),
		grpc.WithTransportCredentials(transportCreds),
	)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to create gRPC connection", "error", err)
		panic(err)
	}
	if err := v1.RegisterCalculationServiceHandler(ctx, gwmx, conn); err != nil {
		logger.ErrorContext(ctx, "Failed to register gRPC gateway handler", "error", err)
		panic(err)
	}

	server := &GatewayServer{
		logger: logger.With("layer", "server"),
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port()),
			Handler: gwmx,
		},
	}

	logger.InfoContext(ctx, "Gateway server built successfully", "port", config.Port())
	return server
}

func loadCACertPool(path string) (*x509.CertPool, error) {
	pem, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read CA file: %w", err)
	}
	pool, err := x509.SystemCertPool()
	if err != nil {
		pool = x509.NewCertPool()
	}
	if !pool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("no valid CA certificates found in %s", path)
	}
	return pool, nil
}
