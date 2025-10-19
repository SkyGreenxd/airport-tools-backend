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

func (t *TransactionResolutionsRepo) Create(ctx context.Context, transaction *domain.TransactionResolution, toolIds []int64) (*domain.TransactionResolution, error) {
	const op = "TransactionResolutionsRepo.Create"

	model := toTransactionResolutionModel(transaction)

	if len(toolIds) != 0 {
		var tools []*ToolTypeModel
		if err := t.DB.WithContext(ctx).Where("id IN ?", toolIds).Find(&tools).Error; err != nil {
			return nil, e.Wrap(op, err)
		}

		model.Tools = tools
	}

	result := t.DB.WithContext(ctx).Create(&model)
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
		Preload("Transaction.User").Order("id desc").
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
		Preload("Transaction.CvScans").
		Where("reason = ?", domain.ModelError).
		Order("transaction_id desc").
		Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainArrTransactionResolution(models), nil
}

func (t *TransactionResolutionsRepo) GetMlErrorTools(ctx context.Context) ([]*repository.ToolSetWithErrors, error) {
	const op = "TransactionResolutionsRepo.GetMlErrorTools"

	type toolErrorCount struct {
		ToolTypeId   int64
		MLErrorCount int64
	}

	// Считаем ML-ошибки сразу для всех инструментов
	var counts []toolErrorCount
	if err := t.DB.WithContext(ctx).
		Model(&ModelErrItemModel{}).
		Select("tool_type_id, COUNT(*) AS ml_error_count").
		Group("tool_type_id").
		Scan(&counts).Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	// Создаем мапу tool_id -> ml_error_count
	countMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		countMap[c.ToolTypeId] = c.MLErrorCount
	}

	// Загружаем все сеты с инструментами
	var toolSets []ToolSetModel
	if err := t.DB.WithContext(ctx).Preload("Tools").Find(&toolSets).Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	var result []*repository.ToolSetWithErrors
	for _, ts := range toolSets {
		tsWithErrors := repository.ToolSetWithErrors{
			ID:    ts.Id,
			Name:  ts.Name,
			Tools: []repository.ToolWithErrorCount{},
		}

		for _, tool := range ts.Tools {
			tsWithErrors.Tools = append(tsWithErrors.Tools, repository.ToolWithErrorCount{
				ID:           tool.Id,
				Name:         tool.Name,
				MLErrorCount: countMap[tool.Id],
			})
		}

		result = append(result, &tsWithErrors)
	}

	return result, nil
}

func toTransactionResolutionModel(transaction *domain.TransactionResolution) *TransactionResolutionModel {
	model := &TransactionResolutionModel{
		Id:            transaction.Id,
		TransactionId: transaction.TransactionId,
		QAEmployeeId:  transaction.QAEmployeeId,
		Reason:        transaction.Reason,
		Notes:         transaction.Notes,
		CreatedAt:     transaction.CreatedAt,
	}

	if transaction.Tools != nil {
		model.Tools = toArrToolTypeModel(transaction.Tools)
	}

	return model
}

func toDomainTransactionResolution(model *TransactionResolutionModel) *domain.TransactionResolution {
	tr := &domain.TransactionResolution{
		Id:            model.Id,
		TransactionId: model.TransactionId,
		QAEmployeeId:  model.QAEmployeeId,
		Reason:        model.Reason,
		Notes:         model.Notes,
		CreatedAt:     model.CreatedAt,
		Transaction:   toDomainTransaction(model.Transaction),
	}

	if model.Tools != nil {
		tr.Tools = toArrDomainToolType(model.Tools)
	}

	return tr
}

func toDomainArrTransactionResolution(models []*TransactionResolutionModel) []*domain.TransactionResolution {
	res := make([]*domain.TransactionResolution, 0, len(models))
	for _, model := range models {
		res = append(res, toDomainTransactionResolution(model))
	}

	return res
}
