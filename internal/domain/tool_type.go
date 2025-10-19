package domain

import "airport-tools-backend/pkg/e"

type ToolType struct {
	Id                 int64
	PartNumber         string
	Name               string
	ReferenceEmbedding []float32

	ToolSets []*ToolSet
}

func NewToolType(partNumber, name string, referenceEmbedding []float32) *ToolType {
	return &ToolType{
		PartNumber:         partNumber,
		Name:               name,
		ReferenceEmbedding: referenceEmbedding,
	}
}

func (t *ToolType) ValidateName(newName string) error {
	if t.Name == newName {
		return e.ErrNothingToChange
	}

	return nil
}
