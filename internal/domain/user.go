package domain

import "airport-tools-backend/pkg/e"

type Role string

const (
	Engineer Role = "Engineer" // Авиатехник / Инженер
	// QualityAuditor Role = "Quality Auditor" // Специалист службы качества / аудит
	// SupplyManager  Role = "Supply Manager"  // Руководитель материально-технического снабжения
)

type User struct {
	Id               int64
	EmployeeId       string
	FullName         string
	Role             Role
	DefaultToolSetId int64

	Transactions []*Transaction
}

func NewUser(fullName, employeeId string, role Role, toolSetId int64) *User {
	return &User{
		FullName:         fullName,
		EmployeeId:       employeeId,
		Role:             role,
		DefaultToolSetId: toolSetId,
	}
}

func (u *User) CanCheckout() error {
	if len(u.Transactions) == 0 {
		return nil
	}

	for _, transaction := range u.Transactions {
		if transaction.Status == OPEN || transaction.Status == MANUAL {
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
		if transaction.Status == OPEN || transaction.Status == MANUAL {
			return nil
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

func (u *User) ValidateRole(newRole Role) error {
	if u.Role == newRole {
		return e.ErrNothingToChange
	}

	u.Role = newRole
	return nil
}

func (u *User) ValidateDefaultToolSetId(newToolSetId int64) error {
	if u.DefaultToolSetId == newToolSetId {
		return e.ErrNothingToChange
	}

	u.DefaultToolSetId = newToolSetId
	return nil
}
