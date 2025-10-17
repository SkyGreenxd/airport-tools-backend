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

func (u *UserRepository) GetByEmployeeId(ctx context.Context, employeeId string) (*domain.User, error) {
	const op = "UserRepository.GetByEmployeeId"

	var model UserModel
	result := u.DB.WithContext(ctx).First(&model, "employee_id = ?", employeeId)
	if err := checkGetQueryResult(result, e.ErrUserNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(&model), nil
}

func (u *UserRepository) GetByEmployeeIdWithTransactions(ctx context.Context, employeeId string) (*domain.User, error) {
	const op = "UserRepository.GetByIdWithTransactions"

	var model UserModel
	result := u.DB.WithContext(ctx).Preload("Transactions").First(&model, "employee_id = ?", employeeId)
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

	return toArrDomainUser(models), nil
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
		"employee_id": user.EmployeeId,
		"full_name":   user.FullName,
		"role":        user.Role,
	}

	var updUser UserModel
	result := u.DB.WithContext(ctx).Model(&UserModel{}).Where("id = ?", user.Id).Updates(updates).Scan(&updUser)
	if err := postgresDuplicate(result, e.ErrUserNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrUserNotFound
	}

	return toDomainUser(&updUser), nil
}

func toUserModel(u *domain.User) *UserModel {
	model := &UserModel{
		Id:         u.Id,
		EmployeeId: u.EmployeeId,
		FullName:   u.FullName,
		Role:       u.Role,
	}

	if u.Transactions != nil {
		model.Transactions = toModelArrTransactions(u.Transactions)
	}

	return model
}

func toDomainUser(u *UserModel) *domain.User {
	user := &domain.User{
		Id:         u.Id,
		EmployeeId: u.EmployeeId,
		FullName:   u.FullName,
		Role:       u.Role,
	}

	if u.Transactions != nil {
		user.Transactions = toDomainArrTransactions(u.Transactions)
	}

	return user
}

func toArrDomainUser(u []*UserModel) []*domain.User {
	result := make([]*domain.User, len(u))
	for i, u := range u {
		result[i] = toDomainUser(u)
	}

	return result
}
