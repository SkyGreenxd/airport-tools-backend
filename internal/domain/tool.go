package domain

import "time"

// Tool описывает конкретный инструмент
type Tool struct {
	Id         int64
	TypeToolId int64      // ID на тип инструмента
	ToirId     int64      // Id инструмента в системе ТОиР
	LocationId int64      // местоположение инструмента
	SnBn       string     // FIXME: Я НЕ ЗНАЮ ЧТО ЭТО
	ExpiresAt  *time.Time // дата проверки
	// Status string - на складе, у инженера, в ремонте
}

func NewTool(typeToolId, toirId, locationId int64, snBn string, expiresAt *time.Time) *Tool {
	return &Tool{
		TypeToolId: typeToolId,
		ToirId:     toirId,
		LocationId: locationId,
		SnBn:       snBn,
		ExpiresAt:  expiresAt,
	}
}
