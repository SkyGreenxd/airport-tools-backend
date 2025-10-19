package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/repository"
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

func (t *TransactionResolutionsRepo) GetByQAId(ctx context.Context, qaId int64) ([]*domain.TransactionResolution, error) {
	const op = "TransactionResolutionsRepo.GetByQAId"

	var models []*TransactionResolutionModel
	result := t.DB.WithContext(ctx).
		Preload("Transaction", func(db *gorm.DB) *gorm.DB {
			return db.Order("updated_at DESC")
		}).
		Preload("Transaction.User").
		Find(&models, "qa_employee_id = ?", qaId)
	if err := checkGetQueryResult(result, e.ErrTransactionResolutionsNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactionResolution(models), nil
}

func (t *TransactionResolutionsRepo) GetAllModelError(ctx context.Context) ([]*domain.TransactionResolution, error) {
	return t.getTransactionsWithErrorType(ctx, domain.ModelError)
}

func (t *TransactionResolutionsRepo) GetAllHumanError(ctx context.Context) ([]*domain.TransactionResolution, error) {
	return t.getTransactionsWithErrorType(ctx, domain.HumanError)
}

func (t *TransactionResolutionsRepo) GetTopHumanErrorUsers(ctx context.Context) ([]repository.HumanErrorStats, error) {
	const op = "TransactionRepository.GetTopHumanErrorUsers"

	var stats []repository.HumanErrorStats
	result := t.DB.WithContext(ctx).
		Table("transaction_resolutions AS tr").
		Select(`
			u.full_name AS full_name,
			u.employee_id AS employee_id,
			COUNT(tr.id) AS qa_hits_count
		`).
		Joins("JOIN transactions t ON tr.transaction_id = t.id").
		Joins("JOIN users u ON t.user_id = u.id").
		Where("tr.reason = ?", "HUMAN_ERR").
		Group("u.full_name, u.employee_id").
		Order("qa_hits_count DESC").
		Find(&stats)

	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return stats, nil
}

func (t *TransactionResolutionsRepo) getTransactionsWithErrorType(ctx context.Context, typeOfError domain.Reason) ([]*domain.TransactionResolution, error) {
	const op = "TransactionResolutionsRepo.getTransactionsWithErrorType"
	var models []*TransactionResolutionModel
	result := t.DB.WithContext(ctx).Where("reason = ?", typeOfError).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrTransactionResolutionsNotFound)
	}

	return toDomainArrTransactionResolution(models), nil
}

func (t *TransactionResolutionsRepo) GetMlErrorTransactions(ctx context.Context) ([]*domain.TransactionResolution, error) {
	const op = "TransactionResolutionsRepo.GetMlErrorTransactions"

	var models []*TransactionResolutionModel
	result := t.DB.WithContext(ctx).
		Preload("Transaction.CvScans", func(db *gorm.DB) *gorm.DB {
			return db.Order("updated_at DESC")
		}).
		Where("reason = ?", domain.ModelError).Find(&models)
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
		Reason:        transaction.Reason,
		Notes:         transaction.Notes,
		CreatedAt:     transaction.CreatedAt,
	}
}

func toDomainTransactionResolution(model *TransactionResolutionModel) *domain.TransactionResolution {
	return &domain.TransactionResolution{
		Id:            model.Id,
		TransactionId: model.TransactionId,
		QAEmployeeId:  model.QAEmployeeId,
		Reason:        model.Reason,
		Notes:         model.Notes,
		CreatedAt:     model.CreatedAt,
		Transaction:   toDomainTransaction(model.Transaction),
	}
}

func toDomainArrTransactionResolution(models []*TransactionResolutionModel) []*domain.TransactionResolution {
	res := make([]*domain.TransactionResolution, 0, len(models))
	for _, model := range models {
		res = append(res, toDomainTransactionResolution(model))
	}

	return res
}
