package services

import (
	"log/slog"

	"github.com/tugascript/loan_calculator/calculation_api/internal/providers/database"
)

type Services struct {
	*LoanCalculationRequestsService
}

func buildLogger(logger *slog.Logger, requestId, method string) *slog.Logger {
	return logger.With("requestId", requestId, "method", method)
}

func New(log *slog.Logger, db *database.Database) *Services {
	logger := log.With("layer", "services")
	return &Services{
		LoanCalculationRequestsService: newLoanCalculationRequestsService(logger, db),
	}
}
