package repository

import (
	"airport-tools-backend/internal/domain"
	"context"
)

type ToolTypeRepository interface {
	Create(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error)
	GetById(ctx context.Context, id int64) (*domain.ToolType, error)
	GetAll(ctx context.Context) ([]*domain.ToolType, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error)
}

type ToolSetRepository interface {
	Create(ctx context.Context, toolSet *domain.ToolSet) (*domain.ToolSet, error)
	GetById(ctx context.Context, id int64) (*domain.ToolSet, error)
	GetAll(ctx context.Context) ([]*domain.ToolSet, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, toolSet *domain.ToolSet) (*domain.ToolSet, error)
	GetByIdWithTools(ctx context.Context, id int64) (*domain.ToolSet, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetByEmployeeId(ctx context.Context, employeeId string) (*domain.User, error)
	GetByIdWithTransactions(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetById(ctx context.Context, id int64) (*domain.Transaction, error)
	GetByUserId(ctx context.Context, userId int64) (*domain.Transaction, error)
	GetByUserIdWhereStatusIsOpenOrManual(ctx context.Context, userId int64) (*domain.Transaction, error)
	GetByIdWithCvScans(ctx context.Context, id int64) (*domain.Transaction, error)
	GetByIdWithUser(ctx context.Context, id int64) (*domain.Transaction, error)
	GetAll(ctx context.Context) ([]*domain.Transaction, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
}

type CvScanRepository interface {
	Create(ctx context.Context, cvScan *domain.CvScan) (*domain.CvScan, error)
	GetById(ctx context.Context, id int64) (*domain.CvScan, error)
	GetByTransactionId(ctx context.Context, transactionId int64) (*domain.CvScan, error)
	GetByIdWithTransaction(ctx context.Context, id int64) (*domain.CvScan, error)
	GetByIdWithDetectedTools(ctx context.Context, id int64) (*domain.CvScan, error)
}

type CvScanDetailRepository interface {
	Create(ctx context.Context, cvScanDetail *domain.CvScanDetail) (*domain.CvScanDetail, error)
	GetById(ctx context.Context, id int64) (*domain.CvScanDetail, error)
	GetByCvScanId(ctx context.Context, cvScanId int64) ([]*domain.CvScanDetail, error)
}
