package postgres

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"

	"gorm.io/gorm"
)

type StationRepository struct {
	DB *gorm.DB
}

func NewStationRepository(db *gorm.DB) *StationRepository {
	return &StationRepository{
		DB: db,
	}
}

func (s *StationRepository) Create(ctx context.Context, station *domain.Station) (*domain.Station, error) {
	const op = "StationRepository.Create"
	model := toStationModel(station)

	result := s.DB.WithContext(ctx).Create(model)
	if err := postgresDuplicate(result, e.ErrStationExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainStation(model), nil
}

// TODO: мб нужен Preload
func (s *StationRepository) GetById(ctx context.Context, id int64) (*domain.Station, error) {
	const op = "StationRepository.GetById"
	var model StationModel

	result := s.DB.WithContext(ctx).First(&model, "id = ?", id)
	if err := checkGetQueryResult(result, e.ErrStationNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainStation(&model), nil
}

// TODO: мб нужен Preload
func (s *StationRepository) GetAll(ctx context.Context) ([]*domain.Station, error) {
	const op = "StationRepository.GetAll"
	var models []*StationModel

	result := s.DB.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, e.Wrap(op, err)
	}

	stations := make([]*domain.Station, len(models))
	for i, model := range models {
		stations[i] = toDomainStation(model)
	}

	return stations, nil
}

func (s *StationRepository) Delete(ctx context.Context, id int64) error {
	const op = "StationRepository.Delete"
	var model StationModel

	result := s.DB.WithContext(ctx).Delete(&model, "id = ?", id)
	if err := result.Error; err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (s *StationRepository) Update(ctx context.Context, station *domain.Station) (*domain.Station, error) {
	const op = "StationRepository.Update"

	updates := map[string]interface{}{
		"code": station.Code,
	}

	result := s.DB.WithContext(ctx).Model(&StationModel{}).Where("id = ?", station.Id).Updates(updates)
	if err := postgresDuplicate(result, e.ErrStationExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	if result.RowsAffected == 0 {
		return nil, e.ErrStationNotFound
	}

	updStation, err := s.GetById(ctx, station.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return updStation, nil
}

func toStationModel(s *domain.Station) *StationModel {
	return &StationModel{
		Code: s.Code,
	}
}

func toDomainStation(s *StationModel) *domain.Station {
	return &domain.Station{
		Code: s.Code,
	}
}
