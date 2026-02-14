package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"github.com/tugascript/loan_calculator/calculation_api/internal/config"
	"github.com/tugascript/loan_calculator/calculation_api/internal/controllers"
	"github.com/tugascript/loan_calculator/calculation_api/internal/providers/database"
	v1 "github.com/tugascript/loan_calculator/calculation_api/internal/providers/proto/loan_calculator/v1"
	"github.com/tugascript/loan_calculator/calculation_api/internal/services"
)

type GRPCServer struct {
	db     *database.Database
	logger *slog.Logger
	*grpc.Server
}

// TODO: looks like CTRL+C is not delayed, fix me at the end
func (s *GRPCServer) Close(ctx context.Context, done chan bool) {
	stpCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-stpCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.InfoContext(shutdownCtx, "Closing gRPC server...")
	grpcDone := make(chan any)
	go func() {
		s.Server.GracefulStop()
		close(grpcDone)
	}()

	select {
	case <-grpcDone:
		s.logger.InfoContext(shutdownCtx, "gRPC server closed")
	case <-shutdownCtx.Done():
		s.logger.WarnContext(shutdownCtx, "gRPC graceful stop timed out, forcing stop")
		s.Server.Stop()
		<-grpcDone
	}

	s.logger.InfoContext(shutdownCtx, "Closing database...")
	s.db.Close()
	s.logger.InfoContext(shutdownCtx, "database closed")

	s.logger.InfoContext(shutdownCtx, "All resources closed")
	done <- true
}

func (s *GRPCServer) Start(ctx context.Context, port int64) {
	s.logger.InfoContext(ctx, "Starting gRPC server...", "port", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to listen", "error", err)
		panic(err)
	}

	if err := s.Server.Serve(lis); err != nil {
		s.logger.ErrorContext(ctx, "Failed to serve", "error", err)
		panic(err)
	}
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	cfg *config.Config,
) *GRPCServer {
	logger.InfoContext(ctx, "Building database connection pool...")
	pgCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL())
	if err != nil {
		logger.ErrorContext(ctx, "Failed to parse database URL", "error", err)
		panic(err)
	}

	dbConn, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to create database connection pool", "error", err)
		panic(err)
	}

	db := database.NewDatabase(dbConn)
	logger.InfoContext(ctx, "Database connection pool created")

	logger.InfoContext(ctx, "Building services...")
	services := services.New(logger, db)
	logger.InfoContext(ctx, "Services built")

	logger.InfoContext(ctx, "Building controllers...")
	controllers := controllers.New(logger, services)
	logger.InfoContext(ctx, "Controllers built")

	logger.InfoContext(ctx, "Building gRPC server...")
	grpcServer := &GRPCServer{
		Server: grpc.NewServer(),
		db:     db,
		logger: logger.With("layer", "server"),
	}
	v1.RegisterCalculationServiceServer(grpcServer.Server, controllers.LoanCalculationRequestsController)
	logger.InfoContext(ctx, "gRPC server built")

	return grpcServer
}
