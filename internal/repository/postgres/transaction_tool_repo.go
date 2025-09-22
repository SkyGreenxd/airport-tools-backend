package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type TransactionToolRepository struct {
	DB *gorm.DB
}

func NewTransactionToolRepository(db *gorm.DB) *TransactionToolRepository {
	return &TransactionToolRepository{
		DB: db,
	}
}

func (t *TransactionToolRepository) Create(ctx context.Context, transactionTool *domain.TransactionTool) (*domain.TransactionTool, error) {
	const op = "TransactionToolRepository.Create"

	var model TransactionToolModel
	result := t.DB.WithContext(ctx).Create(&model)
	if err := postgresDuplicate(result, e.ErrTransactionToolExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransactionTool(&model), nil
}

func (t *TransactionToolRepository) GetById(ctx context.Context, id int64) (*domain.TransactionTool, error) {
	const op = "TransactionToolRepository.GetById"

	var model TransactionToolModel
	result := t.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrTransactionToolNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransactionTool(&model), nil
}

func (t *TransactionToolRepository) GetByTransactionId(ctx context.Context, transactionId int64) ([]*domain.TransactionTool, error) {
	const op = "TransactionToolRepository.GetByTransactionId"

	var models []*TransactionToolModel
	result := t.DB.WithContext(ctx).Find(&models, "transaction_id = ?", transactionId)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrTransactionNotFound
	}

	return toDomainArrTransactionsTools(models), nil
}

func (t *TransactionToolRepository) GetByToolId(ctx context.Context, toolId int64) ([]*domain.TransactionTool, error) {
	const op = "TransactionToolRepository.GetByToolId"

	var models []*TransactionToolModel
	result := t.DB.WithContext(ctx).Find(&models, "tool_id = ?", toolId)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrTransactionNotFound
	}

	return toDomainArrTransactionsTools(models), nil
}

func (t *TransactionToolRepository) Delete(ctx context.Context, id int64) error {
	const op = "TransactionToolRepository.Delete"

	var model TransactionToolModel
	result := t.DB.WithContext(ctx).Delete(&model, "id = ?", id)
	if err := result.Error; err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

// TODO: перепроверить
func (t *TransactionToolRepository) GetUnreturnedByTransactionID(ctx context.Context, transactionId int64) ([]*domain.TransactionTool, error) {
	const op = "TransactionToolRepository.GetUnreturnedByTransactionID"

	var models []*TransactionToolModel
	result := t.DB.WithContext(ctx).
		Joins("JOIN transactions ON transactions.id = transaction_tools.transaction_id").
		Where("transaction_tools.transaction_id = ? AND transactions.returned_at IS NULL", transactionId).
		Preload("Tools").
		Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrTransactionNotFound
	}

	return toDomainArrTransactionsTools(models), nil
}

func toTransactionToolModel(t *domain.TransactionTool) *TransactionToolModel {
	return &TransactionToolModel{
		Id:            t.Id,
		TransactionId: t.TransactionId,
		ToolId:        t.ToolId,
		Transaction:   toTransactionModel(t.TransactionObj),
		Tool:          toToolModel(t.ToolObj),
	}
}

func toDomainTransactionTool(t *TransactionToolModel) *domain.TransactionTool {
	return &domain.TransactionTool{
		Id:             t.Id,
		TransactionId:  t.TransactionId,
		ToolId:         t.ToolId,
		TransactionObj: toDomainTransaction(t.Transaction),
		ToolObj:        toDomainTool(t.Tool),
	}
}

func toModelArrTransactionsTools(transactionsTools []*domain.TransactionTool) []*TransactionToolModel {
	models := make([]*TransactionToolModel, len(transactionsTools))
	for i, transactionTool := range transactionsTools {
		models[i] = toTransactionToolModel(transactionTool)
	}

	return models
}

func toDomainArrTransactionsTools(models []*TransactionToolModel) []*domain.TransactionTool {
	transactionsTools := make([]*domain.TransactionTool, len(models))
	for i, model := range models {
		transactionsTools[i] = toDomainTransactionTool(model)
	}

	return transactionsTools
}
