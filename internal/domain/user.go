package domain

import "airport-tools-backend/pkg/e"

type Role string

const (
	Engineer       Role = "Engineer"        // Авиатехник / Инженер
	QualityAuditor Role = "Quality Auditor" // Специалист службы качества / аудит
	SupplyManager  Role = "Supply Manager"  // Руководитель материально-технического снабжения
)

// User описывает пользователя
// TODO: мб стоит привязать пользователя к аэропорту
type User struct {
	Id         int64
	EmployeeId string // табельный номер
	FullName   string // фио
	Role       Role   // роль пользователя

	Transactions []*Transaction
}

func NewUser(employeeId, fullName, role string) *User {
	return &User{
		EmployeeId: employeeId,
		FullName:   fullName,
		Role:       Role(role),
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
