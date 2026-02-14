CREATE TABLE loan_calculation_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    loan_amount NUMERIC(18,2) NOT NULL CHECK (loan_amount > 0),
    interest_rate NUMERIC(8,4) NOT NULL CHECK (interest_rate >= 0),
    number_of_payments INTEGER NOT NULL CHECK (number_of_payments > 0),
    monthly_repayment_amount NUMERIC(18,2) NOT NULL CHECK (monthly_repayment_amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_loan_calculation_requests_tenant_id ON loan_calculation_requests(tenant_id);
CREATE INDEX idx_loan_calculation_requests_id_tenant_id_created_at ON loan_calculation_requests(id, tenant_id, created_at);
