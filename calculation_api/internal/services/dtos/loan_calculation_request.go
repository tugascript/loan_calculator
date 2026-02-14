package dtos

import (
	"github.com/tugascript/loan_calculator/calculation_api/internal/api_errors"
	"github.com/tugascript/loan_calculator/calculation_api/internal/providers/database"
	v1 "github.com/tugascript/loan_calculator/calculation_api/internal/providers/proto/loan_calculator/v1"
)

func MapLoanCalculationRequestToDTO(loanCalculationRequest *database.LoanCalculationRequest) (*v1.CalculateMonthlyRepaymentResponse, *api_errors.ServiceError) {
	monthlyRepaymentAmount, err := loanCalculationRequest.MonthlyRepaymentAmount.Float64Value()
	if err != nil {
		return nil, api_errors.NewInternalServerError("Failed to convert monthly repayment amount to float64")
	}
	if !monthlyRepaymentAmount.Valid {
		return nil, api_errors.NewInternalServerError("Monthly repayment amount is not valid")
	}

	return &v1.CalculateMonthlyRepaymentResponse{
		Id:                     loanCalculationRequest.ID.String(),
		MonthlyRepaymentAmount: monthlyRepaymentAmount.Float64,
	}, nil
}
