package main

import (
	"fmt"
)

// ----------------- Models ------------------

// User represents a user in the system
type User struct {
	ID   string
	Name string
}

// Expense represents an expense
// It can have different split strategies
type Expense struct {
	PaidBy       *User
	Amount       float64
	Participants []*User
	Splits       []Split
	Description  string
}

// Split represents how much a user owes
// This can be extended for unequal splits, percentage splits, etc.
type Split interface {
	GetUser() *User
	GetAmount() float64
}

type EqualSplit struct {
	User   *User
	Amount float64
}

func (es *EqualSplit) GetUser() *User {
	return es.User
}

func (es *EqualSplit) GetAmount() float64 {
	return es.Amount
}

// ----------------- Service Layer ------------------

type ExpenseService struct {
	expenses     []*Expense
	userExpenses map[string][]*Expense
	balances     map[string]map[string]float64 // balances[user1][user2] = amount user1 owes user2
}

func NewExpenseService() *ExpenseService {
	return &ExpenseService{
		expenses:     []*Expense{},
		userExpenses: make(map[string][]*Expense),
		balances:     make(map[string]map[string]float64),
	}
}

func (es *ExpenseService) AddExpense(paidBy *User, amount float64, participants []*User, description string) {
	splits := es.createEqualSplits(participants, amount)
	expense := &Expense{
		PaidBy:       paidBy,
		Amount:       amount,
		Participants: participants,
		Splits:       splits,
		Description:  description,
	}
	es.expenses = append(es.expenses, expense)

	for _, participant := range participants {
		// Attach expense to user
		es.userExpenses[participant.ID] = append(es.userExpenses[participant.ID], expense)
	}

	// Update balances
	for _, split := range splits {
		if split.GetUser().ID != paidBy.ID {
			if es.balances[split.GetUser().ID] == nil {
				es.balances[split.GetUser().ID] = make(map[string]float64)
			}
			es.balances[split.GetUser().ID][paidBy.ID] += split.GetAmount()
		}
	}
}

func (es *ExpenseService) createEqualSplits(participants []*User, amount float64) []Split {
	splits := []Split{}
	share := amount / float64(len(participants))
	for _, user := range participants {
		splits = append(splits, &EqualSplit{
			User:   user,
			Amount: share,
		})
	}
	return splits
}

func (es *ExpenseService) GetAllExpensesByUser(userID string) []*Expense {
	return es.userExpenses[userID]
}

func (es *ExpenseService) GetBalance(userID string) map[string]float64 {
	return es.balances[userID]
}

func (es *ExpenseService) GetSimplifiedBalance(userID string) {
	fmt.Printf("Simplified Balance for %s:\n", userID)
	for otherUser, amount := range es.balances[userID] {
		if amount > 0 {
			fmt.Printf("%s owes %s: %.2f\n", userID, otherUser, amount)
		}
	}
}

// ----------------- Main ------------------

func main() {
	service := NewExpenseService()

	a := &User{ID: "u1", Name: "Alice"}
	b := &User{ID: "u2", Name: "Bob"}
	c := &User{ID: "u3", Name: "Charlie"}

	service.AddExpense(a, 300, []*User{a, b, c}, "Dinner")
	service.AddExpense(b, 150, []*User{b, c}, "Snacks")

	fmt.Println("--- Expenses by Alice ---")
	for _, e := range service.GetAllExpensesByUser("u1") {
		fmt.Printf("Expense: %s Amount: %.2f\n", e.Description, e.Amount)
	}

	fmt.Println("--- Balances for Bob ---")
	fmt.Println(service.GetBalance("u2"))

	fmt.Println("--- Simplified Balance for Charlie ---")
	service.GetSimplifiedBalance("u3")
}
