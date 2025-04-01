package billing

type BillingEngine struct {
	repo LoanRepository
}

func NewBillingEngine(repo LoanRepository) *BillingEngine {
	return &BillingEngine{repo: repo}
}

func (be *BillingEngine) NewLoan(principal float64, rate float64, weeks int) (string, error) {
	loan, err := NewLoan(principal, rate, weeks)
	if err != nil {
		return "", err
	}
	return be.repo.Save(loan)
}

func (be *BillingEngine) GetOutstanding(loanID string) (float64, error) {
	loan, err := be.repo.FindByID(loanID)
	if err != nil {
		return 0, err
	}
	return loan.OutstandingAmount(), nil
}

func (be *BillingEngine) GeneratePaymentSchedule(loanID string) ([]PaymentSchedule, error) {
	loan, err := be.repo.FindByID(loanID)
	if err != nil {
		return nil, err
	}
	return loan.GenerateSchedule(), nil
}

func (be *BillingEngine) MakePayment(loanID string, week int, amount float64) error {
	loan, err := be.repo.FindByID(loanID)
	if err != nil {
		return err
	}

	err = loan.ApplyPayment(week, amount)
	if err != nil {
		return err
	}
	_, err = be.repo.Save(loan)

	return err
}

func (be *BillingEngine) IsDelinquent(loanID string, currentWeek int) (bool, error) {
	loan, err := be.repo.FindByID(loanID)
	if err != nil {
		return false, err
	}
	return loan.IsDelinquent(currentWeek)
}
