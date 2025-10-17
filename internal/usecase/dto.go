package usecase

import (
	"airport-tools-backend/internal/domain"
	"time"
)

// ListTransactionsRes список транзакций
type ListTransactionsRes struct {
	Transactions []*TransactionDTO
}

type GetQAVerificationRes struct {
	TransactionId    int64
	ToolSetId        int64
	CreatedAt        time.Time
	User             UserDto
	AccessTools      []*domain.RecognizedTool
	ProblematicTools *ProblematicTools
	ImageUrl         string
}

type ProblematicTools struct {
	ManualCheckTools []*domain.RecognizedTool
	UnknownTools     []*domain.RecognizedTool
	MissingTools     []*ToolTypeDTO
}

func NewGetQAVerificationRes(id, toolSetId int64, createdAt time.Time, user UserDto, accessTools []*domain.RecognizedTool, problematicTools *ProblematicTools, imageurl string) *GetQAVerificationRes {
	return &GetQAVerificationRes{
		TransactionId:    id,
		ToolSetId:        toolSetId,
		CreatedAt:        createdAt,
		User:             user,
		AccessTools:      accessTools,
		ProblematicTools: problematicTools,
		ImageUrl:         imageurl,
	}
}

func NewProblematicTools(manualCheckTools, unknownTools []*domain.RecognizedTool, missingTools []*ToolTypeDTO) *ProblematicTools {
	return &ProblematicTools{
		ManualCheckTools: manualCheckTools,
		UnknownTools:     unknownTools,
		MissingTools:     missingTools,
	}
}

type Verification struct {
	TransactionID int64
	QAEmployeeId  string
	Notes         string
}

type VerificationRes struct {
	TransactionID int64
	Status        string
	VerifiedBy    string
	CreatedAt     time.Time
}

func NewVerificationRes(id int64, status string, verifiedBy string, createdAt time.Time) *VerificationRes {
	return &VerificationRes{
		TransactionID: id,
		Status:        status,
		VerifiedBy:    verifiedBy,
		CreatedAt:     createdAt,
	}
}

type RegisterReq struct {
	EmployeeId string
	FullName   string
	Role       domain.Role
}

type RegisterRes struct {
	Id int64
}

type TransactionDTO struct {
	Id        int64
	ToolSetId int64
	CreatedAt time.Time
	User      UserDto
}

type UserDto struct {
	FullName   string
	EmployeeId string
}

type LoginReq struct {
	EmployeeId string
}

type LoginRes struct {
	Role domain.Role
}

type GetRolesRes struct {
	Roles []domain.Role
}

// TransactionProcess внутренняя структура для сдачи/выдачи инструментов
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
	ProblematicTools *ProblematicTools
	TransactionType  string
	Status           string
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
	CosineSimCompare  float32
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

func NewCheckinRes(imageUrl, debugImageUrl string, accessTools, manualCheckTools, unknownTools []*domain.RecognizedTool, missingTools []*ToolTypeDTO, transactionType, status string) *CheckRes {
	return &CheckRes{
		ImageUrl:         imageUrl,
		DebugImageUrl:    debugImageUrl,
		ProblematicTools: NewProblematicTools(manualCheckTools, unknownTools, missingTools),
		TransactionType:  transactionType,
		Status:           status,
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

func NewFilterReq(confidenceCompare, cosineSimCompare float32, Tools []*domain.RecognizedTool, referenceTools []*domain.ToolType) *FilterReq {
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

func NewGetRolesRes(roles []domain.Role) *GetRolesRes {
	return &GetRolesRes{
		Roles: roles,
	}
}

func NewLoginRes(role domain.Role) *LoginRes {
	return &LoginRes{
		Role: role,
	}
}

func NewListTransactionsRes(tools []*TransactionDTO) *ListTransactionsRes {
	return &ListTransactionsRes{
		Transactions: tools,
	}
}

func toListTransactionsRes(tools []*domain.Transaction) []*TransactionDTO {
	result := make([]*TransactionDTO, len(tools))
	for i, tool := range tools {
		result[i] = toTransactionDTO(tool)
	}

	return result
}

func toTransactionDTO(transaction *domain.Transaction) *TransactionDTO {
	var userDto UserDto
	if transaction.User != nil {
		userDto = toUserDTO(*transaction.User)
	}

	return &TransactionDTO{
		Id:        transaction.Id,
		ToolSetId: transaction.ToolSetId,
		CreatedAt: transaction.CreatedAt,
		User:      userDto,
	}
}

func toUserDTO(user domain.User) UserDto {
	return UserDto{
		FullName:   user.FullName,
		EmployeeId: user.EmployeeId,
	}
}

func NewRegisterRes(id int64) *RegisterRes {
	return &RegisterRes{
		Id: id,
	}
}

func NewVerification(transactionID int64, qaEmployeeId string, notes string) *Verification {
	return &Verification{
		TransactionID: transactionID,
		QAEmployeeId:  qaEmployeeId,
		Notes:         notes,
	}
}

func NewUserDto(fullname, employeeId string) UserDto {
	return UserDto{
		FullName:   fullname,
		EmployeeId: employeeId,
	}
}
