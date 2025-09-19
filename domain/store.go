package domain

// Store описывает склад
type Store struct {
	Id        int64
	StationId int64
	Name      string
}

func NewStore(stationId int64, name string) *Store {
	return &Store{
		StationId: stationId,
		Name:      name,
	}
}
