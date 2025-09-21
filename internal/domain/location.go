package domain

import "airport-tools-backend/pkg/e"

// Location описывает ячейку на складе
type Location struct {
	Id      int64
	StoreId int64
	Name    string

	StoreObj *Store
	Tools    []*Tool
}

func NewLocation(storeId int64, name string) *Location {
	return &Location{
		StoreId: storeId,
		Name:    name,
	}
}

func (l *Location) ChangeName(name string) error {
	if l.Name == name {
		return e.ErrNothingToChange
	}

	l.Name = name
	return nil
}
