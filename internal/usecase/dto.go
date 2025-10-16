package usecase

import "airport-tools-backend/internal/domain"

type TransactionProcess struct {
	UserId    int64
	Data      string
	ToolSetId int64
}

// CheckReq представляет запрос на выдачу/сдачу инструментов
type CheckReq struct {
	EmployeeId string
	Data       string
	ToolSetId  int64
}

// CheckRes содержит результат проверки инструментов после сканирования.
type CheckRes struct {
	ImageUrl         string
	DebugImageUrl    string
	AccessTools      []*domain.RecognizedTool
	ManualCheckTools []*domain.RecognizedTool
	UnknownTools     []*domain.RecognizedTool
	MissingTools     []*ToolTypeDTO
}

type UploadImageReq struct {
	Data string
	Mode string
}

type ToolTypeDTO struct {
	Id         int64
	PartNumber string
	Name       string
}

// ScanRequest используется для передачи изображения в ML-сервис.
type ScanRequest struct {
	ImageId   string
	ImageUrl  string
	Threshold float32
}

// ScanResult возвращает распознанные инструменты ML-сервиса.
type ScanResult struct {
	Tools         []*domain.RecognizedTool
	DebugImageUrl string
}

type CreateScanReq struct {
	TransactionId int64
	ScanType      domain.ScanType
	ImageUrl      string
	DebugImageUrl string
	Tools         []*domain.RecognizedTool
}

type FilterReq struct {
	ConfidenceCompare float32
	CosineSimCompare  float64
	Tools             []*domain.RecognizedTool
	ReferenceTools    []*domain.ToolType
}

type FilterRes struct {
	AccessTools      []*domain.RecognizedTool
	ManualCheckTools []*domain.RecognizedTool
	UnknownTools     []*domain.RecognizedTool
	MissingTools     []*ToolTypeDTO
}

type UploadImageRes struct {
	Key      string
	ImageUrl string
}

func ToToolTypeDTO(tool *domain.ToolType) *ToolTypeDTO {
	return &ToolTypeDTO{
		Id:         tool.Id,
		PartNumber: tool.PartNumber,
		Name:       tool.Name,
	}
}

func NewUploadImageRes(key string, imageUrl string) *UploadImageRes {
	return &UploadImageRes{
		Key:      key,
		ImageUrl: imageUrl,
	}
}

func NewCheckinRes(imageUrl, debugImageUrl string, accessTools, manualCheckTools, unknownTools []*domain.RecognizedTool, missingTools []*ToolTypeDTO) *CheckRes {
	return &CheckRes{
		ImageUrl:         imageUrl,
		DebugImageUrl:    debugImageUrl,
		AccessTools:      accessTools,
		ManualCheckTools: manualCheckTools,
		UnknownTools:     unknownTools,
		MissingTools:     missingTools,
	}
}

func NewFilterRes(accessTools, manualCheckTools, unknownTools []*domain.RecognizedTool, missingTools []*ToolTypeDTO) *FilterRes {
	return &FilterRes{
		AccessTools:      accessTools,
		ManualCheckTools: manualCheckTools,
		UnknownTools:     unknownTools,
		MissingTools:     missingTools,
	}
}

func NewFilterReq(confidenceCompare float32, cosineSimCompare float64, Tools []*domain.RecognizedTool, referenceTools []*domain.ToolType) *FilterReq {
	return &FilterReq{
		ConfidenceCompare: confidenceCompare,
		CosineSimCompare:  cosineSimCompare,
		Tools:             Tools,
		ReferenceTools:    referenceTools,
	}
}

func NewCreateScanReq(transactionId int64, scanType domain.ScanType, imageUrl, debugImageUrl string, tools []*domain.RecognizedTool) *CreateScanReq {
	return &CreateScanReq{
		TransactionId: transactionId,
		ScanType:      scanType,
		ImageUrl:      imageUrl,
		DebugImageUrl: debugImageUrl,
		Tools:         tools,
	}
}

func NewScanReq(imageId, imageUrl string, threshold float32) *ScanRequest {
	return &ScanRequest{
		ImageId:   imageId,
		ImageUrl:  imageUrl,
		Threshold: threshold,
	}
}

func NewUploadImageReq(data, mode string) *UploadImageReq {
	return &UploadImageReq{
		Data: data,
		Mode: mode,
	}
}

func NewTransactionProcess(userId int64, data string, toolSetId int64) *TransactionProcess {
	return &TransactionProcess{
		UserId:    userId,
		Data:      data,
		ToolSetId: toolSetId,
	}
}
