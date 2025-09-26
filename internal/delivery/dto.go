package delivery

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/usecase"
)

type CheckReq struct {
	EmployeeId string   `json:"employee_id" binding:"required"`
	Image      ImageReq `json:"image" binding:"required"`
}

type ImageReq struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	Data        []byte `json:"data" binding:"required"`
}

type CheckRes struct {
	ImageUrl         string               `json:"image_url"`
	AccessTools      []*RecognizedToolDTO `json:"access_tools"`
	ManualCheckTools []*RecognizedToolDTO `json:"manual_check_tools"`
	UnknownTools     []*RecognizedToolDTO `json:"unknown_tools"`
	MissingTools     []*ToolTypeDTO       `json:"missing_tools"`
}

type RecognizedToolDTO struct {
	ToolTypeId int64   `json:"tool_type_id"`
	Confidence float32 `json:"confidence"`
}

type ToolTypeDTO struct {
	Id         int64  `json:"id"`
	PartNumber string `json:"part_number"`
	Name       string `json:"name"`
}

func toDeliveryRecognizedToolDTO(tool *domain.RecognizedTool) *RecognizedToolDTO {
	return &RecognizedToolDTO{
		ToolTypeId: tool.ToolTypeId,
		Confidence: tool.Confidence,
	}
}

func toDeliveryToolTypeDTO(dto *usecase.ToolTypeDTO) *ToolTypeDTO {
	return &ToolTypeDTO{
		Id:         dto.Id,
		PartNumber: dto.PartNumber,
		Name:       dto.Name,
	}
}

func toArrDeliveryToolTypeDTO(tools []*usecase.ToolTypeDTO) []*ToolTypeDTO {
	result := make([]*ToolTypeDTO, len(tools))
	for i, tool := range tools {
		result[i] = toDeliveryToolTypeDTO(tool)
	}

	return result
}

func toArrDeliveryRecognizedToolDTO(tools []*domain.RecognizedTool) []*RecognizedToolDTO {
	result := make([]*RecognizedToolDTO, len(tools))
	for i, tool := range tools {
		result[i] = toDeliveryRecognizedToolDTO(tool)
	}

	return result
}

func ToDeliveryCheckRes(res *usecase.CheckRes) *CheckRes {
	return &CheckRes{
		ImageUrl:         res.ImageUrl,
		AccessTools:      toArrDeliveryRecognizedToolDTO(res.AccessTools),
		ManualCheckTools: toArrDeliveryRecognizedToolDTO(res.ManualCheckTools),
		UnknownTools:     toArrDeliveryRecognizedToolDTO(res.UnknownTools),
		MissingTools:     toArrDeliveryToolTypeDTO(res.MissingTools),
	}
}

func ToUseCaseCheckReq(req *CheckReq) *usecase.CheckReq {
	return &usecase.CheckReq{
		EmployeeId: req.EmployeeId,
		Image:      *ToUsecaseImageReq(&req.Image),
	}
}

func ToUsecaseImageReq(req *ImageReq) *usecase.ImageReq {
	return &usecase.ImageReq{
		Filename:    req.Filename,
		ContentType: req.ContentType,
		Data:        req.Data,
	}
}
