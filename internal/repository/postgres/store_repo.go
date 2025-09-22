package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type StoreRepository struct {
	DB *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{DB: db}
}

func (s *StoreRepository) Create(ctx context.Context, store *domain.Store) (*domain.Store, error) {
	const op = "StoreRepository.Create"

	model := toStoreModel(store)
	result := s.DB.WithContext(ctx).Create(model)
	if err := postgresDuplicate(result, e.ErrStoreExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainStoreModel(model), nil
}

func (s *StoreRepository) GetById(ctx context.Context, id int64) (*domain.Store, error) {
	const op = "StoreRepository.GetById"

	var model StoreModel
	result := s.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrStoreNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainStoreModel(&model), nil
}

func (s *StoreRepository) GetAll(ctx context.Context) ([]*domain.Store, error) {
	const op = "StoreRepository.GetAll"

	var models []*StoreModel
	result := s.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	stores := make([]*domain.Store, len(models))
	for i, model := range models {
		stores[i] = toDomainStoreModel(model)
	}

	return stores, nil
}

func (s *StoreRepository) Delete(ctx context.Context, id int64) error {
	const op = "StoreRepository.Delete"
	var model StoreModel

	result := s.DB.WithContext(ctx).Delete(&model, "id = ?", id)
	if result.Error != nil {
		return e.Wrap(op, result.Error)
	}

	return nil
}

func (s *StoreRepository) Update(ctx context.Context, store *domain.Store) (*domain.Store, error) {
	const op = "StoreRepository.Update"

	updates := map[string]interface{}{
		"name": store.Name,
	}

	result := s.DB.WithContext(ctx).Model(&StoreModel{}).Where("id = ?", store.Id).Updates(updates)
	if err := postgresDuplicate(result, e.ErrStoreExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.Wrap(op, e.ErrStoreNotFound)
	}

	updStores, err := s.GetById(ctx, store.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updStores, nil
}

func (s *StoreRepository) GetByIdWithStation(ctx context.Context, id int64) (*domain.Store, error) {
	const op = "StoreRepository.GetByIdWithStation"

	var model StoreModel
	result := s.DB.WithContext(ctx).Preload("Station").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrStoreNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainStoreModel(&model), nil
}

func (s *StoreRepository) GetByIdWithLocations(ctx context.Context, id int64) (*domain.Store, error) {
	const op = "StoreRepository.GetByIdWithLocations"

	var model StoreModel
	result := s.DB.WithContext(ctx).Preload("Locations").First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrStoreNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainStoreModel(&model), nil
}

func toStoreModel(s *domain.Store) *StoreModel {
	return &StoreModel{
		Id:        s.Id,
		StationId: s.StationId,
		Name:      s.Name,
		Station:   toStationModel(s.StationObj),
		Locations: toModelArrLocations(s.Locations),
	}
}

func toDomainStoreModel(s *StoreModel) *domain.Store {
	return &domain.Store{
		Id:         s.Id,
		StationId:  s.StationId,
		Name:       s.Name,
		StationObj: toDomainStation(s.Station),
		Locations:  toDomainArrLocations(s.Locations),
	}
}

func toModelArrStores(stores []*domain.Store) []*StoreModel {
	models := make([]*StoreModel, len(stores))
	for i, store := range stores {
		models[i] = toStoreModel(store)
	}

	return models
}

func toDomainArrStores(models []*StoreModel) []*domain.Store {
	stores := make([]*domain.Store, len(models))
	for i, model := range models {
		stores[i] = toDomainStoreModel(model)
	}

	return stores
}
