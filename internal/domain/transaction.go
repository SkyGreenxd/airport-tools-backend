package domain

import "time"

// Transaction представляет запись о выдаче или возврате инструментов
type Transaction struct {
	Id         int64
	UserId     int64      // Received в UI, у кого инструмент
	IssuedAt   time.Time  // дата выдачи
	ReturnedAt *time.Time // фактическая дата возврата, может быть nil

	UserObj *User
	Tools   []*TransactionTool // вложенный список инструментов

}

func NewTransaction(userId int64) *Transaction {
	return &Transaction{
		UserId:     userId,
		ReturnedAt: nil,
	}
}

// MarkReturned устанавливает время возврата инструмента
func (t *Transaction) MarkReturned(returnTime time.Time) {
	t.ReturnedAt = &returnTime
}
