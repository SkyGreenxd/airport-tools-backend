package domain

import "airport-tools-backend/pkg/e"

// Station описывает аэропорт
type Station struct {
	Id   int64
	Code string // код аэропорта IATA = 3 буквы (на скрине нет, ICAO = 4 буквы)

	Stores []*Store
}

func NewStation(code string) *Station {
	return &Station{
		Code: code,
	}
}

func (s *Station) ChangeCode(code string) error {
	if s.Code == code {
		return e.ErrNothingToChange
	}

	s.Code = code
	return nil
}
