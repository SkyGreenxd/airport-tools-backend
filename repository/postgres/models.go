package postgres

import (
	"airport-tools-backend/domain"
	"time"
)

// TODO: добавить GORM теги
type StationModel struct {
	Id   int64
	Code string
}

type LocationModel struct {
	Id      int64
	StoreId int64
	Name    string
}

type StoreModel struct {
	Id        int64
	StationId int64
	Name      string
}

type ToolModel struct {
	Id         int64
	TypeToolId int64
	ToirId     int64
	LocationId int64
	SnBn       string
	ExpiresAt  time.Time
}

type ToolType struct {
	Id          int64
	PartNumber  string
	Name        string
	Description string
	Co          string
	MC          string
}

type TransactionModel struct {
	Id               int64
	UserId           int64
	Type             domain.TypeTransaction
	IssuedAt         time.Time
	ExpectedReturnAt time.Time
	ReturnedAt       *time.Time
	Tools            []domain.TransactionTool
}

type TransactionToolModel struct {
	ID            int64
	TransactionID int64
	ToolID        int64
	Qty           int
}

type UserModel struct {
	Id         int64
	EmployeeId string
	FullName   string
	Role       domain.Role
}
