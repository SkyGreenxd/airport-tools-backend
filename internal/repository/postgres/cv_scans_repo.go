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

func (c *CvStoreRepository) Update(ctx context.Context, cvScan *domain.CvScan) (*domain.CvScan, error) {
	const op = "CvStoreRepository.Update"

	updates := map[string]interface{}{
		"status": cvScan.Status,
		"reason": cvScan.Reason,
	}

	result := c.DB.WithContext(ctx).Model(&CvScanModel{}).Where("id = ?", cvScan.Id).Updates(updates)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrCvScanNotFound)
	}

	updCvScan, err := c.GetById(ctx, cvScan.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updCvScan, nil
}

func toCvScanModel(c *domain.CvScan) *CvScanModel {
	return &CvScanModel{
		Id:            c.Id,
		TransactionId: c.TransactionId,
		Status:        c.Status,
		Reason:        c.Reason,
		Photo:         c.Photo,
	}
}

func toDomainCvScan(c *CvScanModel) *domain.CvScan {
	return &domain.CvScan{
		Id:            c.Id,
		TransactionId: c.TransactionId,
		Status:        c.Status,
		Reason:        c.Reason,
		Photo:         c.Photo,
	}
}
