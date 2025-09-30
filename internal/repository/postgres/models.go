package postgres

import (
	"airport-tools-backend/internal/domain"
	"time"

	"github.com/pgvector/pgvector-go"
)

type ToolTypeModel struct {
	Id                 int64
	PartNumber         string
	Name               string
	ReferenceEmbedding pgvector.Vector `gorm:"type:vector(1280)"`

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
	Reason    *domain.Reason
	CreatedAt time.Time
	UpdatedAt time.Time

	User    *UserModel     `gorm:"foreignkey:UserId"`
	CvScans []*CvScanModel `gorm:"foreignkey:TransactionId"`
}

type CvScanModel struct {
	Id            int64
	TransactionId int64
	ScanType      domain.ScanType
	ImageUrl      string
	CreatedAt     time.Time

	Transaction   *TransactionModel    `gorm:"foreignKey:TransactionId"`
	DetectedTools []*CvScanDetailModel `gorm:"foreignKey:CvScanId"`
}

type CvScanDetailModel struct {
	Id                 int64
	CvScanId           int64
	DetectedToolTypeId int64
	Confidence         float32
	Embedding          pgvector.Vector `gorm:"type:vector(1280)"`
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
