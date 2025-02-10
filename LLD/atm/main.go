package main

import (
	"fmt"
)

// --- Account Interface ---
type Account interface {
	Withdraw(amount float64) error
	Deposit(amount float64)
	GetBalance() float64
}

// --- Concrete Account Implementations ---
type SavingsAccount struct {
	balance float64
}

func (s *SavingsAccount) Withdraw(amount float64) error {
	if s.balance < amount {
		return fmt.Errorf("insufficient funds")
	}
	s.balance -= amount
	return nil
}

func (s *SavingsAccount) Deposit(amount float64) {
	s.balance += amount
}

func (s *SavingsAccount) GetBalance() float64 {
	return s.balance
}

// --- Account Factory ---
type AccountFactory struct{}

func (f *AccountFactory) CreateAccount(accountType string, initialBalance float64) Account {
	switch accountType {
	case "savings":
		return &SavingsAccount{balance: initialBalance}
	default:
		return nil
	}
}

// --- Strategy Pattern for Transactions ---
type TransactionStrategy interface {
	Execute(account Account, amount float64) error
}

type WithdrawStrategy struct{}

func (w *WithdrawStrategy) Execute(account Account, amount float64) error {
	return account.Withdraw(amount)
}

type DepositStrategy struct{}

func (d *DepositStrategy) Execute(account Account, amount float64) error {
	account.Deposit(amount)
	return nil
}

// --- State Pattern for ATM ---
type ATMState interface {
	InsertCard(atm *ATM, account Account)
	EjectCard(atm *ATM)
	EnterPin(atm *ATM, pin int)
	RequestTransaction(account Account, requestType string, amount float64)
}

// Idle State
type IdleState struct{}

func (i *IdleState) InsertCard(atm *ATM, account Account) {
	fmt.Println("Card Inserted. Please enter PIN.")
	atm.SetState(&HasCardState{
		Account: account,
	})
}
func (i *IdleState) EjectCard(atm *ATM) {
	fmt.Println("No card to eject.")
}
func (i *IdleState) EnterPin(atm *ATM, pin int) {
	fmt.Println("Insert card first.")
}
func (i *IdleState) RequestTransaction(account Account, requestType string, amount float64) {
	fmt.Println("Insert card first.")
}

// Has Card State
type HasCardState struct {
	Account Account
}

func (h *HasCardState) InsertCard(atm *ATM, account Account) {
	fmt.Println("Card already inserted.")
}
func (h *HasCardState) EjectCard(atm *ATM) {
	fmt.Println("Card Ejected.")
}
func (h *HasCardState) EnterPin(atm *ATM, pin int) {
	fmt.Println("PIN accepted. You may proceed with a transaction.")
}
func (h *HasCardState) RequestTransaction(account Account, requestType string, amount float64) {
	fmt.Println("Enter PIN first.")
}

// Pin Entered State

type IAtmProcessFactory interface {
	CreateProcess(requestType string, amount float64) IAtmProcessExecute
}

type AtmProcessFactory struct{}

func (f *AtmProcessFactory) CreateProcess(requestType string, amount float64) IAtmProcessExecute {
	switch requestType {
	case "withdraw":
		return &WithdrawProcess{
			amount: amount,
		}
	case "deposit":
		return &DepositProcess{
			amount: amount,
		}
	case "check balance":
		return &CheckBalanceProcess{}
	}
	return nil
}

type IAtmProcessExecute interface {
	Execute(account Account) error
}

type WithdrawProcess struct {
	amount float64
}

func (w *WithdrawProcess) Execute(account Account) error {
	return nil
}

type DepositProcess struct {
	amount float64
}

func (d *DepositProcess) Execute(account Account) error {
	return nil
}

type CheckBalanceProcess struct{}

func (c *CheckBalanceProcess) Execute(account Account) error {
	return nil
}

type PinEnteredState struct {
	atmProcessFactory IAtmProcessFactory
}

func (p *PinEnteredState) InsertCard(atm *ATM, account Account) {
	fmt.Println("Card already inserted.")
}
func (p *PinEnteredState) EjectCard(atm *ATM) {
	fmt.Println("Card Ejected.")
	atm.SetState(&IdleState{})
}
func (p *PinEnteredState) EnterPin(atm *ATM, pin int) {
	fmt.Println("PIN already entered.")
}
func (p *PinEnteredState) RequestTransaction(account Account, requestType string, amount float64) {
	process := p.atmProcessFactory.CreateProcess(requestType, amount)
	process.Execute(account)
}

// ATM Context
type ATM struct {
	state ATMState
}

func (a *ATM) SetState(state ATMState) {
	a.state = state
}
func (a *ATM) InsertCard(account Account) {
	a.state.InsertCard(a, account)
	a.state = &HasCardState{Account: account}
}
func (a *ATM) EjectCard() {
	a.state.EjectCard(a)
	a.state = &IdleState{}
}
func (a *ATM) EnterPin(pin int) {
	a.state.EnterPin(a, pin)
	a.state = &PinEnteredState{}
}
func (a *ATM) RequestTransaction(account Account, requestType string, amount float64) {
	a.state.RequestTransaction(account, requestType, amount)
	a.state = &HasCardState{}
}

func main() {
	factory := &AccountFactory{}
	account := factory.CreateAccount("savings", 1000)
	atm := &ATM{state: &IdleState{}}

	atm.InsertCard(account)
	atm.EnterPin(1234)
	atm.RequestTransaction(account, "withdraw", 500)
	atm.EjectCard()
}
