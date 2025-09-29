package usecase

import (
	"airport-tools-backend/internal/domain"
	"math"
)

// cosineSimilarity вычисляет косинусное сходство между двумя векторами
func cosineSimilarity(reference, recognized []float32) float64 {
	var dot, normReference, normRecognized float64
	for i := range reference {
		dot += float64(reference[i] * recognized[i])
		normReference += float64(reference[i] * reference[i])
		normRecognized += float64(recognized[i] * recognized[i])
	}
	return dot / (math.Sqrt(normReference) * math.Sqrt(normRecognized))
}

// filterRecognizedTools разделяет инструменты на категории
func filterRecognizedTools(req *FilterReq) (*FilterRes, error) {
	accessTools := make([]*domain.RecognizedTool, 0, len(req.Tools))
	manualCheckTools := make([]*domain.RecognizedTool, 0, len(req.Tools))
	unknownTools := make([]*domain.RecognizedTool, 0, len(req.Tools))
	missingTools := make([]*ToolTypeDTO, 0)

	refMap := make(map[int64]*domain.ToolType)
	for _, r := range req.ReferenceTools {
		refMap[r.Id] = r
	}

	recognizedMap := make(map[int64]*domain.RecognizedTool)
	for _, recognized := range req.Tools {
		ref, exists := refMap[recognized.ToolTypeId]
		if !exists {
			unknownTools = append(unknownTools, recognized)
			continue
		}

		cosSim := cosineSimilarity(ref.ReferenceEmbedding, recognized.Embedding)
		// fmt.Printf("DEBUG: toolId: %d, Confidence: %f, cosSim: %f\n", recognized.ToolTypeId, recognized.Confidence, cosSim)
		if cosSim >= req.CosineSimCompare && recognized.Confidence >= req.ConfidenceCompare {
			accessTools = append(accessTools, recognized)
		} else {
			manualCheckTools = append(manualCheckTools, recognized)
		}

		recognizedMap[recognized.ToolTypeId] = recognized
	}

	for _, ref := range req.ReferenceTools {
		if _, ok := recognizedMap[ref.Id]; !ok {
			missingTools = append(missingTools, ToToolTypeDTO(ref))
		}
	}

	return NewFilterRes(accessTools, manualCheckTools, unknownTools, missingTools), nil
}
