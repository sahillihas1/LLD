package main

import "errors"

type AddExpenseRequest struct {
	Title       string
	Amount      float64
	PaidBy      string
	SplitType   SplitType
	GroupId     string
	Description string
	Map         map[string]Amount
	TotalAmount Amount
}

type IExpenseService interface {
	AddExpense(request *AddExpenseRequest) error
	GetPaymentGraph(expense map[string]float64) (map[string]map[string]float64, error)
}

type ExpenseService struct {
	expenseRepo IExpenseRepo
	groupRepo   IGroupRepo
}

func (service *ExpenseService) AddExpense(expenseRequest *AddExpenseRequest) error {
	group, err := service.groupRepo.GetGroupById(expenseRequest.GroupId)
	if group != nil || err != nil {
		return errors.New("group not found")
	}
	expense := &Expense{
		ID:           "2",
		Title:        expenseRequest.Title,
		Description:  expenseRequest.Description,
		GroupId:      expenseRequest.GroupId,
		userBalances: expenseRequest.Map,
		TotalAmount:  expenseRequest.TotalAmount,
		SplitType:    expenseRequest.SplitType,
	}
	if err := service.expenseRepo.AddExpense(expense); err != nil {
		return err
	}
	return nil
}

func (service *ExpenseService) GetPaymentGraph(expense map[string]float64) (map[string]map[string]float64, error) {

}
