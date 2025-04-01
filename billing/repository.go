package billing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type LoanRepository interface {
	Save(loan *Loan) (string, error)
	FindByID(id string) (*Loan, error)
}

type JSONLoanRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewLoanRepository(filePath string) *JSONLoanRepository {
	return &JSONLoanRepository{
		filePath: filePath,
	}
}

func (repo *JSONLoanRepository) Save(loan *Loan) (string, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	loans, err := repo.loadAll()
	if err != nil {
		return "", err
	}

	loans[loan.ID] = loan
	return loan.ID, repo.saveAll(loans)
}

func (repo *JSONLoanRepository) saveAll(loans map[string]*Loan) error {
	file, err := os.Create(repo.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(loans)
}

func (repo *JSONLoanRepository) loadAll() (map[string]*Loan, error) {
	file, err := os.Open(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*Loan), nil
		}
		return nil, err
	}
	defer file.Close()

	loans := make(map[string]*Loan)
	err = json.NewDecoder(file).Decode(&loans)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return make(map[string]*Loan), nil
		}
		return nil, err
	}

	return loans, nil
}

func (repo *JSONLoanRepository) FindByID(id string) (*Loan, error) {
	loans, err := repo.loadAll()
	if err != nil {
		return nil, err
	}

	loan, exists := loans[id]
	if !exists {
		return nil, fmt.Errorf("loan with ID %s not found", id)
	}

	return loan, nil
}
