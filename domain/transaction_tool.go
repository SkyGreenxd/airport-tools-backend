package domain

// TransactionTool связывает транзакцию с инструментами.
// Хранит количество выданного инструмента (Qty).
// TODO: мб перенести сущность к Transaction
type TransactionTool struct {
	ID            int64
	TransactionID int64 // fk Transaction
	ToolID        int64 // fk Tool
	Qty           int   // сколько выдали
}
