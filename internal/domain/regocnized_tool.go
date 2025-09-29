package domain

// RecognizedTool описывает инструмент, который был распознан CV
type RecognizedTool struct {
	ToolTypeId int64
	Confidence float32
	Embedding  []float32
}

func NewRecognizedTool(toolTypeId int64, confidence float32, embedding []float32) *RecognizedTool {
	return &RecognizedTool{
		ToolTypeId: toolTypeId,
		Confidence: confidence,
		Embedding:  embedding,
	}
}
