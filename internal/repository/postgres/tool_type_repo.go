package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

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

func (t *ToolTypeRepository) GetByIdWithTools(ctx context.Context, id int64) (*domain.ToolType, error) {
	const op = "ToolTypeRepository.GetByIdWithTools"

	var model ToolTypeModel
	result := t.DB.WithContext(ctx).Preload("Tools").First(&model, "id = ?", id)
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
	result := t.DB.WithContext(ctx).Delete(&model)
	if err := postgresForeignKeyViolation(result, e.ErrToolTypeIsUsed); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (t *ToolTypeRepository) Update(ctx context.Context, toolType *domain.ToolType) (*domain.ToolType, error) {
	const op = "ToolTypeRepository.Update"

	updates := map[string]interface{}{
		"description": toolType.Description,
	}
	result := t.DB.WithContext(ctx).Model(&ToolTypeModel{}).Where("id = ?", toolType.Id).Updates(updates)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrToolTypeNotFound
	}

	updToolType, err := t.GetById(ctx, toolType.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updToolType, nil
}

func toToolTypeModel(t *domain.ToolType) *ToolTypeModel {
	return &ToolTypeModel{
		Id:          t.Id,
		PartNumber:  t.PartNumber,
		Description: t.Description,
		Tools:       toModelArrTools(t.Tools),
	}
}

func toDomainToolType(t *ToolTypeModel) *domain.ToolType {
	return &domain.ToolType{
		Id:          t.Id,
		PartNumber:  t.PartNumber,
		Description: t.Description,
		Tools:       toDomainArrTools(t.Tools),
	}
}
