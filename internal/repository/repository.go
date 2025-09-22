package repository

import (
	"airport-tools-backend/internal/domain"
	"context"
)

type StationRepository interface {
	Create(ctx context.Context, station *domain.Station) (*domain.Station, error)
	GetById(ctx context.Context, id int64) (*domain.Station, error)
	GetAll(ctx context.Context) ([]*domain.Station, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, station *domain.Station) (*domain.Station, error)
	GetByIdWithStores(ctx context.Context, id int64) (*domain.Station, error)
}

type StoreRepository interface {
	Create(ctx context.Context, store *domain.Store) (*domain.Store, error)
	GetById(ctx context.Context, id int64) (*domain.Store, error)
	GetAll(ctx context.Context) ([]*domain.Store, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, store *domain.Store) (*domain.Store, error)
	GetByIdWithStation(ctx context.Context, id int64) (*domain.Store, error)
	GetByIdWithLocations(ctx context.Context, id int64) (*domain.Store, error)
}

type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) (*domain.Location, error)
	GetById(ctx context.Context, id int64) (*domain.Location, error)
	GetAll(ctx context.Context) ([]*domain.Location, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, location *domain.Location) (*domain.Location, error)
	GetByIdWithTools(ctx context.Context, id int64) (*domain.Location, error)
}

type ToolRepository interface {
	Create(ctx context.Context, tool *domain.Tool) (*domain.Tool, error)
	GetById(ctx context.Context, id int64) (*domain.Tool, error)
	GetAll(ctx context.Context) ([]*domain.Tool, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, tool *domain.Tool) (*domain.Tool, error)
	GetByIdWithToolType(ctx context.Context, toolId int64) (*domain.Tool, error)
	GetByIdWithLocation(ctx context.Context, toolId int64) (*domain.Tool, error)
	GetByIdWithAllData(ctx context.Context, toolId int64) (*domain.Tool, error)
}

type ToolTypeRepository interface {
	Create(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error)
	GetById(ctx context.Context, id int64) (*domain.ToolType, error)
	GetAll(ctx context.Context) ([]*domain.ToolType, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error)
	GetByIdWithTools(ctx context.Context, id int64) (*domain.ToolType, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetById(ctx context.Context, id int64) (*domain.Transaction, error)
	GetByIdWithTools(ctx context.Context, id int64) (*domain.Transaction, error)
	GetByIdWithUser(ctx context.Context, id int64) (*domain.Transaction, error)
	GetAll(ctx context.Context) ([]*domain.Transaction, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
}

type TransactionToolRepository interface {
	Create(ctx context.Context, transactionTool *domain.TransactionTool) (*domain.TransactionTool, error)
	GetById(ctx context.Context, id int64) (*domain.TransactionTool, error)
	GetByTransactionId(ctx context.Context, transactionId int64) ([]*domain.TransactionTool, error)
	GetByToolId(ctx context.Context, toolId int64) ([]*domain.TransactionTool, error)
	Delete(ctx context.Context, id int64) error
	GetUnreturnedByTransactionID(ctx context.Context, transactionId int64) ([]*domain.TransactionTool, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetByEmployeeId(ctx context.Context, employeeId int64) (*domain.User, error)
	GetByFullName(ctx context.Context, fullName string) (*domain.User, error)
	GetByIdWithTransactions(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
}
