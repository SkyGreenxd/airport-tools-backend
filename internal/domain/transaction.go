package domain

import (
	"airport-tools-backend/pkg/e"
	"time"
)

type Status string

const (
	OPEN   Status = "OPEN"
	CLOSED Status = "CLOSED"
	QA     Status = "QA VERIFICATION"
)

type Transaction struct {
	Id            int64
	UserId        int64 // Received в UI, у кого инструмент
	ToolSetId     int64
	CountOfChecks int64
	Status        Status
	CreatedAt     time.Time
	UpdatedAt     time.Time

	User    *User
	CvScans []*CvScan
}

func NewTransaction(userId, toolSetId int64) *Transaction {
	return &Transaction{
		UserId:        userId,
		Status:        OPEN,
		ToolSetId:     toolSetId,
		CountOfChecks: 0,
	}
}

// TODO: подумать подходит ли ща под QA проверку
// мб для QA сделать отдельную функцию
func (t *Transaction) EvaluateStatus(manualCheckCount, unknownCount, missingCount int) {
	var status Status

	sum := 0
	sum += manualCheckCount
	sum += unknownCount
	sum += missingCount

	if sum == 0 {
		status = CLOSED
	} else if sum >= 4 {
		status = QA
	} else {
		status = OPEN
	}

	t.Status = status
}

func (t *Transaction) CheckCountOfChecks() error {
	if t.CountOfChecks >= 3 {
		return e.ErrTransactionLimit
	}

	return nil
}
