package domain

import "airport-tools-backend/pkg/e"

type User struct {
	Id         int64
	EmployeeId string
	FullName   string
	RoleId     int64

	Role                   *Role
	Transactions           []*Transaction
	TransactionResolutions []*TransactionResolution
}

func NewUser(fullName, employeeId string, roleId int64) *User {
	return &User{
		FullName:   fullName,
		EmployeeId: employeeId,
		RoleId:     roleId,
	}
}

func (u *User) CanCheckout() error {
	if len(u.Transactions) == 0 {
		return nil
	}

	for _, transaction := range u.Transactions {
		if transaction.Status == OPEN || transaction.Status == QA {
			return e.ErrTransactionUnfinished
		}
	}

	return nil
}

func (u *User) CanCheckin() error {
	if len(u.Transactions) == 0 {
		return e.ErrTransactionAllFinished
	}

	for _, transaction := range u.Transactions {
		if transaction.Status == OPEN {
			return nil
		} else if transaction.Status == QA {
			return e.ErrTransactionCheckQA
		}
	}

	return e.ErrTransactionAllFinished
}

func (u *User) ValidateEmployeeId(newEmployeeId string) error {
	if u.EmployeeId == newEmployeeId {
		return e.ErrNothingToChange
	}

	u.EmployeeId = newEmployeeId
	return nil
}

func (u *User) ValidateFullName(newFullName string) error {
	if u.FullName == newFullName {
		return e.ErrNothingToChange
	}

	u.FullName = newFullName
	return nil
}
