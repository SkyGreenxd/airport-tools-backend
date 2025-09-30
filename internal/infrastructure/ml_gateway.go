package infrastructure

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/usecase"
	"airport-tools-backend/pkg/e"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// MlGateway клиент для взаимодействия с ML-сервисом распознавания инструментов
type MlGateway struct {
	client  *http.Client
	baseUrl string
	s3      usecase.ImageStorage
}

func NewMlGateway(client *http.Client, baseUrl string, s3 usecase.ImageStorage) *MlGateway {
	return &MlGateway{
		client:  client,
		baseUrl: baseUrl,
		s3:      s3,
	}
}

// mlAPIResponse структура для декодирования ответа ML-сервиса
type mlAPIResponse struct {
	ImageId     string `json:"image_id"`
	Instruments []struct {
		Bbox       []float32 `json:"bbox"`
		ToolTypeId int64     `json:"class"`
		Confidence float32   `json:"confidence"`
		Embedding  []float32 `json:"embedding"`
	} `json:"instruments"`
	DebugImage string `json:"debug_image"`
}

// ScanTools отправляет изображение на ML-сервис и возвращает распознанные инструменты
func (ml *MlGateway) ScanTools(ctx context.Context, req *usecase.ScanRequest) (*usecase.ScanResult, error) {
	const op = "MlGateway.ScanTools"

	getUrl := fmt.Sprintf("%s/predict/?image_id=%s&url=%s&thresh=%s",
		ml.baseUrl,
		url.QueryEscape(req.ImageId),
		url.QueryEscape(req.ImageUrl),
		url.QueryEscape(fmt.Sprintf("%f", req.Threshold)),
	)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, getUrl, nil)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	res, err := ml.client.Do(httpReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, e.Wrap(op, fmt.Errorf("%s: %d", e.ErrMLServiceNonOK, res.StatusCode))
	}

	var apiResp mlAPIResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		return nil, e.Wrap(op, fmt.Errorf("%s: %w", e.ErrMLServiceDecode, err))
	}

	var scanResult usecase.ScanResult
	uploadImageRes, err := ml.s3.UploadImage(ctx, apiResp.DebugImage)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	scanResult.DebugImageUrl = uploadImageRes.ImageUrl
	for _, instrument := range apiResp.Instruments {
		recognizedTool := domain.NewRecognizedTool(instrument.ToolTypeId+1, instrument.Confidence, instrument.Embedding)
		scanResult.Tools = append(scanResult.Tools, recognizedTool)
	}

	return &scanResult, nil
}
