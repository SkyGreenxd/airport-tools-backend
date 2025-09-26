package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

func (t *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	const op = "TransactionRepository.Create"

	model := toTransactionModel(transaction)
	result := t.DB.WithContext(ctx).Create(model)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransaction(model), nil
}

func (t *TransactionRepository) GetById(ctx context.Context, id int64) (*domain.Transaction, error) {
	const op = "TransactionRepository.GetById"

	var model TransactionModel
	result := t.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrTransactionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransaction(&model), nil
}

func (t *TransactionRepository) GetByUserId(ctx context.Context, userId int64) (*domain.Transaction, error) {
	const op = "TransactionRepository.GetByEmployeeId"

	var model TransactionModel
	result := t.DB.WithContext(ctx).First(&model, "user_id = ?", userId)
	if err := checkGetQueryResult(result, e.ErrTransactionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransaction(&model), nil
}

func (t *TransactionRepository) GetByUserIdWhereStatusIsOpenOrManual(ctx context.Context, userId int64) (*domain.Transaction, error) {
	const op = "TransactionRepository.GetByUserIdWhereStatusIsOpenOrManual"
	var model TransactionModel
	result := t.DB.WithContext(ctx).
		Where("user_id = ? AND status IN ?", userId, []domain.Status{domain.OPEN, domain.MANUAL}).
		First(&model)

	if err := checkGetQueryResult(result, e.ErrTransactionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransaction(&model), nil
}

func (t *TransactionRepository) GetByIdWithCvScans(ctx context.Context, id int64) (*domain.Transaction, error) {
	const op = "TransactionRepository.GetByIdWithCvScans"

	var model TransactionModel
	result := t.DB.WithContext(ctx).Preload("CvScans").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrTransactionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransaction(&model), nil
}

func (t *TransactionRepository) GetByIdWithUser(ctx context.Context, id int64) (*domain.Transaction, error) {
	const op = "TransactionRepository.GetByIdWithUser"

	var model TransactionModel
	result := t.DB.WithContext(ctx).Preload("User").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrTransactionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransaction(&model), nil
}

func (t *TransactionRepository) GetAll(ctx context.Context) ([]*domain.Transaction, error) {
	const op = "TransactionRepository.GetAll"

	var models []*TransactionModel
	result := t.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactions(models), nil
}

func (t *TransactionRepository) Delete(ctx context.Context, id int64) error {
	const op = "TransactionRepository.Delete"

	var model TransactionModel
	result := t.DB.WithContext(ctx).Delete(&model, "id = ?", id)
	if result.Error != nil {
		return e.Wrap(op, result.Error)
	}

	return nil
}

func (t *TransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	const op = "TransactionRepository.Update"

	updates := map[string]interface{}{
		"status": transaction.Status,
		"reason": transaction.Reason,
	}

	result := t.DB.WithContext(ctx).Model(&TransactionModel{}).Where("id = ?", transaction.Id).Updates(updates)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrTransactionNotFound)
	}

	updTransaction, err := t.GetById(ctx, transaction.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updTransaction, nil
}

func toTransactionModel(t *domain.Transaction) *TransactionModel {
	model := &TransactionModel{
		Id:        t.Id,
		UserId:    t.UserId,
		ToolSetId: t.ToolSetId,
		Status:    t.Status,
		Reason:    t.Reason,
	}

	if t.CvScans != nil {
		model.CvScans = toArrCvScansModel(t.CvScans)
	}

	if t.User != nil {
		model.User = toUserModel(t.User)
	}

	return model
}

func toDomainTransaction(t *TransactionModel) *domain.Transaction {
	transaction := &domain.Transaction{
		Id:        t.Id,
		UserId:    t.UserId,
		ToolSetId: t.ToolSetId,
		Status:    t.Status,
		Reason:    t.Reason,
	}

	if t.CvScans != nil {
		transaction.CvScans = toArrDomainCvScans(t.CvScans)
	}

	if t.User != nil {
		transaction.User = toDomainUser(t.User)
	}

	return transaction
}

func toModelArrTransactions(transactions []*domain.Transaction) []*TransactionModel {
	models := make([]*TransactionModel, len(transactions))
	for i, transaction := range transactions {
		models[i] = toTransactionModel(transaction)
	}

	return models
}

func toDomainArrTransactions(models []*TransactionModel) []*domain.Transaction {
	transactions := make([]*domain.Transaction, len(models))
	for i, model := range models {
		transactions[i] = toDomainTransaction(model)
	}

	return transactions
}
