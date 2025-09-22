package domain

import "airport-tools-backend/pkg/e"

// Store описывает склад
type Store struct {
	Id        int64
	StationId int64
	Name      string

	StationObj *Station
	Locations  []*Location
}

func NewStore(stationId int64, name string) *Store {
	return &Store{
		StationId: stationId,
		Name:      name,
	}
}

func (s *Store) ChangeName(name string) error {
	if s.Name == name {
		return e.ErrNothingToChange
	}

	s.Name = name
	return nil
}
