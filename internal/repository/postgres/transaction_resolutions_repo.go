package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type TransactionResolutionsRepo struct {
	DB *gorm.DB
}

func NewTransactionResolutionsRepo(db *gorm.DB) *TransactionResolutionsRepo {
	return &TransactionResolutionsRepo{
		DB: db,
	}
}

func (t *TransactionResolutionsRepo) Create(ctx context.Context, transaction *domain.TransactionResolution) (*domain.TransactionResolution, error) {
	const op = "TransactionResolutionsRepo.Create"

	model := toTransactionResolutionModel(transaction)

	result := t.DB.WithContext(ctx).Create(model)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransactionResolution(model), nil
}

func (t *TransactionResolutionsRepo) GetAll(ctx context.Context) ([]*domain.TransactionResolution, error) {
	const op = "TransactionResolutionsRepo.GetAll"

	var models []*TransactionResolutionModel
	result := t.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactionResolution(models), nil
}

func (t *TransactionResolutionsRepo) GetById(ctx context.Context, id int64) (*domain.TransactionResolution, error) {
	const op = "TransactionResolutionsRepo.GetById"

	var model TransactionResolutionModel
	result := t.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransactionResolution(&model), nil
}

func (t *TransactionResolutionsRepo) GetByQAId(ctx context.Context, qaEmployeeId string) ([]*domain.TransactionResolution, error) {
	const op = "TransactionResolutionsRepo.GetByQAId"

	var models []*TransactionResolutionModel
	result := t.DB.WithContext(ctx).Find(&models, "qa_employee_id = ?", qaEmployeeId)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactionResolution(models), nil
}

func toTransactionResolutionModel(transaction *domain.TransactionResolution) *TransactionResolutionModel {
	return &TransactionResolutionModel{
		Id:            transaction.Id,
		TransactionId: transaction.TransactionId,
		QAEmployeeId:  transaction.QAEmployeeId,
		Notes:         transaction.Notes,
		CreatedAt:     transaction.CreatedAt,
	}
}

func toDomainTransactionResolution(model *TransactionResolutionModel) *domain.TransactionResolution {
	return &domain.TransactionResolution{
		Id:            model.Id,
		TransactionId: model.TransactionId,
		QAEmployeeId:  model.QAEmployeeId,
		Notes:         model.Notes,
		CreatedAt:     model.CreatedAt,
	}
}

func toDomainArrTransactionResolution(models []*TransactionResolutionModel) []*domain.TransactionResolution {
	res := make([]*domain.TransactionResolution, len(models))
	for i, model := range models {
		res[i] = toDomainTransactionResolution(model)
	}

	return res
}
