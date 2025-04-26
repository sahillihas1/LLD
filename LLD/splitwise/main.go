package main

type SplitWiseService struct {
	expenseService IExpenseService
	groupService   IGroupService
}

func main() {
	splitWiseService := &SplitWiseService{
		expenseService: &ExpenseService{},
		groupService:   &GroupService{},
	}
	splitWiseService.expenseService.AddExpense(&AddExpenseRequest{})
	splitWiseService.groupService.GetGroupPaymentGraph("1")
}
