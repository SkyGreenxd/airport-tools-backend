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

type MlGateway struct {
	client  *http.Client
	baseUrl string
}

func NewMlGateway(client *http.Client) *MlGateway {
	return &MlGateway{}
}

type mlAPIResponse struct {
	ImageId     string `json:"image_id"`
	Instruments []struct {
		Bbox       []int     `json:"bbox"`
		ToolTypeId int64     `json:"class"`
		Confidence float32   `json:"confidence"`
		Embedding  []float32 `json:"embedding"`
		Hash       string    `json:"hash"`
	} `json:"instruments"`
	ImageUrl string `json:"debug_image_url"`
}

// http://127.0.0.1:8000/analyze/?image_id={id}&url={url}
func (ml *MlGateway) ScanTools(ctx context.Context, req *usecase.ScanRequest) (*usecase.ScanResult, error) {
	const op = "MlGateway.ScanTools"

	getUrl := fmt.Sprintf("%s/analyze/?image_id=%s&url=%s",
		ml.baseUrl,
		url.QueryEscape(req.ImageId),
		url.QueryEscape(req.ImageUrl),
	)
	res, err := ml.client.Get(getUrl)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, e.Wrap(op, fmt.Errorf("ml service returned non-200 status: %d", res.StatusCode))
	}

	var apiResp mlAPIResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode ml service response: %w", err)
	}

	var scanResult usecase.ScanResult
	for _, instrument := range apiResp.Instruments {
		recognizedTool := domain.NewRecognizedTool(instrument.ToolTypeId, instrument.Confidence, instrument.Hash, instrument.Embedding)
		scanResult.Tools = append(scanResult.Tools, recognizedTool)
	}

	return &scanResult, nil
}
