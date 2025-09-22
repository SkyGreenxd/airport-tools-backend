package domain

// TransactionTool связывает транзакцию с инструментами.
type TransactionTool struct {
	Id            int64
	TransactionId int64 // fk Transaction
	ToolId        int64 // fk Tool

	TransactionObj *Transaction
	ToolObj        *Tool
}

func NewTransactionTool(transactionId, toolId int64) *TransactionTool {
	return &TransactionTool{
		TransactionId: transactionId,
		ToolId:        toolId,
	}
}
