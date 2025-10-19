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
	FAILED Status = "FAILED"
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

func NewTransaction(userId, toolSetId int64, status Status) *Transaction {
	return &Transaction{
		UserId:        userId,
		Status:        status,
		ToolSetId:     toolSetId,
		CountOfChecks: 0,
	}
}

func (t *Transaction) EvaluateStatus(manualCheckCount, unknownCount, missingCount int) {
	var status Status

	sum := 0
	sum += manualCheckCount
	sum += unknownCount
	sum += missingCount

	if sum == 0 {
		status = CLOSED
	} else if sum >= 4 || t.CountOfChecks >= 3 {
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

func ValidateStatus(status string) (Status, error) {
	switch status {
	case string(OPEN):
		return OPEN, nil
	case string(CLOSED):
		return CLOSED, nil
	case string(QA), "QA":
		return QA, nil
	case string(FAILED):
		return FAILED, nil
	}

	return "", e.ErrTransactionStatusNotFound
}
