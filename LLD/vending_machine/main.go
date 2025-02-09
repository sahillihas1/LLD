package main

import (
	"errors"
	"fmt"
	"sync"
)

// Product represents a product in the vending machine
type Product struct {
	Name     string
	Price    int
	Quantity int
}

// VendingMachine represents the vending machine
type VendingMachine struct {
	Products      map[string]*Product // product name -> quantity
	Balance       int                 // current balance in the machine
	State         VendingMachineState
	PaymentMethod PaymentStrategy
	mu            sync.Mutex
}

// VendingMachineState defines the interface for vending machine states
type VendingMachineState interface {
	SelectProduct(vm *VendingMachine, productName string) error
	InsertMoney(vm *VendingMachine, amount int) error
	DispenseProduct(vm *VendingMachine) error
}

// IdleState represents the idle state of the vending machine
type IdleState struct{}

func (i *IdleState) SelectProduct(vm *VendingMachine, productName string) error {
	if _, exists := vm.Products[productName]; !exists {
		return errors.New("product not found")
	}
	if vm.Products[productName].Quantity <= 0 {
		return errors.New("product out of stock")
	}
	vm.State = &ProcessingState{SelectedProduct: productName}
	return nil
}

func (i *IdleState) InsertMoney(vm *VendingMachine, amount int) error {
	return errors.New("please select a product first")
}

func (i *IdleState) DispenseProduct(vm *VendingMachine) error {
	return errors.New("please select a product first")
}

// ProcessingState represents the state when a product is selected
type ProcessingState struct {
	SelectedProduct string
}

func (p *ProcessingState) SelectProduct(vm *VendingMachine, productName string) error {
	return errors.New("already processing a product")
}

func (p *ProcessingState) InsertMoney(vm *VendingMachine, amount int) error {
	if vm.PaymentMethod == nil {
		return errors.New("no payment method selected")
	}
	err := vm.PaymentMethod.Pay(amount)
	if err != nil {
		return err
	}
	vm.Balance += amount
	if vm.Balance >= vm.Products[p.SelectedProduct].Price {
		vm.State = &DispensingState{}
	}
	return nil
}

func (p *ProcessingState) DispenseProduct(vm *VendingMachine) error {
	return errors.New("please insert more money")
}

// DispensingState represents the state when the product is being dispensed
type DispensingState struct{}

func (d *DispensingState) SelectProduct(vm *VendingMachine, productName string) error {
	return errors.New("currently dispensing a product")
}

func (d *DispensingState) InsertMoney(vm *VendingMachine, amount int) error {
	return errors.New("currently dispensing a product")
}

func (d *DispensingState) DispenseProduct(vm *VendingMachine) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	productPrice := vm.Products[vm.State.(*ProcessingState).SelectedProduct].Price
	if vm.Balance < productPrice {
		return errors.New("insufficient funds")
	}

	vm.Products[vm.State.(*ProcessingState).SelectedProduct].Quantity--
	vm.Balance -= productPrice
	fmt.Printf("Dispensing %s\n", vm.State.(*ProcessingState).SelectedProduct)

	if vm.Balance > 0 {
		fmt.Printf("Returning change: %d\n", vm.Balance)
		vm.Balance = 0
	}

	vm.State = &IdleState{}
	return nil
}

// PaymentStrategy defines the interface for payment methods
type PaymentStrategy interface {
	Pay(amount int) error
}

// CoinPayment represents payment using coins
type CoinPayment struct{}

func (c *CoinPayment) Pay(amount int) error {
	fmt.Printf("Paid %d using coins\n", amount)
	return nil
}

// NotePayment represents payment using notes
type NotePayment struct{}

func (n *NotePayment) Pay(amount int) error {
	fmt.Printf("Paid %d using notes\n", amount)
	return nil
}

// CardPayment represents payment using a card
type CardPayment struct{}

func (c *CardPayment) Pay(amount int) error {
	fmt.Printf("Paid %d using card\n", amount)
	return nil
}

// VendingMachineService implements the business logic for the vending machine
type VendingMachineService struct {
	vm *VendingMachine
}

// NewVendingMachineService creates a new service
func NewVendingMachineService(vm *VendingMachine) *VendingMachineService {
	return &VendingMachineService{vm: vm}
}

// SelectProduct selects a product
func (s *VendingMachineService) SelectProduct(productName string) error {
	return s.vm.State.SelectProduct(s.vm, productName)
}

// InsertMoney inserts money into the vending machine using the selected payment method
func (s *VendingMachineService) InsertMoney(amount int, paymentMethod PaymentStrategy) error {
	s.vm.PaymentMethod = paymentMethod
	return s.vm.State.InsertMoney(s.vm, amount)
}

// DispenseProduct dispenses the selected product
func (s *VendingMachineService) DispenseProduct() error {
	return s.vm.State.DispenseProduct(s.vm)
}

// Restock adds more products to the vending machine
func (s *VendingMachineService) Restock(productName string, quantity int) {
	s.vm.mu.Lock()
	defer s.vm.mu.Unlock()
	s.vm.Products[productName].Quantity += quantity
}

// CollectMoney retrieves the money from the vending machine
func (s *VendingMachineService) CollectMoney() int {
	s.vm.mu.Lock()
	defer s.vm.mu.Unlock()
	money := s.vm.Balance
	s.vm.Balance = 0
	return money
}

func main() {
	// Initialize vending machine
	vm := &VendingMachine{
		Products: make(map[string]*Product),
		Balance:  0,
		State:    &IdleState{},
	}

	// Add products
	vm.Products["Coke"] = &Product{Name: "Coke", Price: 10, Quantity: 10}
	vm.Products["Pepsi"] = &Product{Name: "Pepsi", Price: 15, Quantity: 10}
	// Initialize service
	vmService := NewVendingMachineService(vm)

	// Simulate a transaction
	err := vmService.SelectProduct("Coke")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Use CoinPayment strategy
	err = vmService.InsertMoney(20, &CoinPayment{})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = vmService.DispenseProduct()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Restock products
	vmService.Restock("Coke", 10)

	// Collect money
	money := vmService.CollectMoney()
	fmt.Printf("Collected money: %d\n", money)
}
