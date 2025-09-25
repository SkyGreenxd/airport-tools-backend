package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type CvScanDetailRepository struct {
	DB *gorm.DB
}

func NewCvScanDetailRepository(db *gorm.DB) *CvScanDetailRepository {
	return &CvScanDetailRepository{
		DB: db,
	}
}

func (c *CvScanDetailRepository) Create(ctx context.Context, cvScanDetail *domain.CvScanDetail) (*domain.CvScanDetail, error) {
	const op = "CvScanDetailRepository.Create"

	model := toCvScanDetailModel(cvScanDetail)
	result := c.DB.WithContext(ctx).Create(model)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScanDetail(model), nil
}

func (c *CvScanDetailRepository) GetById(ctx context.Context, id int64) (*domain.CvScanDetail, error) {
	const op = "CvScanDetailRepository.GetById"

	var model CvScanDetailModel
	result := c.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrCvScanDetailNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainCvScanDetail(&model), nil
}

func (c *CvScanDetailRepository) GetByCvScanId(ctx context.Context, cvScanId int64) ([]*domain.CvScanDetail, error) {
	const op = "CvScanDetailRepository.GetByCvScanId"

	var model []*CvScanDetailModel
	result := c.DB.WithContext(ctx).Find(&model, "cv_scan_id = ?", cvScanId)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrCvScanDetailNotFound)
	}

	return toArrDomainCvScanDetail(model), nil
}

func toCvScanDetailModel(c *domain.CvScanDetail) *CvScanDetailModel {
	return &CvScanDetailModel{
		Id:                 c.Id,
		CvScanId:           c.CvScanId,
		DetectedToolTypeId: c.DetectedToolTypeId,
		Confidence:         c.Confidence,
		ImageHash:          c.ImageHash,
		Embedding:          pgvector.NewVector(c.Embedding),
	}
}

func toDomainCvScanDetail(c *CvScanDetailModel) *domain.CvScanDetail {
	return &domain.CvScanDetail{
		Id:                 c.Id,
		CvScanId:           c.CvScanId,
		DetectedToolTypeId: c.DetectedToolTypeId,
		Confidence:         c.Confidence,
		ImageHash:          c.ImageHash,
		Embedding:          c.Embedding.Slice(),
	}
}

func toArrDomainCvScanDetail(models []*CvScanDetailModel) []*domain.CvScanDetail {
	scans := make([]*domain.CvScanDetail, len(models))
	for i, model := range models {
		scans[i] = toDomainCvScanDetail(model)
	}

	return scans
}

func toArrCvScanDetailModel(scans []*domain.CvScanDetail) []*CvScanDetailModel {
	models := make([]*CvScanDetailModel, len(scans))
	for i, model := range scans {
		models[i] = toCvScanDetailModel(model)
	}

	return models
}
