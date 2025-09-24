package domain

// Transaction представляет запись о выдаче или возврате инструментов
type Transaction struct {
	Id     int64
	UserId int64 // Received в UI, у кого инструмент
	Status string
	Reason *string

	User    *User
	CvScans []*CvScan
}

func NewTransaction(userId int64, status string, reason *string) *Transaction {
	return &Transaction{
		UserId: userId,
		Status: status,
		Reason: reason,
	}
}
