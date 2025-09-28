package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type ToolSetRepository struct {
	DB *gorm.DB
}

func NewToolSetRepository(db *gorm.DB) *ToolSetRepository {
	return &ToolSetRepository{DB: db}
}

func (t *ToolSetRepository) Create(ctx context.Context, toolSet *domain.ToolSet) (*domain.ToolSet, error) {
	const op = "ToolSetRepository.Create"

	model := toToolSetModel(toolSet)
	result := t.DB.Where(model).Create(model)
	if err := postgresDuplicate(result, e.ErrToolSetExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainToolSet(model), nil
}

func (t *ToolSetRepository) GetById(ctx context.Context, id int64) (*domain.ToolSet, error) {
	const op = "ToolSetRepository.GetById"

	var model ToolSetModel
	result := t.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrToolSetNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainToolSet(&model), nil
}

func (t *ToolSetRepository) GetAll(ctx context.Context) ([]*domain.ToolSet, error) {
	const op = "ToolSetRepository.GetAll"

	var models []*ToolSetModel
	result := t.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	return toArrDomainToolSet(models), nil
}

func (t *ToolSetRepository) Delete(ctx context.Context, id int64) error {
	const op = "ToolSetRepository.Delete"

	var model ToolSetModel
	result := t.DB.WithContext(ctx).Delete(&model, "id = ?", id)
	if result.Error != nil {
		return e.Wrap(op, result.Error)
	}

	return nil
}

func (t *ToolSetRepository) Update(ctx context.Context, toolSet *domain.ToolSet) (*domain.ToolSet, error) {
	const op = "ToolSetRepository.Update"

	updates := map[string]interface{}{
		"name": toolSet.Name,
	}

	var updSet ToolSetModel
	result := t.DB.WithContext(ctx).Model(&ToolSetModel{}).Where("id = ?", toolSet.Id).Updates(updates).Scan(&updSet)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrTransactionNotFound)
	}

	return toDomainToolSet(&updSet), nil
}

func (t *ToolSetRepository) GetByIdWithTools(ctx context.Context, id int64) (*domain.ToolSet, error) {
	const op = "ToolSetRepository.GetByIdWithTools"

	var model ToolSetModel
	result := t.DB.WithContext(ctx).Preload("Tools").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrToolSetNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainToolSet(&model), nil
}

func toToolSetModel(t *domain.ToolSet) *ToolSetModel {
	model := &ToolSetModel{
		Id:   t.Id,
		Name: t.Name,
	}

	if t.Tools != nil {
		model.Tools = toArrToolTypeModel(t.Tools)
	}

	return model
}

func toDomainToolSet(t *ToolSetModel) *domain.ToolSet {
	set := &domain.ToolSet{
		Id:   t.Id,
		Name: t.Name,
	}

	if t.Tools != nil {
		set.Tools = toArrDomainToolType(t.Tools)
	}

	return set
}

func toArrDomainToolSet(models []*ToolSetModel) []*domain.ToolSet {
	sets := make([]*domain.ToolSet, len(models))
	for i, model := range models {
		sets[i] = toDomainToolSet(model)
	}

	return sets
}
