package controllers

import (
	"log/slog"
	"math/big"

	"github.com/google/uuid"
	"github.com/tugascript/loan_calculator/calculation_api/internal/services"
)

type Controllers struct {
	*LoanCalculationRequestsController
}

func New(log *slog.Logger, services *services.Services) *Controllers {
	logger := log.With("layer", "controllers")
	return &Controllers{
		LoanCalculationRequestsController: newLoanCalculationRequestController(services, logger),
	}
}

func base62Encode(bytes []byte) string {
	return new(big.Int).SetBytes(bytes).Text(62)
}

func generateRequestID() string {
	id := uuid.New()
	return base62Encode(id[:])
}
