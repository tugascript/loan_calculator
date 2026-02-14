package controllers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tugascript/loan_calculator/calculation_api/internal/api_errors"
	v1 "github.com/tugascript/loan_calculator/calculation_api/internal/providers/proto/loan_calculator/v1"
	"github.com/tugascript/loan_calculator/calculation_api/internal/services"
)

type LoanCalculationRequestsController struct {
	services *services.Services
	logger   *slog.Logger
	v1.UnimplementedCalculationServiceServer
}

func newLoanCalculationRequestController(services *services.Services, logger *slog.Logger) *LoanCalculationRequestsController {
	return &LoanCalculationRequestsController{
		services: services,
		logger:   logger,
	}
}

func (lc *LoanCalculationRequestsController) CalculateMonthlyRepayment(
	ctx context.Context,
	req *v1.CalculateMonthlyRepaymentRequest,
) (*v1.CalculateMonthlyRepaymentResponse, error) {
	requestID := generateRequestID()
	logger := lc.logger.With("method", "CalculateMonthlyRepayment", "requestID", requestID)
	logger.InfoContext(ctx, "Calculating monthly repayment", "req", req)

	loanCalculationRequest, serviceErr := lc.services.LoanCalculationRequestsService.CreateLoanCalculationRequest(
		ctx,
		services.CreateLoanCalculationRequestOptions{
			RequestID:        requestID,
			TenantID:         uuid.New(),
			LoanAmount:       req.LoanAmount,
			InterestRate:     req.InterestRate,
			NumberOfPayments: req.NumberOfPayments,
		},
	)
	if serviceErr != nil {
		return nil, api_errors.GRPCStatus(serviceErr)
	}

	logger.InfoContext(ctx, "Loan calculation request created successfully")
	return loanCalculationRequest, nil
}
