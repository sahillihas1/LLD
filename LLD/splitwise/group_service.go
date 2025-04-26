package main

import "errors"

type IGroupService interface {
	GetGroupPaymentGraph(groupId string) (map[string]map[string]float64, error)
}

type GroupService struct {
	groupRepo      IGroupRepo
	expenseRepo    IExpenseRepo
	expenseService IExpenseService
}

func (service *GroupService) GetGroupPaymentGraph(groupId string) (map[string]map[string]float64, error) {
	group, err := service.groupRepo.GetGroupById(groupId)
	if group == nil || err != nil {
		return nil, errors.New("group not found")
	}
	expenses, err := service.expenseRepo.GetExpenseByGroupId(groupId)
	groupBalance := make(map[string]float64)
	for _, expense := range expenses {
		for user, amount := range expense.userBalances {
			groupBalance[user] += amount.Value
		}
	}
	return service.expenseService.GetPaymentGraph(groupBalance)
}
