package domain

// Transaction представляет запись о выдаче или возврате инструментов

type Status string
type Reason string

const (
	OPEN   Status = "OPEN"
	CLOSED Status = "CLOSED"
	MANUAL Status = "MANUAL VERIFICATION"

	RETURNED Reason = "All instruments have been handed over"
	PROBLEMS Reason = "There are problems with the tools, a manual check is needed"
)

type Transaction struct {
	Id        int64
	UserId    int64 // Received в UI, у кого инструмент
	ToolSetId int64
	Status    Status
	Reason    *Reason

	User    *User
	CvScans []*CvScan
}

func NewTransaction(userId, toolSetId int64) *Transaction {
	return &Transaction{
		UserId:    userId,
		Status:    OPEN,
		ToolSetId: toolSetId,
	}
}
