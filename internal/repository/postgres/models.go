package postgres

import (
	"airport-tools-backend/internal/domain"
)

type ToolTypeModel struct {
	Id                 int64
	PartNumber         string
	Name               string
	ReferenceImageHash string
	ReferenceEmbedding []float32
}

type ToolSetModel struct {
	Id   int64
	Name string

	Tools []*ToolTypeModel
}
type ToolSetItemModel struct {
	Id         int64
	ToolSetId  int64
	ToolTypeId int64
}

type UserModel struct {
	Id               int64
	EmployeeId       int64
	FullName         string
	Role             domain.Role
	DefaultToolSetId int64

	Transactions []*TransactionModel
}

type TransactionModel struct {
	Id     int64
	UserId int64
	Status string
	Reason string

	User    *UserModel
	CvScans []*CvScanModel
}

type CvScanModel struct {
	Id            int64
	TransactionId int64
	ScanType      string
	ImageUrl      string

	Transaction   *TransactionModel
	DetectedTools []*CvScanDetailModel
}

type CvScanDetailModel struct {
	Id                 int64
	CvScanId           int64
	DetectedToolTypeId int64
	ImageHash          string
	Embedding          []float32
}

func (ToolTypeModel) TableName() string {
	return "tool_types"
}

func (ToolSetModel) TableName() string {
	return "tool_sets"
}

func (ToolSetItemModel) TableName() string {
	return "tool_set_items"
}

func (UserModel) TableName() string {
	return "users"
}

func (TransactionModel) TableName() string {
	return "transactions"
}

func (CvScanModel) TableName() string {
	return "cv_scans"
}

func (CvScanDetailModel) TableName() string {
	return "cv_scan_details"
}
