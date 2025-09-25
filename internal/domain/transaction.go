package domain

// Transaction представляет запись о выдаче или возврате инструментов

type Status string

const (
	OPEN   Status = "OPEN"
	CLOSED Status = "CLOSED"
	MANUAL Status = "MANUAL VERIFICATION"
)

type Transaction struct {
	Id        int64
	UserId    int64 // Received в UI, у кого инструмент
	ToolSetId int64
	Status    Status
	Reason    *string

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
