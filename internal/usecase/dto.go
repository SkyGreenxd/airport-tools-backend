package usecase

import (
	"airport-tools-backend/internal/domain"
	"time"
)

type MlErrorTransaction struct {
	TransactionID  int64
	SourceImageUrl string
	DebugImageUrl  string
}

func NewMlErrorTransaction(id int64, sUrl, dUrl string) *MlErrorTransaction {
	return &MlErrorTransaction{
		TransactionID:  id,
		SourceImageUrl: sUrl,
		DebugImageUrl:  dUrl,
	}
}

type GetAvgWorkDurationRes struct {
	Transactions []GetAvgWorkDuration
}

type GetAvgWorkDuration struct {
	User            UserDto
	AvgWorkDuration float64
}

func NewGetAvgWorkDuration(user UserDto, avgWorkDuration float64) GetAvgWorkDuration {
	return GetAvgWorkDuration{
		User:            user,
		AvgWorkDuration: avgWorkDuration,
	}
}
func NewGetAvgWorkDurationRes(transactions []GetAvgWorkDuration) *GetAvgWorkDurationRes {
	return &GetAvgWorkDurationRes{
		Transactions: transactions,
	}
}

type GetAllWorkDurationRes struct {
	Transactions []GetWorkDuration
}

type GetWorkDuration struct {
	TransactionId int64
	WorkDuration  float64
}

func NewGetAllWorkDurationRes(arr []GetWorkDuration) *GetAllWorkDurationRes {
	return &GetAllWorkDurationRes{
		Transactions: arr,
	}
}

func NewGetWorkDuration(id int64, workDuration float64) GetWorkDuration {
	return GetWorkDuration{
		TransactionId: id,
		WorkDuration:  workDuration,
	}
}

// ListTransactionsRes список транзакций
type ListTransactionsRes struct {
	Transactions []*TransactionDTO
}

type GetTransactionStatisticsRes struct {
	Transactions       int
	OpenedTransactions int
	ClosedTransactions int
	QATransactions     int
	FailedTransactions int
}

func NewGetTransactionStatisticsRes(transactions, opened, closed, qa, failed int) *GetTransactionStatisticsRes {
	return &GetTransactionStatisticsRes{
		Transactions:       transactions,
		OpenedTransactions: opened,
		ClosedTransactions: closed,
		QATransactions:     qa,
		FailedTransactions: failed,
	}
}

type HumanErrorStats struct {
	FullName    string
	EmployeeId  string
	QAHitsCount int64
}

func NewHumanErrorStats(fullName, employeeId string, QAHitsCount int64) HumanErrorStats {
	return HumanErrorStats{
		FullName:    fullName,
		EmployeeId:  employeeId,
		QAHitsCount: QAHitsCount,
	}
}

type ModelOrHumanStatsRes struct {
	MlErrors    int
	HumanErrors int
}

func NewModelOrHumanStatsRes(ml int, humans int) *ModelOrHumanStatsRes {
	return &ModelOrHumanStatsRes{
		MlErrors:    ml,
		HumanErrors: humans,
	}
}

type QaTransactionsRes struct {
	Qa           UserDto
	Transactions []*TransactionResolutionDTO
}

type TransactionResolutionDTO struct {
	Transaction *TransactionDTO
	Reason      domain.Reason
	Notes       string
	CreatedAt   time.Time
}

func ToListTransactionResolutionDTO(transactions []*domain.TransactionResolution) []*TransactionResolutionDTO {
	result := make([]*TransactionResolutionDTO, len(transactions))
	for i, tr := range transactions {
		result[i] = ToTransactionResolutionDTO(tr.Transaction, tr.Reason, tr.Notes, tr.CreatedAt)
	}

	return result
}

func ToTransactionResolutionDTO(transaction *domain.Transaction, reason domain.Reason, notes string, createdAt time.Time) *TransactionResolutionDTO {
	var transactionDTO *TransactionDTO
	if transaction != nil {
		transactionDTO = toTransactionDTO(transaction)
	}

	return &TransactionResolutionDTO{
		Transaction: transactionDTO,
		Reason:      reason,
		Notes:       notes,
		CreatedAt:   createdAt,
	}
}

func NewQaTransactionsRes(qa UserDto, transactions []*TransactionResolutionDTO) *QaTransactionsRes {
	return &QaTransactionsRes{
		Qa:           qa,
		Transactions: transactions,
	}
}

type UserTransactionsReq struct {
	EmployeeId string
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      *int
}

func NewUserTransactionsReq(employeeId string, startDate, endDate *time.Time, limit *int) *UserTransactionsReq {
	return &UserTransactionsReq{
		EmployeeId: employeeId,
		StartDate:  startDate,
		EndDate:    endDate,
		Limit:      limit,
	}
}

type GetQAVerificationRes struct {
	TransactionId    int64
	ToolSetId        int64
	CreatedAt        time.Time
	User             UserDto
	AccessTools      []*domain.RecognizedTool
	ProblematicTools *ProblematicTools
	ImageUrl         string
	Status           string
}

type ProblematicTools struct {
	ManualCheckTools []*domain.RecognizedTool
	UnknownTools     []*domain.RecognizedTool
	MissingTools     []*ToolTypeDTO
}

func NewGetQAVerificationRes(id, toolSetId int64, createdAt time.Time, user UserDto, accessTools []*domain.RecognizedTool, problematicTools *ProblematicTools, imageurl, status string) *GetQAVerificationRes {
	return &GetQAVerificationRes{
		TransactionId:    id,
		ToolSetId:        toolSetId,
		CreatedAt:        createdAt,
		User:             user,
		AccessTools:      accessTools,
		ProblematicTools: problematicTools,
		ImageUrl:         imageurl,
		Status:           status,
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
	Reason        domain.Reason
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
	Status    domain.Status
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
		AccessTools:      accessTools,
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
		Status:    transaction.Status,
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

func NewVerification(transactionID int64, qaEmployeeId string, reason domain.Reason, notes string) *Verification {
	return &Verification{
		TransactionID: transactionID,
		QAEmployeeId:  qaEmployeeId,
		Reason:        reason,
		Notes:         notes,
	}
}

func NewUserDto(fullname, employeeId string) UserDto {
	return UserDto{
		FullName:   fullname,
		EmployeeId: employeeId,
	}
}
