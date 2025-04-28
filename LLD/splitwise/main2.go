package main

import (
	"fmt"
)

// ========================= Models =========================

type User struct {
	ID   string
	Name string
}

type Expense struct {
	PaidBy       *User
	Amount       float64
	Participants []*User
	Splits       []Split
	Description  string
}

type Split struct {
	User   *User
	Amount float64
}

type Balance struct {
	OwedBy *User
	OwedTo *User
	Amount float64
}

// ========================= Strategy Pattern =========================

type SplitStrategy interface {
	CalculateSplits(paidBy *User, amount float64, participants []*User) []Split
}

type EqualSplit struct{}

func (e *EqualSplit) CalculateSplits(paidBy *User, amount float64, participants []*User) []Split {
	splitAmount := amount / float64(len(participants))
	splits := []Split{}
	for _, user := range participants {
		splits = append(splits, Split{
			User:   user,
			Amount: splitAmount,
		})
	}
	return splits
}

// ========================= Service =========================

type ExpenseService struct {
	users         map[string]*User
	expenses      []*Expense
	balances      []Balance
	splitStrategy SplitStrategy
}

func NewExpenseService() *ExpenseService {
	return &ExpenseService{
		users:         make(map[string]*User),
		expenses:      []*Expense{},
		balances:      []Balance{},
		splitStrategy: &EqualSplit{},
	}
}

func (es *ExpenseService) AddUser(id, name string) {
	es.users[id] = &User{ID: id, Name: name}
}

func (es *ExpenseService) AddExpense(paidByID string, amount float64, participantIDs []string, description string) {
	paidBy := es.users[paidByID]
	participants := []*User{}
	for _, id := range participantIDs {
		participants = append(participants, es.users[id])
	}

	splits := es.splitStrategy.CalculateSplits(paidBy, amount, participants)
	expense := &Expense{
		PaidBy:       paidBy,
		Amount:       amount,
		Participants: participants,
		Splits:       splits,
		Description:  description,
	}
	es.expenses = append(es.expenses, expense)

	// Update balances
	for _, split := range splits {
		if split.User.ID == paidBy.ID {
			continue
		}
		es.updateBalance(split.User, paidBy, split.Amount)
	}
}

func (es *ExpenseService) updateBalance(owedBy, owedTo *User, amount float64) {
	for i, bal := range es.balances {
		if bal.OwedBy.ID == owedBy.ID && bal.OwedTo.ID == owedTo.ID {
			es.balances[i].Amount += amount
			return
		}
		if bal.OwedBy.ID == owedTo.ID && bal.OwedTo.ID == owedBy.ID {
			es.balances[i].Amount -= amount
			// If amount becomes negative, swap owedBy and owedTo
			if es.balances[i].Amount < 0 {
				es.balances[i].Amount = -es.balances[i].Amount
				es.balances[i].OwedBy, es.balances[i].OwedTo = es.balances[i].OwedTo, es.balances[i].OwedBy
			}
			return
		}
	}
	// Not found, create new balance
	es.balances = append(es.balances, Balance{
		OwedBy: owedBy,
		OwedTo: owedTo,
		Amount: amount,
	})
}

func (es *ExpenseService) GetUserExpenses(userID string) []*Expense {
	result := []*Expense{}
	for _, exp := range es.expenses {
		if exp.PaidBy.ID == userID {
			result = append(result, exp)
		}
	}
	return result
}

func (es *ExpenseService) GetExpensesWhereUserIsPart(userID string) []*Expense {
	result := []*Expense{}
	for _, exp := range es.expenses {
		for _, participant := range exp.Participants {
			if participant.ID == userID {
				result = append(result, exp)
				break
			}
		}
	}
	return result
}

func (es *ExpenseService) GetUserBalance(userID string) {
	fmt.Println("Balances for user:", es.users[userID].Name)
	for _, bal := range es.balances {
		if bal.OwedBy.ID == userID {
			fmt.Printf("%s owes %s: %.2f\n", es.users[userID].Name, bal.OwedTo.Name, bal.Amount)
		}
		if bal.OwedTo.ID == userID {
			fmt.Printf("%s owes %s: %.2f\n", bal.OwedBy.Name, es.users[userID].Name, bal.Amount)
		}
	}
}

// ========================= Main to Test =========================

func main() {
	expenseService := NewExpenseService()

	expenseService.AddUser("u1", "Alice")
	expenseService.AddUser("u2", "Bob")
	expenseService.AddUser("u3", "Charlie")

	expenseService.AddExpense("u1", 120, []string{"u1", "u2", "u3"}, "Lunch")
	expenseService.AddExpense("u2", 60, []string{"u2", "u3"}, "Taxi")

	fmt.Println("User Expenses for Alice:")
	for _, exp := range expenseService.GetUserExpenses("u1") {
		fmt.Println(exp.Description, exp.Amount)
	}

	fmt.Println("\nExpenses where Bob is participant:")
	for _, exp := range expenseService.GetExpensesWhereUserIsPart("u2") {
		fmt.Println(exp.Description, exp.Amount)
	}

	fmt.Println("\nBalances:")
	expenseService.GetUserBalance("u1")
	expenseService.GetUserBalance("u2")
	expenseService.GetUserBalance("u3")
}
