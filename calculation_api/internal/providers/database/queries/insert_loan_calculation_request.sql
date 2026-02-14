-- name: InsertLoanCalculationRequest :one
INSERT INTO loan_calculation_requests (
    id,
    tenant_id, 
    loan_amount, 
    interest_rate, 
    number_of_payments, 
    monthly_repayment_amount
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
