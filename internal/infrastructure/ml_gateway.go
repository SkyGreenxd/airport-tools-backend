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
	"strconv"
)

type MlGateway struct {
	client  *http.Client
	baseUrl string
}

func NewMlGateway(client *http.Client, baseUrl string) *MlGateway {
	return &MlGateway{
		client:  client,
		baseUrl: baseUrl,
	}
}

type mlAPIResponse struct {
	ImageId     string `json:"image_id"`
	Instruments []struct {
		Bbox       []int     `json:"bbox"`
		ToolTypeId int64     `json:"class"`
		Confidence float32   `json:"confidence"`
		Embedding  []float32 `json:"embedding"`
		Hash       int       `json:"hash"`
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
	// fmt.Println("Debug: " + getUrl)
	res, err := ml.client.Get(getUrl)
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
	for _, instrument := range apiResp.Instruments {
		recognizedTool := domain.NewRecognizedTool(instrument.ToolTypeId+1, instrument.Confidence, strconv.Itoa(instrument.Hash), instrument.Embedding)
		scanResult.Tools = append(scanResult.Tools, recognizedTool)
	}

	return &scanResult, nil
}
