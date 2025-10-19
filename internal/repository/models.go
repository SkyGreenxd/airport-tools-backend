package repository

// HumanErrorStats — структура для хранения статистики ошибок, допущенных конкретным сотрудником (не Ml моделью)
type HumanErrorStats struct {
	FullName    string
	EmployeeId  string
	QAHitsCount int64
}

type ToolWithErrorCount struct {
	ID           int64
	Name         string
	MLErrorCount int64
}

type ToolSetWithErrors struct {
	ID    int64
	Name  string
	Tools []ToolWithErrorCount
}
