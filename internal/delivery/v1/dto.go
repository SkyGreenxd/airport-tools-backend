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
	QAEmployeeId string        `json:"qa_employee_id" binding:"required"`
	Reason       domain.Reason `json:"reason" binding:"required"`
	Notes        string        `json:"notes"`
}

type VerificationRes struct {
	TransactionID int64     `json:"transaction_id"`
	Status        string    `json:"status"`
	VerifiedBy    string    `json:"verified_by"`
	CreatedAt     time.Time `json:"created_at"`
}

func toDeliveryVerificationRes(res *usecase.VerificationRes) *VerificationRes {
	return &VerificationRes{
		TransactionID: res.TransactionID,
		Status:        res.Status,
		VerifiedBy:    res.VerifiedBy,
		CreatedAt:     res.CreatedAt,
	}
}

type GetQAVerificationRes struct {
	TransactionId    int64                `json:"transaction_id"`
	ToolSetId        int64                `json:"tool_set_id"`
	CreatedAt        time.Time            `json:"created_at"`
	User             UserDto              `json:"user"`
	AccessTools      []*RecognizedToolDTO `json:"access_tools"`
	ProblematicTools *ProblematicTools    `json:"problematic_tools"`
	ImageUrl         string               `json:"image_url"`
}

type ProblematicTools struct {
	ManualCheckTools []*RecognizedToolDTO `json:"manual_check_tools"`
	UnknownTools     []*RecognizedToolDTO `json:"unknown_tools"`
	MissingTools     []*ToolTypeDTO       `json:"missing_tools"`
}

func toDeliveryGetQAVerificationRes(res *usecase.GetQAVerificationRes) *GetQAVerificationRes {
	return &GetQAVerificationRes{
		TransactionId:    res.TransactionId,
		ToolSetId:        res.ToolSetId,
		CreatedAt:        res.CreatedAt,
		User:             toDeliveryUserDto(res.User),
		AccessTools:      toArrDeliveryRecognizedToolDTO(res.AccessTools),
		ProblematicTools: toDeliveryProblematicTools(res.ProblematicTools),
		ImageUrl:         res.ImageUrl,
	}
}

func toDeliveryProblematicTools(tools *usecase.ProblematicTools) *ProblematicTools {
	return &ProblematicTools{
		ManualCheckTools: toArrDeliveryRecognizedToolDTO(tools.ManualCheckTools),
		UnknownTools:     toArrDeliveryRecognizedToolDTO(tools.UnknownTools),
		MissingTools:     toArrDeliveryToolTypeDTO(tools.MissingTools),
	}
}

type ListTransactionsRes struct {
	Transactions []TransactionDTO `json:"transactions"`
}

type TransactionDTO struct {
	Id        int64         `json:"id"`
	ToolSetId int64         `json:"tool_set_id"`
	CreatedAt time.Time     `json:"created_at"`
	User      UserDto       `json:"user"`
	Status    domain.Status `json:"status"`
}

type UserDto struct {
	FullName   string `json:"full_name"`
	EmployeeId string `json:"employee_id"`
}

type RegisterReq struct {
	EmployeeId string      `json:"employee_id" binding:"required"`
	FullName   string      `json:"full_name" binding:"required"`
	Role       domain.Role `json:"role" binding:"required"`
}

type RegisterRes struct {
	Id int64 `json:"id"`
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
	ProblematicTools *ProblematicTools    `json:"problematic_tools"`
	TransactionType  string               `json:"transaction_type"`
	Status           string               `json:"status"`
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
		Status:    transaction.Status,
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
		ProblematicTools: toDeliveryProblematicTools(res.ProblematicTools),
		TransactionType:  res.TransactionType,
		Status:           res.Status,
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

func toUseCaseRegisterReq(req RegisterReq) *usecase.RegisterReq {
	return &usecase.RegisterReq{
		EmployeeId: req.EmployeeId,
		FullName:   req.FullName,
		Role:       req.Role,
	}
}

func toDeliveryRegisterRes(res *usecase.RegisterRes) RegisterRes {
	return RegisterRes{
		Id: res.Id,
	}
}
