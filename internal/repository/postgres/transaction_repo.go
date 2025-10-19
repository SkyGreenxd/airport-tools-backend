package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"
	"time"

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

func (t *TransactionRepository) GetByUserIds(ctx context.Context, userIds []int64) ([]*domain.Transaction, error) {
	const op = "TransactionRepository.GetByUserIds"

	var models []*TransactionModel
	result := t.DB.WithContext(ctx).
		Preload("User").
		Where("user_id IN ?", userIds).
		Find(&models)

	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactions(models), nil
}

func (t *TransactionRepository) GetLastFailedByUserId(ctx context.Context, userId int64) (*domain.Transaction, error) {
	const op = "TransactionRepository.GetLastByUserId"

	var model TransactionModel
	result := t.DB.WithContext(ctx).Where("user_id = ? AND status = ?", userId, domain.FAILED).First(&model)
	if err := checkGetQueryResult(result, e.ErrTransactionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTransaction(&model), nil
}

func (t *TransactionRepository) GetByUserIdWhereStatusIsOpenOrQA(ctx context.Context, userId int64) (*domain.Transaction, error) {
	const op = "TransactionRepository.GetByUserIdWhereStatusIsOpenOrQA"
	var model TransactionModel
	result := t.DB.WithContext(ctx).
		Where("user_id = ? AND status IN ?", userId, []domain.Status{domain.OPEN, domain.QA}).
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

func (t *TransactionRepository) GetAllWithUser(ctx context.Context) ([]*domain.Transaction, error) {
	const op = "TransactionRepository.GetAllWithUser"

	var models []*TransactionModel
	result := t.DB.WithContext(ctx).Preload("User").Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactions(models), nil
}

func (t *TransactionRepository) GetAllWhereStatusIsQAWithUser(ctx context.Context) ([]*domain.Transaction, error) {
	const op = "TransactionRepository.GetAllWhereStatusIsQA"

	var models []*TransactionModel
	result := t.DB.WithContext(ctx).Preload("User").Where("status = ?", domain.QA).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactions(models), nil
}

func (t *TransactionRepository) GetAllWithStatus(ctx context.Context, status domain.Status) ([]*domain.Transaction, error) {
	const op = "TransactionRepository.getAllWithStatus"

	var models []*TransactionModel
	result := t.DB.WithContext(ctx).Where("status = ?", status).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactions(models), nil
}

func (t *TransactionRepository) GetAllByUserId(ctx context.Context, userId int64, startDate, endDate *time.Time, limit *int) ([]*domain.Transaction, error) {
	const op = "TransactionRepository.GetAllByUserId"

	var models []*TransactionModel
	db := t.DB.WithContext(ctx).Preload("User").Where("user_id = ?", userId)

	if startDate != nil {
		db = db.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		db = db.Where("created_at <= ?", *endDate)
	}

	db = db.Order("id DESC")

	if limit != nil {
		db = db.Limit(*limit)
	}

	result := db.Find(&models)
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
		"status":          transaction.Status,
		"updated_at":      time.Now().UTC(),
		"count_of_checks": transaction.CountOfChecks,
	}

	var updTransaction TransactionModel
	result := t.DB.WithContext(ctx).Model(&TransactionModel{}).Where("id = ?", transaction.Id).Updates(updates).Scan(&updTransaction)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrTransactionNotFound)
	}

	return toDomainTransaction(&updTransaction), nil
}

func toTransactionModel(t *domain.Transaction) *TransactionModel {
	model := &TransactionModel{
		Id:            t.Id,
		UserId:        t.UserId,
		ToolSetId:     t.ToolSetId,
		CountOfChecks: t.CountOfChecks,
		Status:        t.Status,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
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
	if t == nil {
		return nil
	}

	transaction := &domain.Transaction{
		Id:            t.Id,
		UserId:        t.UserId,
		ToolSetId:     t.ToolSetId,
		CountOfChecks: t.CountOfChecks,
		Status:        t.Status,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
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
