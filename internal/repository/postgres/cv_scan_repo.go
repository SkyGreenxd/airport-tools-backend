package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type CvScanRepository struct {
	DB *gorm.DB
}

func NewCvScanRepository(db *gorm.DB) *CvScanRepository {
	return &CvScanRepository{
		DB: db,
	}
}

func (c *CvScanRepository) Create(ctx context.Context, cvScan *domain.CvScan) (*domain.CvScan, error) {
	const op = "CvScanRepository.Create"

	model := toCvScanModel(cvScan)
	result := c.DB.WithContext(ctx).Create(model)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(model), nil
}

func (c *CvScanRepository) GetById(ctx context.Context, id int64) (*domain.CvScan, error) {
	const op = "CvScanRepository.GetById"

	var model CvScanModel
	result := c.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrCvScanNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(&model), nil
}

func (c *CvScanRepository) GetByTransactionId(ctx context.Context, transactionId int64) (*domain.CvScan, error) {
	const op = "CvScanRepository.GetByTransactionId"

	var model CvScanModel
	result := c.DB.WithContext(ctx).First(&model, "transaction_id = ?", transactionId)
	if err := checkGetQueryResult(result, e.ErrCvScanNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(&model), nil
}

func (c *CvScanRepository) GetByIdWithTransaction(ctx context.Context, id int64) (*domain.CvScan, error) {
	const op = "CvScanRepository.GetByIdWithTransaction"

	var model CvScanModel
	result := c.DB.WithContext(ctx).Preload("Transaction").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrCvScanNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(&model), nil
}

func (c *CvScanRepository) GetByIdWithDetectedTools(ctx context.Context, id int64) (*domain.CvScan, error) {
	const op = "CvScanRepository.GetByIdWithDetectedTools"

	var model CvScanModel
	result := c.DB.WithContext(ctx).Preload("DetectedTools").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrCvScanNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(&model), nil
}

func toCvScanModel(c *domain.CvScan) *CvScanModel {
	model := &CvScanModel{
		Id:            c.Id,
		TransactionId: c.TransactionId,
		ScanType:      c.ScanType,
		ImageUrl:      c.ImageUrl,
		CreatedAt:     c.CreatedAt,
	}

	if c.DetectedTools != nil {
		model.Transaction = toTransactionModel(c.TransactionObj)
	}

	if c.TransactionObj != nil {
		model.DetectedTools = toArrCvScanDetailModel(c.DetectedTools)
	}

	return model
}

func toDomainCvScan(c *CvScanModel) *domain.CvScan {
	scan := &domain.CvScan{
		Id:            c.Id,
		TransactionId: c.TransactionId,
		ScanType:      c.ScanType,
		ImageUrl:      c.ImageUrl,
		CreatedAt:     c.CreatedAt,
	}

	if c.DetectedTools != nil {
		scan.TransactionObj = toDomainTransaction(c.Transaction)
	}

	if c.Transaction != nil {
		scan.DetectedTools = toArrDomainCvScanDetail(c.DetectedTools)
	}

	return scan
}

func toArrDomainCvScans(models []*CvScanModel) []*domain.CvScan {
	result := make([]*domain.CvScan, len(models))
	for i, model := range models {
		result[i] = toDomainCvScan(model)
	}

	return result
}

func toArrCvScansModel(scans []*domain.CvScan) []*CvScanModel {
	result := make([]*CvScanModel, len(scans))
	for i, model := range scans {
		result[i] = toCvScanModel(model)
	}

	return result
}
