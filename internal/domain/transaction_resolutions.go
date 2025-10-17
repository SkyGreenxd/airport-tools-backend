package domain

import "time"

type TransactionResolution struct {
	Id            int64
	TransactionId int64
	QAEmployeeId  int64
	Notes         string
	CreatedAt     time.Time
}

func NewTransactionResolution(transactionId int64, qaEmployeeId int64, notes string) *TransactionResolution {
	return &TransactionResolution{
		TransactionId: transactionId,
		QAEmployeeId:  qaEmployeeId,
		Notes:         notes,
	}
}
