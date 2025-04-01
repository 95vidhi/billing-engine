package billing

import (
	"os"
	"testing"
)

func TestLoanRepository_SaveAndFind(t *testing.T) {
	filePath := "test_loans.json"
	os.Remove(filePath)
	repo := NewLoanRepository(filePath)

	loan := &Loan{ID: "123", Principal: 1000}
	_, err := repo.Save(loan)
	if err != nil {
		t.Fatalf("failed to save loan: %v", err)
	}

	savedLoan, err := repo.FindByID("123")
	if err != nil {
		t.Fatalf("failed to find loan: %v", err)
	}
	if savedLoan.Principal != 1000 {
		t.Errorf("expected principal 1000, got %v", savedLoan.Principal)
	}

	os.Remove(filePath)
}

func TestLoanRepository_SaveMultipleLoans(t *testing.T) {
	repo := NewLoanRepository("test_loans.json")
	os.Remove("test_loans.json")

	loans := []*Loan{{ID: "1", Principal: 500}, {ID: "2", Principal: 1000}}
	for _, loan := range loans {
		repo.Save(loan)
	}

	savedLoan1, _ := repo.FindByID("1")
	savedLoan2, _ := repo.FindByID("2")

	if savedLoan1.Principal != 500 || savedLoan2.Principal != 1000 {
		t.Errorf("failed to retrieve saved loans")
	}

	os.Remove("test_loans.json")
}

func TestLoanRepository_LoadEmptyFile(t *testing.T) {
	filePath := "empty.json"
	os.WriteFile(filePath, []byte(""), 0644)
	repo := NewLoanRepository(filePath)
	loans, err := repo.loadAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(loans) != 0 {
		t.Errorf("expected empty loans map, got %v", loans)
	}

	os.Remove(filePath)
}

func TestLoanRepository_LoadCorruptedFile(t *testing.T) {
	filePath := "corrupted.json"
	os.WriteFile(filePath, []byte("invalid json"), 0644)
	repo := NewLoanRepository(filePath)
	_, err := repo.loadAll()
	if err == nil {
		t.Errorf("expected error for corrupted file, got nil")
	}

	os.Remove(filePath)
}
