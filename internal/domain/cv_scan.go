package domain

type CvScan struct {
	Id            int64
	TransactionId int64
	ScanType      string
	ImageUrl      string

	TransactionObj *Transaction
	DetectedTools  []*CvScanDetail
}

func NewCvScan(transactionId int64, scanType, imageUrl string) *CvScan {
	return &CvScan{
		TransactionId: transactionId,
		ScanType:      scanType,
		ImageUrl:      imageUrl,
	}
}
