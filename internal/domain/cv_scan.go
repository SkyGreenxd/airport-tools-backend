package domain

type CvScan struct {
	Id            int64
	TransactionId int64
	Status        string
	Reason        string
	ImageUrl      string
	ImageHash     string
}

func NewCvScan(transactionId int64, status, reason, imageUrl, imageHash string) *CvScan {
	return &CvScan{
		TransactionId: transactionId,
		Status:        status,
		Reason:        reason,
		ImageUrl:      imageUrl,
		ImageHash:     imageHash,
	}
}
