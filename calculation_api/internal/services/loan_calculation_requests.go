package services

import (
	"context"
	"log/slog"
	"math/big"
	"strconv"

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

	logger.InfoContext(ctx, "Calculating monthly repayment amount",
		"loanAmount", opts.LoanAmount,
		"interestRate", opts.InterestRate,
		"numberOfPayments", opts.NumberOfPayments,
	)
	mraFloat64 := calculateMonthlyRepaymentAmount(
		opts.LoanAmount,
		opts.InterestRate,
		opts.NumberOfPayments,
	)
	logger.InfoContext(ctx, "Monthly repayment amount", "mraFloat64", mraFloat64)

	loanAmountStr := strconv.FormatFloat(opts.LoanAmount, 'f', 2, 64)
	interestRateStr := strconv.FormatFloat(opts.InterestRate, 'f', 4, 64)
	mraStr := strconv.FormatFloat(mraFloat64, 'f', 2, 64)

	loanAmount, serviceErr := scanNumeric(ctx, logger, "loanAmount", loanAmountStr)
	if serviceErr != nil {
		return nil, serviceErr
	}
	interestRate, serviceErr := scanNumeric(ctx, logger, "interestRate", interestRateStr)
	if serviceErr != nil {
		return nil, serviceErr
	}
	monthlyRepaymentAmount, serviceErr := scanNumeric(ctx, logger, "mraFloat64", mraStr)
	if serviceErr != nil {
		return nil, serviceErr
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
		logger.ErrorContext(ctx, "Failed to insert loan calculation request", "error", err)
		return nil, api_errors.FromDBError(err)
	}

	return dtos.MapLoanCalculationRequestToDTO(&loanCalculationRequest)
}

func scanNumeric(
	ctx context.Context,
	logger *slog.Logger,
	label string,
	value string,
) (pgtype.Numeric, *api_errors.ServiceError) {
	var n pgtype.Numeric
	if err := n.Scan(value); err != nil {
		logger.ErrorContext(ctx, "Failed to scan "+label, "error", err, "value", value)
		return pgtype.Numeric{}, api_errors.NewInternalServerError("Failed to scan " + label)
	}
	return n, nil
}

// Copied CODE, not mine.
// Fixed to use big.Rat instead of int64 to avoid overflow.
func calculateMonthlyRepaymentAmount(loanAmount, interestRate float64, numPayments int32) float64 {
	if numPayments <= 0 {
		return 0
	}

	// Convert inputs to big.Rat
	P := new(big.Rat).SetFloat64(loanAmount)

	// Monthly rate = annual / 12
	r := new(big.Rat).SetFloat64(interestRate)
	r.Quo(r, big.NewRat(12, 1))

	if r.Sign() == 0 {
		// No interest case
		result := new(big.Rat).Quo(P, big.NewRat(int64(numPayments), 1))
		f, _ := result.Float64()
		return f
	}

	n := int64(numPayments)

	// one = 1
	one := big.NewRat(1, 1)

	// (1 + r)
	onePlusR := new(big.Rat).Add(one, r)

	// power = (1 + r)^n
	power := powRat(onePlusR, n)

	// numerator = P * r * power
	num := new(big.Rat).Mul(P, r)
	num.Mul(num, power)

	// denominator = power - 1
	den := new(big.Rat).Sub(power, one)

	if den.Sign() == 0 {
		return 0
	}

	result := new(big.Rat).Quo(num, den)

	f, _ := result.Float64()
	return f
}

func powRat(base *big.Rat, exp int64) *big.Rat {
	result := big.NewRat(1, 1)

	for i := int64(0); i < exp; i++ {
		result.Mul(result, base)
	}

	return result
}
