```markdown
# Billing Engine

## Overview

This project provides a **Billing Engine** for managing loans. It supports features such as loan creation, payments, generating payment schedules, calculating outstanding amounts, and handling loan delinquency. The app uses Go structs and JSON files for persistence, ensuring a simple and lightweight design.

## How Loans Work

- Loans are created with a fixed principal amount, an interest rate, and a number of weeks (e.g., 50 weeks).
- The interest on the loan is calculated using the formula:  
  `Total Loan Amount = Principal + (Principal * Rate / 100)`
- The loan is divided into equal weekly installments, and the amount to be paid weekly is calculated using:  
  `Weekly Installment = Loan Amount / Number of Weeks`
- Payments can be made weekly. Each payment will reduce the outstanding balance by the weekly installment amount.
- Loans can become delinquent if payments are missed for two consecutive weeks. Delinquency status is checked during each payment.

## Features

- **Loan Creation:** Create a loan by specifying principal, interest rate, and number of weeks.
- **Payments:** Apply payments to the loan, reducing the outstanding balance.
- **Delinquency Tracking:** Track if the loan is delinquent due to missed payments.
- **Outstanding Balance:** Get the remaining amount on the loan after payments.
- **Payment Schedule:** Generate a payment schedule showing when payments are due and their status.

## Sample Transactions

Here is an example of how the engine works:

1. **Create a Loan**

   ```go
   loan, err := NewLoan(1000, 5, 50)
   if err != nil {
       log.Fatalf("Error creating loan: %v", err)
   }
   fmt.Println("Loan created successfully!")
   ```

2. **Make a Payment**

   Make a payment for week 1:

   ```go
   err = loan.ApplyPayment(1, 20)
   if err != nil {
       log.Fatalf("Error making payment: %v", err)
   }
   fmt.Println("Payment applied successfully!")
   ```

3. **Check Outstanding Amount**

   Check the remaining balance on the loan:

   ```go
   outstanding := loan.OutstandingAmount()
   fmt.Printf("Outstanding balance: %.2f\n", outstanding)
   ```

4. **Check Delinquency Status**

   Check if the loan is delinquent:

   ```go
   isDelinquent, err := loan.IsDelinquent(2)
   if err != nil {
       log.Fatalf("Error checking delinquency: %v", err)
   }
   fmt.Printf("Delinquent status: %v\n", isDelinquent)
   ```

## How to Run the App

1. **Clone the repository**  
   Clone this repository to your local machine:

   ```bash
   git clone https://github.com/95vidhi/billing-engine
   cd billing-engine
   ```

2. **Run the app**  
   You can run the Go app directly from the command line:

   ```bash
   go run main.go
   ```

   This will run the application and you can interact with the billing engine.

## How to Run Tests

To run all the tests for the project, use the following command:

```bash
go test ./...
```

You can also run tests with coverage:

```bash
go test -cover ./...
```

## Design Philosophy

- **Simplicity:** The application uses simple Go structs for loan management and a JSON file for persistence, making it easy to extend and maintain.
- **Extensibility:** The design is modular, with separate components for managing loans, payments, and delinquency.
- **Test-Driven Development:** Comprehensive tests are written for every major functionality to ensure the correctness of the system.

## Assumptions

- No pre-payments are allowed.
- Partial payments are not supported; full weekly payments are required.
- The system assumes loans are paid weekly until fully repaid.

```
