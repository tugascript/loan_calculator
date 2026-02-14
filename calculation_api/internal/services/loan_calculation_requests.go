package services

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tugascript/loan_calculator/calculation_api/internal/api_errors"
	"github.com/tugascript/loan_calculator/calculation_api/internal/providers/database"
	v1 "github.com/tugascript/loan_calculator/calculation_api/internal/providers/proto/loan_calculator/v1"
	"github.com/tugascript/loan_calculator/calculation_api/internal/services/dtos"
)

type LoanCalculationRequestsService struct {
	logger *slog.Logger
	db     *database.Database
}

func newLoanCalculationRequestsService(logger *slog.Logger, db *database.Database) *LoanCalculationRequestsService {
	return &LoanCalculationRequestsService{
		logger: logger.With("service", "loan_calculation_requests"),
		db:     db,
	}
}

type CreateLoanCalculationRequestOptions struct {
	RequestID        string
	TenantID         uuid.UUID
	LoanAmount       float64
	InterestRate     float64
	NumberOfPayments int32
}

func (lc *LoanCalculationRequestsService) CreateLoanCalculationRequest(
	ctx context.Context,
	opts CreateLoanCalculationRequestOptions,
) (*v1.CalculateMonthlyRepaymentResponse, *api_errors.ServiceError) {
	logger := buildLogger(lc.logger, opts.RequestID, "CreateLoanCalculationRequest")
	logger.InfoContext(ctx, "Creating loan calculation request", "opts", opts)

	mraFloat64 := calculateMonthlyRepaymentAmount(
		opts.LoanAmount,
		opts.InterestRate,
		opts.NumberOfPayments,
	)

	var loanAmount pgtype.Numeric
	if err := loanAmount.Scan(opts.LoanAmount); err != nil {
		return nil, api_errors.NewServiceError(api_errors.CodeInternalServerError, "Failed to scan loan amount")
	}

	var interestRate pgtype.Numeric
	if err := interestRate.Scan(opts.InterestRate); err != nil {
		return nil, api_errors.NewServiceError(api_errors.CodeInternalServerError, "Failed to scan interest rate")
	}

	var monthlyRepaymentAmount pgtype.Numeric
	if err := monthlyRepaymentAmount.Scan(mraFloat64); err != nil {
		return nil, api_errors.NewServiceError(api_errors.CodeInternalServerError, "Failed to scan monthly repayment amount")
	}

	loanCalculationRequest, err := lc.db.InsertLoanCalculationRequest(ctx, database.InsertLoanCalculationRequestParams{
		ID:                     uuid.New(),
		TenantID:               opts.TenantID,
		LoanAmount:             loanAmount,
		InterestRate:           interestRate,
		NumberOfPayments:       opts.NumberOfPayments,
		MonthlyRepaymentAmount: monthlyRepaymentAmount,
	})
	if err != nil {
		return nil, api_errors.NewServiceError(api_errors.CodeInternalServerError, "Failed to insert loan calculation request")
	}

	return dtos.MapLoanCalculationRequestToDTO(&loanCalculationRequest)
}

// Copied CODE, not mine.
func calculateMonthlyRepaymentAmount(loanAmount, interestRate float64, numPayments int32) float64 {
	// PMT = (P * r * (1 + r)^n) / ((1 + r)^n - 1)
	// where:
	// P = loanAmount
	// r = monthly interest rate
	// n = number of payments
	if numPayments <= 0 {
		return 0
	}
	if interestRate == 0 {
		return loanAmount / float64(numPayments)
	}

	// Scale values to work with integers (simulate fixed point)
	const scale = 1_000_000

	P := int64(loanAmount * scale)
	r := int64(interestRate * scale)
	n := int64(numPayments)

	// Use integers for calculation, simulating fixed-point math
	one := int64(scale)

	// Calculate (1 + r) as an integer (still scaled)
	onePlusR := one + r

	// powInt returns (onePlusR ^ n) / (scale ^ (n-1)) by multiplying with scale each time
	power := powInt(onePlusR, n, scale)

	// Numerator: P * r * power
	// (We must divide by scale twice, once for r, once for power)
	num := P * r * power
	num = num / scale / scale

	// Denominator: power - one (both scaled)
	den := power - one

	if den == 0 {
		return 0
	}

	result := float64(num) / float64(den)
	return result
}

// powInt performs integer exponentiation, keeping scale in mind.
// base and scale should have the same scale factor as the rest of the calculation,
// so powInt(onePlusR, n, scale) computes (onePlusR/scale)^n as a scaled integer.
func powInt(base, exp, scale int64) int64 {
	result := scale
	for range exp {
		result = (result * base) / scale
	}
	return result
}
