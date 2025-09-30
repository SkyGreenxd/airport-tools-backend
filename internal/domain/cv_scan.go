package domain

import "time"

type ScanType string

const (
	Checkin  ScanType = "checkin"  // сдача инструментов
	Checkout ScanType = "checkout" // выдача инструментов
)

type CvScan struct {
	Id            int64
	TransactionId int64
	ScanType      ScanType
	ImageUrl      string
	CreatedAt     time.Time

	TransactionObj *Transaction
	DetectedTools  []*CvScanDetail
}

func NewCvScan(transactionId int64, scanType ScanType, imageUrl string) *CvScan {
	return &CvScan{
		TransactionId: transactionId,
		ScanType:      scanType,
		ImageUrl:      imageUrl,
	}
}
