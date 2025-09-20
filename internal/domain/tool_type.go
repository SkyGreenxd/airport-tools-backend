package domain

// TODO: в модели для GORM добавить ф-цию func TableName() string
// ToolType описывает тип инструмента
type ToolType struct {
	Id          int64
	PartNumber  string // партийный номер
	Description string // описание
	Co          string // состав/группа сплавов
	MC          string // код марки/сплава
}

func NewToolType(partNumber, description, co, mc string) *ToolType {
	return &ToolType{
		PartNumber:  partNumber,
		Description: description,
		Co:          co,
		MC:          mc,
	}
}
