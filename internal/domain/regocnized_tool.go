package domain

// RecognizedTool более легкая сущность для работы в сервисе
type RecognizedTool struct {
	ToolTypeId int64
	Confidence float32
	HashTool   string
	Embedding  []float32
}

func NewRecognizedTool(toolTypeId int64, confidence float32, hashTool string, embedding []float32) *RecognizedTool {
	return &RecognizedTool{
		ToolTypeId: toolTypeId,
		Confidence: confidence,
		HashTool:   hashTool,
		Embedding:  embedding,
	}
}
