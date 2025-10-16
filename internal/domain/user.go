package domain

import "airport-tools-backend/pkg/e"

type Role string

const (
	Engineer       Role = "Engineer"        // Авиатехник / Инженер
	QualityAuditor Role = "Quality Auditor" // Специалист службы качества / аудит
	// SupplyManager  Role = "Supply Manager"  // Руководитель материально-технического снабжения
)

type User struct {
	Id         int64
	EmployeeId string
	FullName   string
	Role       Role

	Transactions []*Transaction
}

func NewUser(fullName, employeeId string, role Role) *User {
	return &User{
		FullName:   fullName,
		EmployeeId: employeeId,
		Role:       role,
	}
}

func (u *User) CanCheckout() error {
	if len(u.Transactions) == 0 {
		return nil
	}

	for _, transaction := range u.Transactions {
		if transaction.Status == OPEN {
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

func ValidateRole(role Role) error {
	switch role {
	case Engineer, QualityAuditor:
		return nil
	default:
		return e.ErrUserRoleNotFound
	}
}

func (u *User) ChangeRole(newRole Role) error {
	if u.Role == newRole {
		return e.ErrNothingToChange
	}

	u.Role = newRole
	return nil
}
