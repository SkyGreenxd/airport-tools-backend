package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type CvStoreRepository struct {
	DB *gorm.DB
}

func NewCvStoreRepository(db *gorm.DB) *CvStoreRepository {
	return &CvStoreRepository{
		DB: db,
	}
}

func (c *CvStoreRepository) Create(ctx context.Context, cvScan *domain.CvScan) (*domain.CvScan, error) {
	const op = "CvStoreRepository.Create"

	model := toCvScanModel(cvScan)
	result := c.DB.WithContext(ctx).Create(model)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(model), nil
}

func (c *CvStoreRepository) GetById(ctx context.Context, id int64) (*domain.CvScan, error) {
	const op = "CvStoreRepository.GetById"

	var model CvScanModel
	result := c.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrCvScanNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(&model), nil
}

func (c *CvStoreRepository) GetByTransactionId(ctx context.Context, transactionId int64) (*domain.CvScan, error) {
	const op = "CvStoreRepository.GetByTransactionId"

	var model CvScanModel
	result := c.DB.WithContext(ctx).First(&model, "transaction_id = ?", transactionId)
	if err := checkGetQueryResult(result, e.ErrCvScanNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(&model), nil
}

func (c *CvStoreRepository) GetByIdWithTransaction(ctx context.Context, id int64) (*domain.CvScan, error) {
	const op = "CvStoreRepository.GetByIdWithTransaction"

	var model CvScanModel
	result := c.DB.WithContext(ctx).Preload("Transaction").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrCvScanNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(&model), nil
}

func (c *CvStoreRepository) GetByIdWithDetectedTools(ctx context.Context, id int64) (*domain.CvScan, error) {
	const op = "CvStoreRepository.GetByIdWithDetectedTools"

	var model CvScanModel
	result := c.DB.WithContext(ctx).Preload("DetectedTools").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrCvScanNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScan(&model), nil
}

func toCvScanModel(c *domain.CvScan) *CvScanModel {
	return &CvScanModel{
		Id:            c.Id,
		TransactionId: c.TransactionId,
		ScanType:      c.ScanType,
		ImageUrl:      c.ImageUrl,
		Transaction:   toTransactionModel(c.TransactionObj),
		DetectedTools: toArrCvScanDetailModel(c.DetectedTools),
	}
}

func toDomainCvScan(c *CvScanModel) *domain.CvScan {
	return &domain.CvScan{
		Id:             c.Id,
		TransactionId:  c.TransactionId,
		ScanType:       c.ScanType,
		ImageUrl:       c.ImageUrl,
		TransactionObj: toDomainTransaction(c.Transaction),
		DetectedTools:  toArrDomainCvScanDetail(c.DetectedTools),
	}
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
