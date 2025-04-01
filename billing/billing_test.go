package billing

import (
	"errors"
	"testing"
)

type mockLoanRepository struct {
	loans map[string]*Loan
}

func (m *mockLoanRepository) Save(loan *Loan) (string, error) {
	m.loans[loan.ID] = loan
	return loan.ID, nil
}

func (m *mockLoanRepository) FindByID(loanID string) (*Loan, error) {
	loan, exists := m.loans[loanID]
	if !exists {
		return nil, errors.New("loan not found")
	}
	return loan, nil
}

func setupBillingEngine() (*BillingEngine, *mockLoanRepository) {
	repo := &mockLoanRepository{loans: make(map[string]*Loan)}
	return NewBillingEngine(repo), repo
}

func TestMakePayment(t *testing.T) {
	be, _ := setupBillingEngine()
	loanID, _ := be.NewLoan(1000, 5, 10)

	tests := []struct {
		name       string
		week       int
		amount     float64
		expectsErr bool
	}{
		{"Already paid week", 1, 100, true},
		{"Paying after 50 weeks", 51, 100, true},
		{"Overpaying", 2, 2000, true},
		{"Multiple missed payments", 3, 105, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := be.MakePayment(loanID, tt.week, tt.amount)
			if (err != nil) != tt.expectsErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestIsDelinquent(t *testing.T) {
	be, _ := setupBillingEngine()
	loanID, _ := be.NewLoan(1000, 5, 10)
	be.MakePayment(loanID, 1, 100)

	tests := []struct {
		name        string
		currentWeek int
		expected    bool
	}{
		{"Loan becomes delinquent", 4, true},
		{"Loan recovers from delinquency", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isDelinquent, err := be.IsDelinquent(loanID, tt.currentWeek)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if isDelinquent != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, isDelinquent)
			}
		})
	}
}

func TestGetOutstanding(t *testing.T) {
	be, _ := setupBillingEngine()
	loanID, _ := be.NewLoan(5000000, 10, 50)
	be.MakePayment(loanID, 1, 100)

	outstanding, err := be.GetOutstanding(loanID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if outstanding != float64(5500000) {
		t.Errorf("expected outstanding to be 900, got %f", outstanding)
	}
}

func TestGeneratePaymentSchedule(t *testing.T) {
	// Setup mock or real repository and BillingEngine
	repo := &mockLoanRepository{loans: make(map[string]*Loan)}
	be := &BillingEngine{repo: repo}

	// Create a loan
	loanID := "1234"
	loan := &Loan{
		ID:                loanID,
		Principal:         1000,
		Rate:              5,
		Weeks:             10,
		LoanAmount:        1050, // 1000 principal + 5% interest
		Outstanding:       1050,
		WeeklyInstallment: 105,
		Schedule:          make([]PaymentSchedule, 10), // Ensure the schedule is initialized
	}

	// Populate the schedule
	for i := 0; i < loan.Weeks; i++ {
		loan.Schedule[i] = PaymentSchedule{Week: i + 1, Status: "Pending"}
	}

	repo.loans[loanID] = loan // Save loan to the mock repository

	// Call GeneratePaymentSchedule
	schedule, err := be.GeneratePaymentSchedule(loanID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that the schedule has the expected number of weeks
	if len(schedule) != loan.Weeks {
		t.Errorf("expected %d weeks in the schedule, got %d", loan.Weeks, len(schedule))
	}

	// Verify that the first week's status is "Pending"
	if schedule[0].Status != "Pending" {
		t.Errorf("expected first week's status to be 'Pending', got '%s'", schedule[0].Status)
	}

	// Verify the last week's status is still "Pending"
	if schedule[len(schedule)-1].Status != "Pending" {
		t.Errorf("expected last week's status to be 'Pending', got '%s'", schedule[len(schedule)-1].Status)
	}

	// Verify the loan ID is correctly used
	if schedule[0].Week != 1 {
		t.Errorf("expected first week to be 1, got %d", schedule[0].Week)
	}
}
