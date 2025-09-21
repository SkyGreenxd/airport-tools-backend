package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type LocationsRepository struct {
	DB *gorm.DB
}

func NewLocationsRepository(db *gorm.DB) LocationsRepository {
	return LocationsRepository{
		DB: db,
	}
}

func (l *LocationsRepository) Create(ctx context.Context, location *domain.Location) (*domain.Location, error) {
	const op = "LocationRepository.Create"

	model := toLocationModel(location)
	result := l.DB.WithContext(ctx).Create(model)
	if err := postgresDuplicate(result, e.ErrLocationExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainLocation(model), nil
}

func (l *LocationsRepository) GetById(ctx context.Context, id int64) (*domain.Location, error) {
	const op = "LocationRepository.GetById"

	var model LocationModel
	result := l.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrLocationNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainLocation(&model), nil
}

func (l *LocationsRepository) GetByIdWithTools(ctx context.Context, id int64) (*domain.Location, error) {
	const op = "LocationRepository.GetByIdWithTools"

	var model LocationModel
	result := l.DB.WithContext(ctx).Preload("Tools").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrLocationNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainLocation(&model), nil
}

func (l *LocationsRepository) GetAll(ctx context.Context) ([]*domain.Location, error) {
	const op = "LocationRepository.GetAll"

	var models []*LocationModel
	result := l.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	locations := make([]*domain.Location, len(models))
	for i, model := range models {
		locations[i] = toDomainLocation(model)
	}

	return locations, nil
}

func (l *LocationsRepository) Delete(ctx context.Context, id int64) error {
	const op = "LocationRepository.Delete"

	var model LocationModel
	result := l.DB.WithContext(ctx).Delete(&model, "id = ?", id)
	if result.Error != nil {
		return e.Wrap(op, result.Error)
	}

	return nil
}

func (l *LocationsRepository) Update(ctx context.Context, location *domain.Location) (*domain.Location, error) {
	const op = "LocationRepository.Update"

	updates := map[string]interface{}{
		"name": location.Name,
	}

	result := l.DB.WithContext(ctx).Model(&LocationModel{}).Where("id = ?", location.Id).Updates(updates)
	if err := postgresDuplicate(result, e.ErrLocationExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrLocationNotFound)
	}

	updLocation, err := l.GetById(ctx, location.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updLocation, nil
}

func toLocationModel(l *domain.Location) *LocationModel {
	return &LocationModel{
		Id:      l.Id,
		StoreId: l.StoreId,
		Name:    l.Name,
		Store:   toStoreModel(l.StoreObj),
		Tools:   toModelArrTools(l.Tools),
	}
}

func toDomainLocation(l *LocationModel) *domain.Location {
	return &domain.Location{
		Id:       l.Id,
		StoreId:  l.StoreId,
		Name:     l.Name,
		StoreObj: toDomainStoreModel(l.Store),
		Tools:    toDomainArrTools(l.Tools),
	}
}

func toModelArrLocations(locations []*domain.Location) []*LocationModel {
	models := make([]*LocationModel, len(locations))
	for i, location := range locations {
		models[i] = toLocationModel(location)
	}

	return models
}

func toDomainArrLocations(models []*LocationModel) []*domain.Location {
	locations := make([]*domain.Location, len(models))
	for i, model := range models {
		locations[i] = toDomainLocation(model)
	}

	return locations
}
