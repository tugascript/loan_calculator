package controllers

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/tugascript/loan_calculator/calculation_api/internal/api_errors"
	v1 "github.com/tugascript/loan_calculator/calculation_api/internal/providers/proto/loan_calculator/v1"
	"github.com/tugascript/loan_calculator/calculation_api/internal/services"
)

type LoanCalculationRequestsController struct {
	services  *services.Services
	logger    *slog.Logger
	validator *validator.Validate
	v1.UnimplementedCalculationServiceServer
}

func newLoanCalculationRequestController(logger *slog.Logger, validate *validator.Validate, services *services.Services) *LoanCalculationRequestsController {
	return &LoanCalculationRequestsController{
		services:  services,
		logger:    logger,
		validator: validate,
	}
}

type calculateMonthlyRepaymentRequest struct {
	LoanAmount       float64 `validate:"required,min=0.01"`
	InterestRate     float64 `validate:"required,min=0,max=1"`
	NumberOfPayments int32   `validate:"required,min=1"`
}

func (lc *LoanCalculationRequestsController) CalculateMonthlyRepayment(
	ctx context.Context,
	req *v1.CalculateMonthlyRepaymentRequest,
) (*v1.CalculateMonthlyRepaymentResponse, error) {
	requestID := generateRequestID()
	logger := lc.logger.With("method", "CalculateMonthlyRepayment", "requestID", requestID)
	logger.InfoContext(ctx, "Calculating monthly repayment", "req", req)

	request := calculateMonthlyRepaymentRequest{
		LoanAmount:       req.LoanAmount,
		InterestRate:     req.InterestRate,
		NumberOfPayments: req.NumberOfPayments,
	}
	if err := lc.validator.StructCtx(ctx, &request); err != nil {
		logger.ErrorContext(ctx, "Invalid request", "error", err)
		return nil, api_errors.InvalidArgumentResponseFromErr(err)
	}

	loanCalculationRequest, serviceErr := lc.services.LoanCalculationRequestsService.CreateLoanCalculationRequest(
		ctx,
		services.CreateLoanCalculationRequestOptions{
			RequestID:        requestID,
			TenantID:         uuid.New(),
			LoanAmount:       request.LoanAmount,
			InterestRate:     request.InterestRate,
			NumberOfPayments: request.NumberOfPayments,
		},
	)
	if serviceErr != nil {
		return nil, api_errors.GRPCStatus(serviceErr)
	}

	logger.InfoContext(ctx, "Loan calculation request created successfully")
	return loanCalculationRequest, nil
}
