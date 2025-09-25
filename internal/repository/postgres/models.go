package postgres

import (
	"airport-tools-backend/internal/domain"

	"github.com/pgvector/pgvector-go"
)

type ToolTypeModel struct {
	Id                 int64
	PartNumber         string
	Name               string
	ReferenceImageHash string
	// TODO: заменить на VECTOR(1280)
	ReferenceEmbedding pgvector.Vector `gorm:"type:vector(3)"`

	ToolSets []*ToolSetModel `gorm:"many2many:tool_set_items;joinForeignKey:ToolTypeId;joinReferences:ToolSetId"`
}

type ToolSetModel struct {
	Id   int64
	Name string

	Tools []*ToolTypeModel `gorm:"many2many:tool_set_items;joinForeignKey:ToolSetId;joinReferences:ToolTypeId"`
}
type ToolSetItemModel struct {
	ToolSetId  int64 `gorm:"column:tool_set_id"`
	ToolTypeId int64 `gorm:"column:tool_type_id"`
}

type UserModel struct {
	Id               int64
	EmployeeId       string
	FullName         string
	Role             domain.Role
	DefaultToolSetId int64 `gorm:"column:default_tool_set_id"`

	Transactions []*TransactionModel `gorm:"foreignkey:UserId"`
}

type TransactionModel struct {
	Id        int64
	UserId    int64
	ToolSetId int64
	Status    domain.Status
	Reason    *string

	User    *UserModel     `gorm:"foreignkey:UserId"`
	CvScans []*CvScanModel `gorm:"foreignkey:TransactionId"`
}

type CvScanModel struct {
	Id            int64
	TransactionId int64
	ScanType      domain.ScanType
	ImageUrl      string

	Transaction   *TransactionModel    `gorm:"foreignKey:TransactionId"`
	DetectedTools []*CvScanDetailModel `gorm:"foreignKey:CvScanId"`
}

type CvScanDetailModel struct {
	Id                 int64
	CvScanId           int64
	DetectedToolTypeId int64
	Confidence         float32
	ImageHash          string
	// TODO: заменить на VECTOR(1280)
	Embedding pgvector.Vector `gorm:"type:vector(3)"`
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
