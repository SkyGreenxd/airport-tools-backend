package domain

import "airport-tools-backend/pkg/e"

// ToolType описывает тип инструмента
type ToolType struct {
	Id          int64
	PartNumber  string // партийный номер
	Description string // описание
	//Co          string // состав/группа сплавов
	//MC          string // код марки/сплава

	Tools []*Tool
}

func NewToolType(partNumber, description string) *ToolType {
	return &ToolType{
		PartNumber:  partNumber,
		Description: description,
		//Co:          co,
		//MC:          mc,
	}
}

func (t *ToolType) ChangeDescription(description string) error {
	if t.Description == description {
		return e.ErrNothingToChange
	}

	t.Description = description
	return nil
}
