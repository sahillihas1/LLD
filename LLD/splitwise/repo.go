package main

type IExpenseRepo interface {
	AddExpense(expense *Expense) error
	GetExpenseByGroupId(groupId string) ([]*Expense, error)
	GetExpenseById(expenseId string) (*Expense, error)
}

type ExpenseRepo struct {
	expenses map[string]*Expense
}

func (repo *ExpenseRepo) AddExpense(expense *Expense) error {
	repo.expenses[expense.ID] = expense
	return nil
}

func (repo *ExpenseRepo) GetExpenseByGroupId(groupId string) ([]*Expense, error) {
	expenses := []*Expense{}
	for _, expense := range repo.expenses {
		if expense.GroupId == groupId {
			expenses = append(expenses, expense)
		}
	}
	return expenses, nil
}

func (repo *ExpenseRepo) GetExpenseById(expenseId string) (*Expense, error) {
	return repo.expenses[expenseId], nil
}

type IGroupRepo interface {
	AddGroup(group *Group) error
	GetGroupById(groupId string) (*Group, error)
}

type GroupRepo struct {
	groups map[string]*Group
}

func (repo *GroupRepo) AddGroup(group *Group) error {
	repo.groups[group.ID] = group
	return nil
}

func (repo *GroupRepo) GetGroupById(groupId string) (*Group, error) {
	return repo.groups[groupId], nil
}
