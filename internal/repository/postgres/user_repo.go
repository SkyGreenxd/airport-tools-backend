package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	const op = "UserRepository.Create"

	model := toUserModel(user)
	result := u.DB.WithContext(ctx).Create(&model)
	if err := postgresDuplicate(result, e.ErrUserExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(model), nil
}

func (u *UserRepository) GetById(ctx context.Context, id int64) (*domain.User, error) {
	const op = "UserRepository.GetById"

	var model UserModel
	result := u.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrUserNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(&model), nil
}

func (u *UserRepository) GetByEmployeeId(ctx context.Context, employeeId int64) (*domain.User, error) {
	const op = "UserRepository.GetByEmployeeId"

	var model UserModel
	result := u.DB.WithContext(ctx).First(&model, "employee_id = ?", employeeId)
	if err := checkGetQueryResult(result, e.ErrUserNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(&model), nil
}

func (u *UserRepository) GetByFullName(ctx context.Context, fullName string) (*domain.User, error) {
	const op = "UserRepository.GetByFullName"

	var model UserModel
	result := u.DB.WithContext(ctx).First(&model, "full_name = ?", fullName)
	if err := checkGetQueryResult(result, e.ErrUserNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(&model), nil
}

func (u *UserRepository) GetByIdWithTransactions(ctx context.Context, id int64) (*domain.User, error) {
	const op = "UserRepository.GetByIdWithTransactions"

	var model UserModel
	result := u.DB.WithContext(ctx).Preload("Transactions").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrUserNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(&model), nil
}

func (u *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	const op = "UserRepository.GetAll"

	var models []*UserModel
	result := u.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrUserNotFound
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = toDomainUser(model)
	}

	return users, nil
}

func (u *UserRepository) Delete(ctx context.Context, id int64) error {
	const op = "UserRepository.Delete"

	var model UserModel
	result := u.DB.WithContext(ctx).Delete(&model, "id = ?", id)
	if err := postgresForeignKeyViolation(result, e.ErrUserInUse); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (u *UserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	const op = "UserRepository.Update"

	updates := map[string]interface{}{
		"full_name": user.FullName,
		"role":      user.Role,
	}

	result := u.DB.WithContext(ctx).Model(&UserModel{}).Where("id = ?", user.Id).Updates(updates)
	if err := postgresDuplicate(result, e.ErrLocationExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrUserNotFound
	}

	updUser, err := u.GetById(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return updUser, nil
}

func toUserModel(u *domain.User) *UserModel {
	return &UserModel{
		Id:           u.Id,
		EmployeeId:   u.EmployeeId,
		FullName:     u.FullName,
		Role:         u.Role,
		Transactions: toModelArrTransactions(u.Transactions),
	}
}

func toDomainUser(u *UserModel) *domain.User {
	return &domain.User{
		Id:           u.Id,
		EmployeeId:   u.EmployeeId,
		FullName:     u.FullName,
		Role:         u.Role,
		Transactions: toDomainArrTransactions(u.Transactions),
	}
}
