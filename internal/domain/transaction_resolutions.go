package domain

import (
	"airport-tools-backend/pkg/e"
	"time"
)

type Reason string

const (
	ModelError Reason = "MODEL_ERR"
	HumanError Reason = "HUMAN_ERR"
)

type TransactionResolution struct {
	Id            int64
	TransactionId int64
	QAEmployeeId  int64
	Reason        Reason
	Notes         string
	CreatedAt     time.Time
}

func NewTransactionResolution(transactionId int64, qaEmployeeId int64, reason Reason, notes string) *TransactionResolution {
	return &TransactionResolution{
		TransactionId: transactionId,
		QAEmployeeId:  qaEmployeeId,
		Notes:         notes,
		Reason:        reason,
	}
}

func ValidateReason(reason Reason) error {
	switch reason {
	case ModelError, HumanError:
		return nil
	}

	return e.ErrTransactionReasonInvalid
}
