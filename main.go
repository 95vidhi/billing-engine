package main

import (
	"billing-engine/billing"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

const TransactionFilePath = "transactions.json"

func main() {
	repo := billing.NewLoanRepository(TransactionFilePath)
	engine := billing.NewBillingEngine(repo)

	runCLI(engine, os.Stdin, os.Stdout)
}

func runCLI(engine *billing.BillingEngine, input *os.File, output *os.File) {
	scanner := bufio.NewScanner(input)
	for {
		fmt.Fprintln(output, "\n***************************************************************")
		fmt.Fprintln(output, "Select an option:")
		fmt.Fprintln(output, "1. Create New Loan")
		fmt.Fprintln(output, "2. Get Outstanding Amount")
		fmt.Fprintln(output, "3. Generate Payment Schedule")
		fmt.Fprintln(output, "4. Make Payment")
		fmt.Fprintln(output, "5. Check if borrower is delinquent")
		fmt.Fprintln(output, "6. Exit")
		fmt.Fprint(output, "Enter choice: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			createNewLoan(engine, scanner, output)
		case "2":
			getOutstanding(engine, scanner, output)
		case "3":
			generateSchedule(engine, scanner, output)
		case "4":
			makePayment(engine, scanner, output)
		case "5":
			checkDelinquency(engine, scanner, output)
		case "6":
			fmt.Fprintln(output, "Exiting...")
			return
		default:
			fmt.Fprintln(output, "Invalid choice, try again.")
		}
	}
}

func createNewLoan(engine *billing.BillingEngine, scanner *bufio.Scanner, output *os.File) {
	fmt.Fprint(output, "Enter principal amount: ")
	scanner.Scan()
	principal, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Fprint(output, "Enter interest rate: ")
	scanner.Scan()
	rate, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Fprint(output, "Enter number of weeks: ")
	scanner.Scan()
	weeks, _ := strconv.Atoi(scanner.Text())

	loanID, err := engine.NewLoan(principal, rate, weeks)
	if err != nil {
		fmt.Fprintln(output, "Error creating loan:", err)
		return
	}
	fmt.Fprintln(output, "New Loan ID:", loanID)
}

func getOutstanding(engine *billing.BillingEngine, scanner *bufio.Scanner, output io.Writer) {
	fmt.Fprint(output, "Enter Loan ID: ")
	scanner.Scan()
	loanID := scanner.Text()

	outstanding, err := engine.GetOutstanding(loanID)
	if err != nil {
		fmt.Fprintln(output, "Error:", err)
		return
	}
	fmt.Fprintf(output, "Outstanding amount: %.2f\n", outstanding)
}

func generateSchedule(engine *billing.BillingEngine, scanner *bufio.Scanner, output io.Writer) {
	fmt.Fprint(output, "Enter Loan ID: ")
	scanner.Scan()
	loanID := scanner.Text()

	schedule, err := engine.GeneratePaymentSchedule(loanID)
	if err != nil {
		fmt.Fprintln(output, "Error:", err)
		return
	}
	for _, week := range schedule {
		fmt.Printf("Week %d: %s\n", week.Week, week.Status)
	}
}

func makePayment(engine *billing.BillingEngine, scanner *bufio.Scanner, output *os.File) {
	fmt.Fprint(output, "Enter Loan ID: ")
	scanner.Scan()
	loanID := scanner.Text()

	fmt.Fprint(output, "Enter payment week: ")
	scanner.Scan()
	paymentWeek, _ := strconv.Atoi(scanner.Text())

	fmt.Fprint(output, "Enter amount: ")
	scanner.Scan()
	amount, _ := strconv.ParseFloat(scanner.Text(), 64)

	err := engine.MakePayment(loanID, paymentWeek, amount)
	if err != nil {
		fmt.Fprintln(output, "Error making payment:", err)
		return
	}
	fmt.Fprintln(output, "Payment successful!")
}

func checkDelinquency(engine *billing.BillingEngine, scanner *bufio.Scanner, output io.Writer) {
	fmt.Fprint(output, "Enter Loan ID: ")
	scanner.Scan()
	loanID := scanner.Text()

	fmt.Fprint(output, "Enter current week: ")
	scanner.Scan()
	currentWeek, _ := strconv.Atoi(scanner.Text())

	isDelinquent, err := engine.IsDelinquent(loanID, currentWeek)
	if err != nil {
		fmt.Fprintln(output, "Error checking delinquency:", err)
	} else if isDelinquent {
		fmt.Fprintln(output, "Borrower is delinquent!")
	} else {
		fmt.Fprintln(output, "Borrower is not delinquent.")
	}
}
