package domain

// TransactionTool связывает транзакцию с инструментами.
// Хранит количество выданного инструмента (Qty).
type TransactionTool struct {
	ID            int64
	TransactionID int64 // fk Transaction
	ToolID        int64 // fk Tool
}

func NewTransactionTool(transactionID, toolID int64) *TransactionTool {
	return &TransactionTool{
		TransactionID: transactionID,
		ToolID:        toolID,
	}
}
