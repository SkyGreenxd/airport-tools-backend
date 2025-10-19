package repository

import (
	"airport-tools-backend/internal/domain"
	"context"
	"time"
)

// ToolTypeRepository интерфейс для работы с типами инструментов в базе данных
type ToolTypeRepository interface {
	Create(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error)
	GetById(ctx context.Context, id int64) (*domain.ToolType, error)
	GetAll(ctx context.Context) ([]*domain.ToolType, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error)
}

// ToolSetRepository интерфейс для работы с наборами инструментов в базе данных
type ToolSetRepository interface {
	Create(ctx context.Context, toolSet *domain.ToolSet) (*domain.ToolSet, error)
	GetById(ctx context.Context, id int64) (*domain.ToolSet, error)
	GetAll(ctx context.Context) ([]*domain.ToolSet, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, toolSet *domain.ToolSet) (*domain.ToolSet, error)
	GetByIdWithTools(ctx context.Context, id int64) (*domain.ToolSet, error)
	CreateWithTools(ctx context.Context, toolSet *domain.ToolSet, toolsIds []int64) (*domain.ToolSet, error)
}

// UserRepository интерфейс для работы с пользователями в базе данных
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetByEmployeeId(ctx context.Context, employeeId string) (*domain.User, error)
	GetByEmployeeIdWithTransactions(ctx context.Context, employeeId string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByEmployeeIdWithTransactionResolutions(ctx context.Context, employeeId string) (*domain.User, error)
	GetAllQa(ctx context.Context) ([]*domain.User, error)
	GetAllEngineersWithTransactions(ctx context.Context) ([]*domain.User, error)
}

// TransactionRepository интерфейс для работы с транзакциями инструментов в базе данных
type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetById(ctx context.Context, id int64) (*domain.Transaction, error)
	GetByUserIds(ctx context.Context, userIds []int64) ([]*domain.Transaction, error)
	GetByUserIdWhereStatusIsOpenOrQA(ctx context.Context, userId int64) (*domain.Transaction, error)
	GetByIdWithCvScans(ctx context.Context, id int64) (*domain.Transaction, error)
	GetByIdWithUser(ctx context.Context, id int64) (*domain.Transaction, error)
	GetAll(ctx context.Context) ([]*domain.Transaction, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetAllWithUser(ctx context.Context) ([]*domain.Transaction, error)
	GetAllWithStatusAndUser(ctx context.Context, status domain.Status) ([]*domain.Transaction, error)
	GetLastFailedByUserId(ctx context.Context, userId int64) (*domain.Transaction, error)
	GetAllByUserId(ctx context.Context, userId int64, startDate, endDate *time.Time, limit *int) ([]*domain.Transaction, error)
	GetAllWithStatus(ctx context.Context, status domain.Status) ([]*domain.Transaction, error)
}

// CvScanRepository интерфейс для работы со сканами инструментов в базе данных
type CvScanRepository interface {
	Create(ctx context.Context, cvScan *domain.CvScan) (*domain.CvScan, error)
	GetById(ctx context.Context, id int64) (*domain.CvScan, error)
	GetByTransactionId(ctx context.Context, transactionId int64) (*domain.CvScan, error)
	GetByIdWithTransaction(ctx context.Context, id int64) (*domain.CvScan, error)
	GetByTransactionIdWithDetectedToolsAndTransaction(ctx context.Context, transactionId int64) (*domain.CvScan, error)
}

// CvScanDetailRepository интерфейс для работы с детализацией сканов в базе данных
type CvScanDetailRepository interface {
	Create(ctx context.Context, cvScanDetail *domain.CvScanDetail) (*domain.CvScanDetail, error)
	GetById(ctx context.Context, id int64) (*domain.CvScanDetail, error)
	GetByCvScanId(ctx context.Context, cvScanId int64) ([]*domain.CvScanDetail, error)
}

// ImageRepository интерфейс для работы с хранением изображений
type ImageRepository interface {
	Save(ctx context.Context, img *domain.Image) (*domain.UploadImage, error)
}

// TransactionResolutionsRepository интерфейс для работы с QA проверками
type TransactionResolutionsRepository interface {
	Create(ctx context.Context, transaction *domain.TransactionResolution, toolIds []int64) (*domain.TransactionResolution, error)
	GetAll(ctx context.Context) ([]*domain.TransactionResolution, error)
	GetById(ctx context.Context, id int64) (*domain.TransactionResolution, error)
	GetByQAId(ctx context.Context, qaId int64) ([]*domain.TransactionResolution, error)
	GetAllModelError(ctx context.Context) ([]*domain.TransactionResolution, error)
	GetAllHumanError(ctx context.Context) ([]*domain.TransactionResolution, error)
	GetTopHumanErrorUsers(ctx context.Context) ([]HumanErrorStats, error)
	GetMlErrorTransactions(ctx context.Context) ([]*domain.TransactionResolution, error)
	GetMlErrorTools(ctx context.Context) ([]*ToolSetWithErrors, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role *domain.Role) (*domain.Role, error)
	GetAll(ctx context.Context) ([]*domain.Role, error)
	GetById(ctx context.Context, id int64) (*domain.Role, error)
	GetByName(ctx context.Context, name string) (*domain.Role, error)
}
