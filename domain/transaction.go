package domain

import "time"

type TypeTransaction string

const (
	Issuance TypeTransaction = "Issuance"
	Return   TypeTransaction = "Return"
)

// Transaction представляет запись о выдаче или возврате инструментов
type Transaction struct {
	Id               int64
	UserId           int64             // Received в UI, у кого инструмент
	Type             TypeTransaction   // тип транзакции
	IssuedAt         time.Time         // дата выдачи
	ExpectedReturnAt time.Time         // ожидаемая дата возврата
	ReturnedAt       *time.Time        // фактическая дата возврата, может быть nil
	Tools            []TransactionTool // вложенный список инструментов
}

// TODO: если по БЛ известно, когда возвращают инструменты
// то логичнее заполнять автоматически через код/бд
func NewTransaction(userId int64, typeTransaction TypeTransaction, ExpectedReturnAt time.Time) *Transaction {
	return &Transaction{
		UserId:           userId,
		Type:             typeTransaction,
		ExpectedReturnAt: ExpectedReturnAt,
		ReturnedAt:       nil,
	}
}

func (t *Transaction) AddTool(toolId int64, qty int) {
	tool := TransactionTool{
		TransactionID: t.Id,
		ToolID:        toolId,
		Qty:           qty,
	}

	t.Tools = append(t.Tools, tool)
}

// MarkReturned устанавливает время возврата инструмента
func (t *Transaction) MarkReturned(returnTime time.Time) {
	t.ReturnedAt = &returnTime
}
