package domain

import (
	"airport-tools-backend/pkg/e"
	"time"
)

// Tool описывает конкретный инструмент
type Tool struct {
	Id         int64
	TypeToolId int64      // ID на тип инструмента
	ToirId     int64      // Id инструмента в системе ТОиР
	LocationId int64      // местоположение инструмента
	ExpiresAt  *time.Time // дата проверки
	// SnBn       string     // FIXME: Я НЕ ЗНАЮ ЧТО ЭТО
	// Status string - на складе, у инженера, в ремонте

	ToolTypeObj *ToolType
	LocationObj *Location
}

func NewTool(typeToolId, toirId, locationId int64, expiresAt *time.Time) *Tool {
	return &Tool{
		TypeToolId: typeToolId,
		ToirId:     toirId,
		LocationId: locationId,
		ExpiresAt:  expiresAt,
		//SnBn:       snBn,
	}
}

func (t *Tool) ChangeLocationId(id int64) error {
	if t.LocationId == id {
		return e.ErrNothingToChange
	}

	t.LocationId = id
	return nil
}

func (t *Tool) ChangeExpiresAt(expiresAt time.Time) error {
	if t.ExpiresAt == nil {
		t.ExpiresAt = &expiresAt
		return nil
	}

	if *t.ExpiresAt == expiresAt {
		return e.ErrNothingToChange
	}

	if expiresAt.Before(*t.ExpiresAt) {
		return e.ErrIncorrectDate
	}

	t.ExpiresAt = &expiresAt
	return nil
}
