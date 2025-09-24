package domain

import "airport-tools-backend/pkg/e"

type Role string

const (
	Engineer       Role = "Engineer"        // Авиатехник / Инженер
	QualityAuditor Role = "Quality Auditor" // Специалист службы качества / аудит
	SupplyManager  Role = "Supply Manager"  // Руководитель материально-технического снабжения
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

// Функция для проверки роли, если она пришла извне
func ValidateRole(r Role) bool {
	switch r {
	case Engineer, QualityAuditor, SupplyManager:
		return true
	}

	return false
}

func (u *User) ChangeFullName(newFullName string) error {
	if u.FullName == newFullName {
		return e.ErrNothingToChange
	}

	u.FullName = newFullName
	return nil
}

func (u *User) ChangeRole(r Role) error {
	if u.Role == r {
		return e.ErrNothingToChange
	}

	u.Role = r
	return nil
}
