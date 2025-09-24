package domain

import "airport-tools-backend/pkg/e"

// ToolType описывает тип инструмента
type ToolType struct {
	Id                 int64
	PartNumber         string    // партийный номер
	Name               string    // имя
	ReferenceImageHash string    // хэш упражнения
	ReferenceEmbedding []float32 // эмбеддинг

	ToolSets []*ToolSet
}

func NewToolType(partNumber, name, referenceImageHash string, referenceEmbedding []float32) *ToolType {
	return &ToolType{
		PartNumber:         partNumber,
		Name:               name,
		ReferenceImageHash: referenceImageHash,
		ReferenceEmbedding: referenceEmbedding,
	}
}

func (t *ToolType) ChangName(name string) error {
	if t.Name == name {
		return e.ErrNothingToChange
	}

	t.Name = name
	return nil
}
