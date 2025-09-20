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
}

type StoreRepository interface {
	Create(ctx context.Context, station *domain.Store) (*domain.Store, error)
	GetById(ctx context.Context, id int64) (*domain.Store, error)
	GetAll(ctx context.Context) ([]*domain.Store, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, station *domain.Store) (*domain.Store, error)
}

type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) (*domain.Location, error)
	GetById(ctx context.Context, id int64) (*domain.Location, error)
	GetAll(ctx context.Context) ([]*domain.Location, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, location *domain.Location) (*domain.Location, error)
}

type ToolRepository interface {
	Create(ctx context.Context, tool *domain.Tool) (*domain.Tool, error)
	GetById(ctx context.Context, id int64) (*domain.Tool, error)
	GetAll(ctx context.Context) ([]*domain.Tool, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, tool *domain.Tool) (*domain.Tool, error)
}

type ToolTypeRepository interface {
	Create(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error)
	GetById(ctx context.Context, id int64) (*domain.ToolType, error)
	GetAll(ctx context.Context) ([]*domain.ToolType, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, station *domain.Transaction) (*domain.Transaction, error)
	GetById(ctx context.Context, id int64) (*domain.Transaction, error)
	GetAll(ctx context.Context) ([]*domain.Transaction, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, station *domain.Transaction) (*domain.Transaction, error)
}

type TransactionToolRepository interface {
	Create(ctx context.Context, station *domain.TransactionTool) (*domain.TransactionTool, error)
	GetById(ctx context.Context, id int64) (*domain.TransactionTool, error)
	GetByTransactionId(ctx context.Context, transactionId int64) ([]*domain.TransactionTool, error)
	Update(ctx context.Context, station *domain.TransactionTool) (*domain.TransactionTool, error)
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]*domain.TransactionTool, error)
	GetUnreturnedByTransactionID(ctx context.Context, transactionId int64) ([]*domain.TransactionTool, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
}
