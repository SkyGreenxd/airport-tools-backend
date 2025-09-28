package domain

import "airport-tools-backend/pkg/e"

type ToolType struct {
	Id                 int64
	PartNumber         string
	Name               string
	ReferenceImageHash string
	ReferenceEmbedding []float32

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

func (t *ToolType) ValidateName(newName string) error {
	if t.Name == newName {
		return e.ErrNothingToChange
	}

	return nil
}
