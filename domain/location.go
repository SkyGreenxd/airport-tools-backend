package domain

// Station описывает аэропорт
type Station struct {
	Id   int64
	Code string // код аэропорта IATA = 3 буквы (на скрине нет, ICAO = 4 буквы)
}

// Store описывает склад
type Store struct {
	Id        int64
	StationId int64
	Name      string
}

// Location описывает ячейку на складе
type Location struct {
	Id      int64
	StoreId int64
	Name    string
}

func NewStation(code string) *Station {
	return &Station{
		Code: code,
	}
}

func NewStore(stationId int64, name string) *Store {
	return &Store{
		StationId: stationId,
		Name:      name,
	}
}

func NewLocation(storeId int64, name string) *Location {
	return &Location{
		StoreId: storeId,
		Name:    name,
	}
}
