package domain

// Location описывает ячейку на складе
type Location struct {
	Id      int64
	StoreId int64
	Name    string
}

func NewLocation(storeId int64, name string) *Location {
	return &Location{
		StoreId: storeId,
		Name:    name,
	}
}
