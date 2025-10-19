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
	result := u.DB.WithContext(ctx).Preload("Role").First(&model, "employee_id = ?", employeeId)
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

func (u *UserRepository) GetByEmployeeIdWithTransactionResolutions(ctx context.Context, employeeId string) (*domain.User, error) {
	const op = "UserRepository.GetByEmployeeIdWithTransactionResolutions"

	var model UserModel
	result := u.DB.WithContext(ctx).Preload("TransactionResolutions").First(&model, "employee_id = ?", employeeId)
	if err := checkGetQueryResult(result, e.ErrUserNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(&model), nil
}

func (u *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	const op = "UserRepository.GetAll"

	var models []*UserModel
	result := u.DB.WithContext(ctx).Preload("Role").Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrUserNotFound
	}

	return toArrDomainUser(models), nil
}

func (u *UserRepository) GetAllEngineersWithTransactions(ctx context.Context) ([]*domain.User, error) {
	const op = "UserRepository.GetAllWithTransactions"

	var models []*UserModel
	var engineerRole RoleModel
	if err := u.DB.WithContext(ctx).Where("name = ?", domain.Engineer).First(&engineerRole).Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	// Получить пользователей с этим role_id
	result := u.DB.WithContext(ctx).
		Preload("Transactions", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC")
		}).
		Where("role_id = ?", engineerRole.Id).
		Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrUserNotFound
	}

	return toArrDomainUser(models), nil
}

func (u *UserRepository) GetAllQa(ctx context.Context) ([]*domain.User, error) {
	const op = "UserRepository.GetAllQa"

	var models []*UserModel
	var qaRole RoleModel
	if err := u.DB.WithContext(ctx).Where("name = ?", domain.QualityAuditor).First(&qaRole).Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	result := u.DB.WithContext(ctx).
		Where("role_id = ?", qaRole.Id).
		Find(&models)
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

func toArrUserModel(models []*domain.User) []*UserModel {
	result := make([]*UserModel, len(models))
	for i, model := range models {
		result[i] = toUserModel(model)
	}

	return result
}

func toUserModel(u *domain.User) *UserModel {
	model := &UserModel{
		Id:         u.Id,
		EmployeeId: u.EmployeeId,
		FullName:   u.FullName,
		RoleId:     u.RoleId,
	}

	if u.Transactions != nil {
		model.Transactions = toModelArrTransactions(u.Transactions)
	}

	if u.Role != nil {
		model.Role = toRoleModel(u.Role)
	}

	return model
}

func toDomainUser(u *UserModel) *domain.User {
	user := &domain.User{
		Id:         u.Id,
		EmployeeId: u.EmployeeId,
		FullName:   u.FullName,
		RoleId:     u.RoleId,
	}

	if u.Role != nil {
		user.Role = toDomainRole(u.Role)
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
