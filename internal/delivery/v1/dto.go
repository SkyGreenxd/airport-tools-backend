package v1

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/usecase"
	"time"
)

type GetRolesRes struct {
	Roles []domain.Role `json:"roles"`
}

func toDeliveryGetRolesRes(roles *usecase.GetRolesRes) GetRolesRes {
	return GetRolesRes{
		roles.Roles,
	}
}

type VerificationReq struct {
	TransactionID int64  `json:"transaction_id" binding:"required"`
	QAEmployeeId  string `json:"qa_employee_id" binding:"required"`
	Notes         string `json:"notes"`
}

type VerificationRes struct {
	TransactionID string    `json:"transaction_id"` // ID транзакции, которую QA завершил
	Status        string    `json:"status"`         // Новый статус
	VerifiedBy    string    `json:"verified_by"`    // Табельный номер или имя QA
	VerifiedAt    time.Time `json:"verified_at"`    // Время завершения проверки
	Message       string    `json:"message"`        // Краткое текстовое подтверждение
}

type ListTransactionsRes struct {
	Transactions []TransactionDTO `json:"transactions"`
}

type TransactionDTO struct {
	Id        int64     `json:"id"`
	ToolSetId int64     `json:"tool_set_id"`
	CreatedAt time.Time `json:"created_at"`
	User      UserDto   `json:"user"`
}

type UserDto struct {
	FullName   string `json:"full_name"`
	EmployeeId string `json:"employee_id"`
}

type RegisterReq struct {
	EmployeeId string `json:"employee_id" binding:"required"`
	FullName   string `json:"full_name" binding:"required"`
	Role       string `json:"role" binding:"required"`
}

type RegisterRes struct {
	Id string `json:"id"`
}

type LoginReq struct {
	EmployeeId string `json:"employee_id" binding:"required"`
}

type LoginRes struct {
	Role domain.Role `json:"role"`
}

type CheckReq struct {
	EmployeeId string `json:"employee_id" binding:"required"`
	Data       string `json:"data" binding:"required"`
	ToolSetId  int64  `json:"tool_set_id"`
}

type CheckRes struct {
	ImageUrl         string               `json:"image_url"`
	DebugImageUrl    string               `json:"debug_image_url"`
	AccessTools      []*RecognizedToolDTO `json:"access_tools"`
	ManualCheckTools []*RecognizedToolDTO `json:"manual_check_tools"`
	UnknownTools     []*RecognizedToolDTO `json:"unknown_tools"`
	MissingTools     []*ToolTypeDTO       `json:"missing_tools"`
	TransactionType  string               `json:"transaction_type"`
}

type RecognizedToolDTO struct {
	ToolTypeId int64     `json:"tool_type_id"`
	Confidence float32   `json:"confidence"`
	Bbox       []float32 `json:"bbox"`
}

type ToolTypeDTO struct {
	Id         int64  `json:"id"`
	PartNumber string `json:"part_number"`
	Name       string `json:"name"`
}

func NewListTransactionsRes(transactions []TransactionDTO) *ListTransactionsRes {
	return &ListTransactionsRes{
		Transactions: transactions,
	}
}

func toDeliveryListTransactionsRes(transactions []*usecase.TransactionDTO) *ListTransactionsRes {
	res := make([]TransactionDTO, len(transactions))
	for i, transaction := range transactions {
		res[i] = *toDeliveryTransactionDTO(transaction)
	}

	return NewListTransactionsRes(res)
}

func toDeliveryTransactionDTO(transaction *usecase.TransactionDTO) *TransactionDTO {
	return &TransactionDTO{
		Id:        transaction.Id,
		ToolSetId: transaction.ToolSetId,
		CreatedAt: transaction.CreatedAt,
		User:      toDeliveryUserDto(transaction.User),
	}
}

func toDeliveryUserDto(user usecase.UserDto) UserDto {
	return UserDto{
		FullName:   user.FullName,
		EmployeeId: user.EmployeeId,
	}
}

func toDeliveryRecognizedToolDTO(tool *domain.RecognizedTool) *RecognizedToolDTO {
	return &RecognizedToolDTO{
		ToolTypeId: tool.ToolTypeId,
		Confidence: tool.Confidence,
		Bbox:       tool.Bbox,
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
		DebugImageUrl:    res.DebugImageUrl,
		AccessTools:      toArrDeliveryRecognizedToolDTO(res.AccessTools),
		ManualCheckTools: toArrDeliveryRecognizedToolDTO(res.ManualCheckTools),
		UnknownTools:     toArrDeliveryRecognizedToolDTO(res.UnknownTools),
		MissingTools:     toArrDeliveryToolTypeDTO(res.MissingTools),
		TransactionType:  res.TransactionType,
	}
}

func ToUseCaseCheckReq(req *CheckReq) *usecase.CheckReq {
	return &usecase.CheckReq{
		EmployeeId: req.EmployeeId,
		Data:       req.Data,
		ToolSetId:  req.ToolSetId,
	}
}

func toUseCaseLoginReq(req LoginReq) *usecase.LoginReq {
	return &usecase.LoginReq{
		EmployeeId: req.EmployeeId,
	}
}

func toDeliveryLoginRes(res *usecase.LoginRes) LoginRes {
	return LoginRes{
		Role: res.Role,
	}
}
