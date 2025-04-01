package billing

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
)

type Loan struct {
	mu                sync.Mutex
	ID                string
	Principal         float64
	Rate              float64
	Weeks             int
	Outstanding       float64
	LoanAmount        float64
	WeeklyInstallment float64
	LastPaidWeek      int
	Delinquent        bool
	Schedule          []PaymentSchedule
	Payments          map[int]float64
}

type PaymentSchedule struct {
	Week   int
	Status string
}

func GenerateLoanID() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func CalculateInitialOutstanding(principal, rate float64) float64 {
	interest := principal * (rate / 100)
	return principal + interest
}

func WeeklyPayment(loanAmount float64, weeks int) float64 {
	return loanAmount / float64(weeks)
}

func NewLoan(principal, rate float64, weeks int) (*Loan, error) {
	loanID, err := GenerateLoanID()
	if err != nil {
		return nil, err
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
		Schedule:          make([]PaymentSchedule, weeks),
		Payments:          make(map[int]float64),
	}
	for i := 0; i < weeks; i++ {
		loan.Schedule[i] = PaymentSchedule{Week: i + 1, Status: "Pending"}
	}
	return loan, nil
}

func (l *Loan) ApplyPayment(week int, amount float64) error {
	l.mu.Lock()

	if week < 1 || week > l.Weeks {
		l.mu.Unlock()
		return fmt.Errorf("invalid payment week, loan tenure was upto %d weeks", l.Weeks)
	}
	if l.Outstanding == 0 {
		l.mu.Unlock()
		return errors.New("loan fully paid")
	}
	if l.WeeklyInstallment != amount {
		l.mu.Unlock()
		return fmt.Errorf("invalid payment amount, expected %.2f", l.WeeklyInstallment)
	}
	l.Payments[week] = amount
	l.Outstanding -= amount

	for i := range l.Schedule {
		if l.Schedule[i].Week == week {
			l.Schedule[i].Status = "Paid"
			break
		}
	}
	l.mu.Unlock()
	l.Delinquent, _ = l.IsDelinquent(week)
	return nil
}

func (l *Loan) OutstandingAmount() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Outstanding
}

func (l *Loan) GenerateSchedule() []PaymentSchedule {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Schedule
}

func (l *Loan) IsDelinquent(currentWeek int) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	missedPayments := 0
	for i := range l.Schedule {
		if l.Schedule[i].Week >= currentWeek-2 && l.Schedule[i].Week < currentWeek {
			if l.Schedule[i].Status == "Pending" {
				missedPayments++
			} else {
				missedPayments = 0
			}
		}
	}
	l.Delinquent = missedPayments >= 2
	return l.Delinquent, nil
}
