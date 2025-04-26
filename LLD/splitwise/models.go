package main

type Amount struct {
	Value float64
}

type Expense struct {
	ID           string
	Title        string
	ImageUrl     string
	Description  string
	GroupId      string
	userBalances map[string]Amount //user to balance
	TotalAmount  Amount
	SplitType    SplitType
}

type SplitType int

const (
	EQUALLY SplitType = iota
	UNEQUALLY
	BYPERCENTAGE
)

type BalanceMap struct {
	Balances map[string]Amount
}

type Group struct {
	ID      string
	Name    string
	Members []*User
}

type User struct {
	ID   string
	Name string
}
