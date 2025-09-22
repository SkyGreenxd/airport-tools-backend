package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type ToolRepository struct {
	DB *gorm.DB
}

func NewToolRepository(db *gorm.DB) *ToolRepository {
	return &ToolRepository{
		DB: db,
	}
}

func (t *ToolRepository) Create(ctx context.Context, tool *domain.Tool) (*domain.Tool, error) {
	const op = "ToolRepository.Create"

	model := toToolModel(tool)
	result := t.DB.WithContext(ctx).Create(model)
	if err := postgresDuplicate(result, e.ErrToolExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTool(model), nil
}

func (t *ToolRepository) GetById(ctx context.Context, id int64) (*domain.Tool, error) {
	const op = "ToolRepository.GetById"

	var model ToolModel
	result := t.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrToolNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTool(&model), nil
}

func (t *ToolRepository) GetAll(ctx context.Context) ([]*domain.Tool, error) {
	const op = "ToolRepository.GetAll"

	var models []*ToolModel
	result := t.DB.WithContext(ctx).Find(&models)
	if result.Error != nil {
		return nil, e.Wrap(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrToolNotFound)
	}

	tools := make([]*domain.Tool, len(models))
	for i, tool := range models {
		tools[i] = toDomainTool(tool)
	}

	return tools, nil
}

func (t *ToolRepository) Delete(ctx context.Context, id int64) error {
	const op = "ToolRepository.Delete"

	var model ToolModel
	result := t.DB.WithContext(ctx).Delete(&model, "id = ?", id)
	if result.Error != nil {
		return e.Wrap(op, result.Error)
	}

	return nil
}

func (t *ToolRepository) Update(ctx context.Context, tool *domain.Tool) (*domain.Tool, error) {
	const op = "ToolRepository.Update"

	updates := map[string]interface{}{
		"type_tool_id": tool.TypeToolId,
		"location_id":  tool.LocationId,
		//"sn_bn":        tool.SnBn,
		"expires_at": tool.ExpiresAt,
	}

	result := t.DB.WithContext(ctx).Model(&ToolModel{}).Where("id = ?", tool.Id).Updates(updates)
	if err := postgresDuplicate(result, e.ErrToolExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrToolNotFound)
	}

	updTool, err := t.GetById(ctx, tool.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updTool, nil
}

func (t *ToolRepository) GetByIdWithToolType(ctx context.Context, toolId int64) (*domain.Tool, error) {
	const op = "ToolRepository.GetByIdWithToolType"

	var model ToolModel
	result := t.DB.WithContext(ctx).Preload("ToolType").First(&model, "id = ?", toolId)
	if err := checkGetQueryResult(result, e.ErrToolNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTool(&model), nil
}

func (t *ToolRepository) GetByIdWithLocation(ctx context.Context, toolId int64) (*domain.Tool, error) {
	const op = "ToolRepository.GetByIdWithLocation"

	var model ToolModel
	result := t.DB.WithContext(ctx).Preload("Location").First(&model, "id = ?", toolId)
	if err := checkGetQueryResult(result, e.ErrToolNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTool(&model), nil
}

func (t *ToolRepository) GetByIdWithAllData(ctx context.Context, toolId int64) (*domain.Tool, error) {
	const op = "ToolRepository.GetByIdWithAllData"

	var model ToolModel
	result := t.DB.WithContext(ctx).
		Preload("ToolType").
		Preload("Location").
		First(&model, "id = ?", toolId)
	if err := checkGetQueryResult(result, e.ErrToolNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTool(&model), nil
}

// TODO: изменить мапперы
func toToolModel(t *domain.Tool) *ToolModel {
	return &ToolModel{
		Id:         t.Id,
		TypeToolId: t.TypeToolId,
		ToirId:     t.ToirId,
		LocationId: t.LocationId,
		//SnBn:       t.SnBn,
		ExpiresAt: t.ExpiresAt,
	}
}

func toDomainTool(tool *ToolModel) *domain.Tool {
	return &domain.Tool{
		Id:         tool.Id,
		TypeToolId: tool.TypeToolId,
		ToirId:     tool.ToirId,
		LocationId: tool.LocationId,
		//SnBn:       tool.SnBn,
		ExpiresAt: tool.ExpiresAt,
	}
}

func toModelArrTools(tools []*domain.Tool) []*ToolModel {
	models := make([]*ToolModel, len(tools))
	for i, tool := range tools {
		models[i] = toToolModel(tool)
	}

	return models
}

func toDomainArrTools(models []*ToolModel) []*domain.Tool {
	tools := make([]*domain.Tool, len(models))
	for i, model := range models {
		tools[i] = toDomainTool(model)
	}

	return tools
}
