package repository

// HumanErrorStats — структура для хранения статистики ошибок, допущенных конкретным сотрудником (не Ml моделью)
type HumanErrorStats struct {
	FullName    string
	EmployeeId  string
	QAHitsCount int64
}
