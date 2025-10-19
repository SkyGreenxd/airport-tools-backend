package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type RoleRepo struct {
	DB *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		DB: db,
	}
}

func (r *RoleRepo) Create(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	const op = "RoleRepo.Create"

	model := toRoleModel(role)
	result := r.DB.WithContext(ctx).Create(&model)
	if err := postgresDuplicate(result, e.ErrRoleExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainRole(model), nil
}

func (r *RoleRepo) GetAll(ctx context.Context) ([]*domain.Role, error) {
	const op = "RoleRepo.GetAll"

	var models []RoleModel
	result := r.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrRoleNotFound)
	}

	return toDomainArrRoles(models), nil
}

func (r *RoleRepo) GetById(ctx context.Context, id int64) (*domain.Role, error) {
	const op = "RoleRepo.GetById"

	var model RoleModel
	result := r.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrRoleNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainRole(&model), nil
}

func (r *RoleRepo) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	const op = "RoleRepo.GetByName"

	var model RoleModel
	result := r.DB.WithContext(ctx).First(&model, "name = ?", name)
	if err := checkGetQueryResult(result, e.ErrRoleNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainRole(&model), nil
}

func toRoleModel(role *domain.Role) *RoleModel {
	model := &RoleModel{
		Id:   role.Id,
		Name: role.Name,
	}

	if role.Users != nil {
		model.Users = toArrUserModel(role.Users)
	}

	return model
}

func toDomainRole(model *RoleModel) *domain.Role {
	role := &domain.Role{
		Id:   model.Id,
		Name: model.Name,
	}

	if model.Users != nil {
		role.Users = toArrDomainUser(model.Users)
	}

	return role
}

func toDomainArrRoles(models []RoleModel) []*domain.Role {
	roles := make([]*domain.Role, len(models))
	for i, model := range models {
		roles[i] = toDomainRole(&model)
	}

	return roles
}
