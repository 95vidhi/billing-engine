package billing

import (
	"os"
	"testing"
)

func setupLoanStore() (*LoanStore, string) {
	store := &LoanStore{Loans: make(map[string]*Loan)}
	tempFile := "test_loans.json"
	_ = os.Remove(tempFile)
	return store, tempFile
}

func TestNewLoan(t *testing.T) {
	store, tempFile := setupLoanStore()
	defer os.Remove(tempFile)

	loanID, err := store.NewLoan(5000000, 10, 50, tempFile)
	if err != nil {
		t.Fatalf("Failed to create a new loan: %v", err)
	}

	loan, exists := store.Loans[loanID]
	if !exists {
		t.Fatalf("Loan was not stored correctly")
	}

	expectedOutstanding := 5500000.0
	if loan.Outstanding != expectedOutstanding {
		t.Errorf("Expected outstanding amount %v, got %v", expectedOutstanding, loan.Outstanding)
	}

	expectedWeeklyInstallment := 110000.0
	if loan.WeeklyInstallment != expectedWeeklyInstallment {
		t.Errorf("Expected weekly installment %v, got %v", expectedWeeklyInstallment, loan.WeeklyInstallment)
	}
}

func TestMakePayment(t *testing.T) {
	store, tempFile := setupLoanStore()
	defer os.Remove(tempFile)

	loanID, _ := store.NewLoan(5000000, 10, 50, tempFile)

	err := store.MakePayment(loanID, 1, tempFile)
	if err != nil {
		t.Errorf("Unexpected error while making payment: %v", err)
	}

	loan := store.Loans[loanID]
	expectedOutstanding := 5500000 - 110000
	if loan.Outstanding != float64(expectedOutstanding) {
		t.Errorf("Outstanding balance mismatch, got %v", loan.Outstanding)
	}

	err = store.MakePayment("invalid-id", 2, tempFile)
	if err == nil {
		t.Errorf("Expected error for invalid loan ID")
	}

	for i := 2; i <= 50; i++ {
		store.MakePayment(loanID, i, tempFile)
	}
	err = store.MakePayment(loanID, 51, tempFile)
	if err == nil {
		t.Errorf("Expected error for already paid loan")
	}
}

func TestGetOutstanding(t *testing.T) {
	store, tempFile := setupLoanStore()
	defer os.Remove(tempFile)

	loanID, _ := store.NewLoan(5000000, 10, 50, tempFile)

	outstanding, _ := store.GetOutstanding(loanID)
	expectedOutstanding := 5500000.0
	if outstanding != expectedOutstanding {
		t.Errorf("Expected outstanding %v, got %v", expectedOutstanding, outstanding)
	}

	store.MakePayment(loanID, 1, tempFile)
	outstanding, _ = store.GetOutstanding(loanID)
	expectedOutstanding = 5390000.0
	if outstanding != expectedOutstanding {
		t.Errorf("Expected outstanding %v, got %v", expectedOutstanding, outstanding)
	}

	for i := 2; i <= 50; i++ {
		store.MakePayment(loanID, i, tempFile)
	}
	outstanding, _ = store.GetOutstanding(loanID)
	if outstanding != 0 {
		t.Errorf("Expected outstanding 0, got %v", outstanding)
	}

	_, err := store.GetOutstanding("invalid-id")
	if err == nil {
		t.Errorf("Expected error for invalid loan ID")
	}
}

func TestIsDelinquent(t *testing.T) {
	store, tempFile := setupLoanStore()
	defer os.Remove(tempFile)

	loanID, _ := store.NewLoan(500000, 10, 10, tempFile)
	_ = store.Loans[loanID]

	isDelinquent, err := store.IsDelinquent(loanID, 1)
	if err != nil || isDelinquent {
		t.Errorf("Expected non-delinquent status at week 1, got delinquent")
	}

	store.MakePayment(loanID, 1, tempFile)

	isDelinquent, _ = store.IsDelinquent(loanID, 4)
	if !isDelinquent {
		t.Errorf("Expected delinquent status at week 3 after missing 2 consecutive payments")
	}

	store.MakePayment(loanID, 4, tempFile)
	isDelinquent, _ = store.IsDelinquent(loanID, 5)
	if isDelinquent {
		t.Errorf("Expected non-delinquent status at week 5 after making a payment")
	}

	isDelinquent, _ = store.IsDelinquent(loanID, 7)
	if !isDelinquent {
		t.Errorf("Expected delinquent status at week 7 after missing 2 consecutive payments")
	}

	_, err = store.IsDelinquent("invalid-id", 1)
	if err == nil {
		t.Errorf("Expected error for invalid loan ID, but got none")
	}
}
