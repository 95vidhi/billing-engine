package billing

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Loan struct {
	ID                string        `json:"id"`
	Principal         float64       `json:"principal"`
	Rate              float64       `json:"rate"`
	Weeks             int           `json:"weeks"`
	Outstanding       float64       `json:"outstanding"`
	LoanAmount        float64       `json:"loanAmount"`
	WeeklyInstallment float64       `json:"weeklyInstallment"`
	LastPaidWeek      int           `json:"last_paid_week"`
	Delinquent        bool          `json:"delinquent"`
	Schedule          []PaymentWeek `json:"schedule"`
}

type PaymentWeek struct {
	Week   int    `json:"week"`
	Status string `json:"status"`
}

type LoanStore struct {
	Loans map[string]*Loan
	mu    sync.Mutex
}

// LoadLoans loads loans from a JSON file
func LoadLoans(filename string) (*LoanStore, error) {
	store := &LoanStore{Loans: make(map[string]*Loan)}

	file, err := os.Open(filename)
	if err != nil {
		return store, err
	}
	defer file.Close()

	store.mu.Lock()
	defer store.mu.Unlock()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&store.Loans)
	if err != nil {
		return nil, err
	}

	return store, nil
}

// SaveLoans saves all loans to a JSON file
func (s *LoanStore) SaveLoans(filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(s.Loans)
	return err
}

// GenerateLoanID generates a unique random loan ID
func GenerateLoanID() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	id := hex.EncodeToString(bytes)
	return id, nil
}

// CalculateInitialOutstanding calculates the initial outstanding balance
func CalculateInitialOutstanding(principal float64, annualInterestRate float64) float64 {
	interest := principal * (annualInterestRate / 100) // Flat interest for one year
	totalLoan := principal + interest                  // Total amount to be repaid
	return totalLoan
}

// NewLoan creates a new loan with minimal storage requirements
func (s *LoanStore) NewLoan(principal, rate float64, weeks int, transactionFilePath string) (string, error) {
	s.mu.Lock()

	loanID, err := GenerateLoanID()
	if err != nil {
		s.mu.Unlock()
		return "", err
	}
	if _, exists := s.Loans[loanID]; exists {
		s.mu.Unlock()
		return loanID, nil
	}

	loanAmount := CalculateInitialOutstanding(principal, rate)
	loan := &Loan{
		ID:                loanID,
		Principal:         principal,
		Rate:              rate,
		Weeks:             weeks,
		LoanAmount:        loanAmount,
		Outstanding:       loanAmount,
		WeeklyInstallment: WeeklyPayment(loanAmount, weeks),
		LastPaidWeek:      0,
		Delinquent:        false,
		Schedule:          make([]PaymentWeek, weeks),
	}

	for i := 0; i < weeks; i++ {
		loan.Schedule[i] = PaymentWeek{Week: i + 1, Status: "Pending"}
	}
	s.Loans[loanID] = loan
	s.mu.Unlock()
	err = s.SaveLoans(transactionFilePath)
	return loanID, err
}

// MakePayment processes a payment for a specific loan
func (s *LoanStore) MakePayment(id string, currentWeek int, transactionFilePath string) error {
	s.mu.Lock()

	loan, exists := s.Loans[id]
	if !exists {
		s.mu.Unlock()
		return errors.New("loan not found")
	}

	if loan.LastPaidWeek >= loan.Weeks {
		s.mu.Unlock()
		return errors.New("loan fully paid")
	}

	for i := range loan.Schedule {
		if loan.Schedule[i].Week == currentWeek {
			loan.Schedule[i].Status = "Paid"
			break
		}
	}

	loan.LastPaidWeek = currentWeek
	loan.Outstanding -= loan.WeeklyInstallment
	s.mu.Unlock()

	loan.Delinquent, _ = s.IsDelinquent(id, currentWeek)

	return s.SaveLoans(transactionFilePath)
}

// WeeklyPayment dynamically calculates weekly payment amount
func WeeklyPayment(loanAmount float64, weeks int) float64 {
	return loanAmount / float64(weeks)
}

// GetOutstanding returns the outstanding amount for a loan
func (s *LoanStore) GetOutstanding(loanID string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	loan, exists := s.Loans[loanID]
	if !exists {
		return 0, errors.New("Loan ID not found")
	}

	return loan.Outstanding, nil
}

// IsDelinquent checks if the borrower has missed two consecutive payments
func (s *LoanStore) IsDelinquent(loanID string, currentWeek int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	loan, exists := s.Loans[loanID]
	if !exists {
		return false, errors.New("Loan ID not found")
	}

	missedPayments := 0
	for i := range loan.Schedule {
		if loan.Schedule[i].Week >= currentWeek-2 && loan.Schedule[i].Week < currentWeek {
			if loan.Schedule[i].Status == "Pending" {
				missedPayments++
			} else {
				missedPayments = 0
			}
		}
	}

	loan.Delinquent = missedPayments >= 2
	return loan.Delinquent, nil
}

func (s *LoanStore) GeneratePaymentSchedule(loanID string) ([]PaymentWeek, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	loan, exists := s.Loans[loanID]
	if !exists {
		return nil, fmt.Errorf("loan not found")
	}

	return loan.Schedule, nil
}
