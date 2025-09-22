package postgres

import (
	"airport-tools-backend/internal/domain"
	"time"
)

type StationModel struct {
	Id   int64
	Code string

	Stores []*StoreModel `gorm:"foreignKey:StationId"`
}

type StoreModel struct {
	Id        int64
	StationId int64
	Name      string

	Station   *StationModel    `gorm:"foreignKey:StationId"`
	Locations []*LocationModel `gorm:"foreignKey:StoreId"`
}

type LocationModel struct {
	Id      int64
	StoreId int64
	Name    string

	Store *StoreModel  `gorm:"foreignKey:StoreId"`
	Tools []*ToolModel `gorm:"foreignKey:LocationId"`
}

type ToolModel struct {
	Id         int64
	TypeToolId int64
	ToirId     int64
	LocationId int64
	ExpiresAt  *time.Time
	// SnBn       string

	ToolType *ToolTypeModel `gorm:"foreignKey:TypeToolId"`
	Location *LocationModel `gorm:"foreignKey:LocationId"`
}

type ToolTypeModel struct {
	Id          int64
	PartNumber  string
	Description string
	//Co          string
	//MC          string

	Tools []*ToolModel `gorm:"foreignKey:TypeToolId"`
}

type TransactionModel struct {
	Id               int64
	UserId           int64
	Type             domain.TypeTransaction
	IssuedAt         time.Time
	ExpectedReturnAt time.Time
	ReturnedAt       *time.Time

	User  *UserModel              `gorm:"foreignKey:Userf"`
	Tools []*TransactionToolModel `gorm:"foreignKey:TransactionId"`
}

type TransactionToolModel struct {
	Id            int64
	TransactionId int64
	ToolId        int64

	Transaction *TransactionModel `gorm:"foreignKey:TransactionId"`
	Tool        *ToolModel        `gorm:"foreignKey:ToolId"`
}

type UserModel struct {
	Id         int64
	EmployeeId string
	FullName   string
	Role       domain.Role

	Transactions []*TransactionModel `gorm:"foreignKey:UserId"`
}

func (UserModel) TableName() string {
	return "users"
}

func (TransactionToolModel) TableName() string {
	return "transactions_tools"
}

func (TransactionModel) TableName() string {
	return "transactions"
}

func (ToolTypeModel) TableName() string {
	return "tool_types"
}

func (ToolModel) TableName() string {
	return "tools"
}

func (StoreModel) TableName() string {
	return "stores"
}

func (LocationModel) TableName() string {
	return "locations"
}

func (StationModel) TableName() string {
	return "stations"
}
