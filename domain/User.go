package domain

type Role string

const (
	Engineer       Role = "Engineer"        // Авиатехник / Инженер
	QualityAuditor Role = "Quality Auditor" // Специалист службы качества / аудит
	SupplyManager  Role = "Supply Manager"  // Руководитель материально-технического снабжения
)

// User описывает пользователя
type User struct {
	Id         int64
	EmployeeId string // табельный номер
	FullName   string // фио
	Role       Role   // роль пользователя
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
