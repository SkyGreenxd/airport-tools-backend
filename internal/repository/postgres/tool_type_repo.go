package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type ToolTypeRepository struct {
	DB *gorm.DB
}

func NewToolTypeRepository(db *gorm.DB) *ToolTypeRepository {
	return &ToolTypeRepository{
		DB: db,
	}
}

func (t *ToolTypeRepository) Create(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error) {
	const op = "ToolTypeRepository.Create"

	model := toToolTypeModel(toolType)
	result := t.DB.WithContext(ctx).Create(model)
	if err := postgresDuplicate(result, e.ErrToolTypeExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainToolType(model), nil
}

func (t *ToolTypeRepository) GetById(ctx context.Context, id int64) (*domain.ToolType, error) {
	const op = "ToolTypeRepository.GetById"

	var model ToolTypeModel
	result := t.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrToolTypeNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainToolType(&model), nil
}

func (t *ToolTypeRepository) GetAll(ctx context.Context) ([]*domain.ToolType, error) {
	const op = "ToolTypeRepository.GetAll"

	var models []*ToolTypeModel
	result := t.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	toolTypes := make([]*domain.ToolType, len(models))
	for i, m := range models {
		toolTypes[i] = toDomainToolType(m)
	}

	return toolTypes, nil
}

func (t *ToolTypeRepository) Delete(ctx context.Context, id int64) error {
	const op = "ToolTypeRepository.Delete"

	var model ToolTypeModel
	result := t.DB.WithContext(ctx).Delete(&model, id)
	if err := postgresForeignKeyViolation(result, e.ErrToolTypeIsUsed); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (t *ToolTypeRepository) Update(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error) {
	const op = "ToolTypeRepository.Update"

	updates := map[string]interface{}{
		"name": toolType.Name,
	}

	var updToolType ToolTypeModel
	result := t.DB.WithContext(ctx).Model(&ToolTypeModel{}).Where("id = ?", toolType.Id).Updates(updates).Scan(&updToolType)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrToolTypeNotFound
	}

	return toDomainToolType(&updToolType), nil
}

func toToolTypeModel(t *domain.ToolType) *ToolTypeModel {
	return &ToolTypeModel{
		Id:                 t.Id,
		PartNumber:         t.PartNumber,
		Name:               t.Name,
		ReferenceEmbedding: pgvector.NewVector(t.ReferenceEmbedding),
	}
}

func toDomainToolType(t *ToolTypeModel) *domain.ToolType {
	return &domain.ToolType{
		Id:                 t.Id,
		PartNumber:         t.PartNumber,
		Name:               t.Name,
		ReferenceEmbedding: t.ReferenceEmbedding.Slice(),
	}
}

func toArrToolTypeModel(tools []*domain.ToolType) []*ToolTypeModel {
	models := make([]*ToolTypeModel, len(tools))
	for i, m := range tools {
		models[i] = toToolTypeModel(m)
	}

	return models
}

func toArrDomainToolType(models []*ToolTypeModel) []*domain.ToolType {
	toolTypes := make([]*domain.ToolType, len(models))
	for i, m := range models {
		toolTypes[i] = toDomainToolType(m)
	}

	return toolTypes
}
