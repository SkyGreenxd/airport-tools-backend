package domain

// Station описывает аэропорт
type Station struct {
	Id   int64
	Code string // код аэропорта IATA = 3 буквы (на скрине нет, ICAO = 4 буквы)
}

func NewStation(code string) *Station {
	return &Station{
		Code: code,
	}
}
