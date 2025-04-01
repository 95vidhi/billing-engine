package billing

import (
	"testing"
)

func TestLoan_ApplyPaymentEdgeCases(t *testing.T) {
	loan, _ := NewLoan(1000, 5, 10)

	tests := []struct {
		name       string
		week       int
		amount     float64
		expectsErr bool
	}{
		{"Invalid week (before start)", 0, 105, true},
		{"Invalid week (after tenure)", 11, 105, true},
		{"Incorrect payment amount", 2, 200, true},
		{"Valid payment", 3, 105, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loan.ApplyPayment(tt.week, tt.amount)
			if (err != nil) != tt.expectsErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestLoan_OutstandingAmount(t *testing.T) {
	loan, _ := NewLoan(1000, 5, 10)
	loan.ApplyPayment(1, 105)
	outstanding := loan.OutstandingAmount()
	if outstanding != 945 {
		t.Errorf("expected outstanding 945, got %v", outstanding)
	}
}

func TestLoan_GenerateSchedule(t *testing.T) {
	loan, _ := NewLoan(1000, 5, 10)
	schedule := loan.GenerateSchedule()
	if len(schedule) != 10 {
		t.Errorf("expected schedule for 10 weeks, got %d", len(schedule))
	}
	for i, entry := range schedule {
		if entry.Week != i+1 || entry.Status != "Pending" {
			t.Errorf("unexpected schedule entry at week %d: %+v", i+1, entry)
		}
	}
}

func TestLoan_IsDelinquent(t *testing.T) {
	loan, _ := NewLoan(1000, 5, 10)
	loan.ApplyPayment(1, 105)
	loan.ApplyPayment(4, 105)

	delinquent, _ := loan.IsDelinquent(4)
	if !delinquent {
		t.Errorf("expected delinquent status, got false")
	}

	loan.ApplyPayment(2, 105)
	delinquent, _ = loan.IsDelinquent(4)
	if delinquent {
		t.Errorf("expected non-delinquent status, got true")
	}
}
