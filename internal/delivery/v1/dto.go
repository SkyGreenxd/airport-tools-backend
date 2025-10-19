package v1

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/repository"
	"airport-tools-backend/internal/usecase"
	"time"
)

type ToolWithErrorCount struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	MLErrorCount int64  `json:"ml_error_count"`
}

type ToolSetWithErrors struct {
	ID    int64                `json:"id"`
	Name  string               `json:"name"`
	Tools []ToolWithErrorCount `json:"tools"`
}

func toArrDeliveryToolSetWithErrors(res []*repository.ToolSetWithErrors) []ToolSetWithErrors {
	resArr := make([]ToolSetWithErrors, len(res))
	for i, r := range res {
		resArr[i] = toDeliveryToolSetWithErrors(*r)
	}

	return resArr
}

func toDeliveryToolSetWithErrors(res repository.ToolSetWithErrors) ToolSetWithErrors {
	return ToolSetWithErrors{
		ID:    res.ID,
		Name:  res.Name,
		Tools: toArrDeliveryToolWithErrorCount(res.Tools),
	}
}

func toArrDeliveryToolWithErrorCount(arr []repository.ToolWithErrorCount) []ToolWithErrorCount {
	result := make([]ToolWithErrorCount, len(arr))
	for i, t := range arr {
		result[i] = toDeliveryToolWithErrorCount(t)
	}
	return result
}

func toDeliveryToolWithErrorCount(res repository.ToolWithErrorCount) ToolWithErrorCount {
	return ToolWithErrorCount{
		ID:           res.ID,
		Name:         res.Name,
		MLErrorCount: res.MLErrorCount,
	}
}

type AddToolSetRes struct {
	Id    int64          `json:"id"`
	Name  string         `json:"name"`
	Tools []*ToolTypeDTO `json:"tools"`
}

type AddToolSetReq struct {
	ToolSetName string  `json:"tool_set_name" binding:"required,min=3"`
	ToolsIds    []int64 `json:"tools_ids" binding:"required,min=1,dive,gt=0"`
}

type GetUsersListTransactionsRes struct {
	Transactions []*TransactionDTO
	Avg          float64
}

type GetAllTransactions struct {
	User         UserDto            `json:"user"`
	Transactions []LightTransaction `json:"transactions"`
}

type LightTransaction struct {
	Id        int64
	CreatedAt time.Time
}

type MlErrorTransaction struct {
	TransactionID  int64  `json:"transaction_id"`
	SourceImageUrl string `json:"source_image_url"`
	DebugImageUrl  string `json:"debug_image_url"`
}

type GetAvgWorkDurationRes struct {
	Transactions []GetAvgWorkDuration `json:"transactions"`
}

type GetAvgWorkDuration struct {
	User            UserDto `json:"user"`
	AvgWorkDuration float64 `json:"avg_work_duration"`
}

type HumanErrorStats struct {
	FullName    string `json:"full_name"`
	EmployeeId  string `json:"employee_id"`
	QAHitsCount int64  `json:"qa_hits_count"`
}

type GetTransactionStatisticsRes struct {
	Transactions       int `json:"transactions"`
	OpenedTransactions int `json:"opened_transactions"`
	ClosedTransactions int `json:"closed_transactions"`
	QATransactions     int `json:"qa_transactions"`
	FailedTransactions int `json:"failed_transactions"`
}

type ModelOrHumanStatsRes struct {
	MlErrors    int `json:"ml_errors"`
	HumanErrors int `json:"human_errors"`
}

type QaTransactionsRes struct {
	Qa           UserDto                     `json:"qa"`
	Transactions []*TransactionResolutionDTO `json:"transactions"`
}

type TransactionResolutionDTO struct {
	Transaction *TransactionDTO `json:"transaction"`
	Reason      domain.Reason   `json:"reason"`
	Notes       string          `json:"notes"`
	CreatedAt   time.Time       `json:"created_at"`
}

type StatisticsRes struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type GetRolesRes struct {
	Roles []string `json:"roles"`
}

type VerificationReq struct {
	QAEmployeeId string        `json:"qa_employee_id" binding:"required"`
	Reason       domain.Reason `json:"reason" binding:"required"`
	Notes        string        `json:"notes"`
	ToolIds      []int64       `json:"tool_ids" binding:"omitempty,dive,gt=0"`
}

type VerificationRes struct {
	TransactionID int64     `json:"transaction_id"`
	Status        string    `json:"status"`
	VerifiedBy    string    `json:"verified_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetQAVerificationRes struct {
	TransactionId    int64                `json:"transaction_id"`
	ToolSetId        int64                `json:"tool_set_id"`
	CreatedAt        time.Time            `json:"created_at"`
	User             UserDto              `json:"user"`
	AccessTools      []*RecognizedToolDTO `json:"access_tools"`
	ProblematicTools *ProblematicTools    `json:"problematic_tools"`
	ImageUrl         string               `json:"image_url"`
	Status           string               `json:"status"`
}

type ProblematicTools struct {
	ManualCheckTools []*RecognizedToolDTO `json:"manual_check_tools"`
	UnknownTools     []*RecognizedToolDTO `json:"unknown_tools"`
	MissingTools     []*ToolTypeDTO       `json:"missing_tools"`
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
	EmployeeId string `json:"employee_id" binding:"required"`
	FullName   string `json:"full_name" binding:"required"`
	Role       string `json:"role" binding:"required"`
}

type RegisterRes struct {
	Id int64 `json:"id"`
}

type LoginReq struct {
	EmployeeId string `json:"employee_id" binding:"required"`
}

type LoginRes struct {
	Role string `json:"role"`
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

func toArrDeliveryUserDto(users []usecase.UserDto) []UserDto {
	res := make([]UserDto, len(users))
	for i, user := range users {
		res[i] = toDeliveryUserDto(user)
	}

	return res
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

func toDeliveryGetUsersListTransactionsRes(res *usecase.GetUsersListTransactionsRes) *GetUsersListTransactionsRes {
	return &GetUsersListTransactionsRes{
		Transactions: toDeliveryArrTransactionDTO(res.Transactions),
		Avg:          res.Avg,
	}
}

func toUseCaseAddToolSetReq(req AddToolSetReq) usecase.AddToolSetReq {
	return usecase.AddToolSetReq{
		ToolSetName: req.ToolSetName,
		ToolsIds:    req.ToolsIds,
	}
}

func toDeliveryArrTransactionDTO(transactions []*usecase.TransactionDTO) []*TransactionDTO {
	res := make([]*TransactionDTO, len(transactions))
	for i, transaction := range transactions {
		res[i] = toDeliveryTransactionDTO(transaction)
	}

	return res
}

func toDeliveryGetAllTransactions(res *usecase.GetAllTransactions) *GetAllTransactions {
	return &GetAllTransactions{
		User:         toDeliveryUserDto(res.User),
		Transactions: toArrDeliveryLightTransaction(res.Transactions),
	}
}

func toArrDeliveryLightTransaction(res []usecase.LightTransaction) []LightTransaction {
	result := make([]LightTransaction, len(res))
	for i, v := range res {
		result[i] = toDeliveryLightTransaction(v)
	}

	return result
}

func toDeliveryLightTransaction(res usecase.LightTransaction) LightTransaction {
	return LightTransaction{
		Id:        res.Id,
		CreatedAt: res.CreatedAt,
	}
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
		Status:           res.Status,
	}
}

func toDeliveryProblematicTools(tools *usecase.ProblematicTools) *ProblematicTools {
	return &ProblematicTools{
		ManualCheckTools: toArrDeliveryRecognizedToolDTO(tools.ManualCheckTools),
		UnknownTools:     toArrDeliveryRecognizedToolDTO(tools.UnknownTools),
		MissingTools:     toArrDeliveryToolTypeDTO(tools.MissingTools),
	}
}

func toDeliveryQaTransactionsRes(res *usecase.QaTransactionsRes) *QaTransactionsRes {
	return &QaTransactionsRes{
		Qa:           toDeliveryUserDto(res.Qa),
		Transactions: toDeliveryArrTransactionResolutionDTO(res.Transactions),
	}
}

func toDeliveryArrTransactionResolutionDTO(dto []*usecase.TransactionResolutionDTO) []*TransactionResolutionDTO {
	res := make([]*TransactionResolutionDTO, len(dto))
	for i, d := range dto {
		res[i] = toDeliveryTransactionResolutionDTO(d)
	}

	return res
}

func toDeliveryTransactionResolutionDTO(d *usecase.TransactionResolutionDTO) *TransactionResolutionDTO {
	return &TransactionResolutionDTO{
		Transaction: toDeliveryTransactionDTO(d.Transaction),
		Reason:      d.Reason,
		Notes:       d.Notes,
		CreatedAt:   d.CreatedAt,
	}
}

func toDeliveryGetRolesRes(roles *usecase.GetRolesRes) GetRolesRes {
	return GetRolesRes{
		Roles: roles.Roles,
	}
}

func toDeliveryVerificationRes(res *usecase.VerificationRes) *VerificationRes {
	return &VerificationRes{
		TransactionID: res.TransactionID,
		Status:        res.Status,
		VerifiedBy:    res.VerifiedBy,
		CreatedAt:     res.CreatedAt,
	}
}

func ToDeliveryAddToolSetRes(res *usecase.AddToolSetRes) *AddToolSetRes {
	return &AddToolSetRes{
		Id:    res.Id,
		Name:  res.Name,
		Tools: toArrDeliveryToolTypeDTO(res.Tools),
	}
}

func toDeliveryGetTransactionStatisticsRes(res usecase.GetTransactionStatisticsRes) GetTransactionStatisticsRes {
	return GetTransactionStatisticsRes{
		Transactions:       res.Transactions,
		OpenedTransactions: res.OpenedTransactions,
		ClosedTransactions: res.ClosedTransactions,
		QATransactions:     res.QATransactions,
		FailedTransactions: res.FailedTransactions,
	}
}

func toDeliveryHumanErrorStats(res usecase.HumanErrorStats) HumanErrorStats {
	return HumanErrorStats{
		FullName:    res.FullName,
		EmployeeId:  res.EmployeeId,
		QAHitsCount: res.QAHitsCount,
	}
}

func toArrDeliveryHumanErrorStats(res []usecase.HumanErrorStats) []HumanErrorStats {
	result := make([]HumanErrorStats, len(res))
	for i, item := range res {
		result[i] = toDeliveryHumanErrorStats(item)
	}

	return result
}

func toDeliveryGetAvgWorkDurationRes(res *usecase.GetAvgWorkDurationRes) GetAvgWorkDurationRes {
	transactions := make([]GetAvgWorkDuration, len(res.Transactions))
	for i, transaction := range res.Transactions {
		transactions[i] = toDeliveryGetAvgWorkDuration(&transaction)
	}

	return GetAvgWorkDurationRes{
		Transactions: transactions,
	}
}

func toDeliveryGetAvgWorkDuration(res *usecase.GetAvgWorkDuration) GetAvgWorkDuration {
	return GetAvgWorkDuration{
		User:            toDeliveryUserDto(res.User),
		AvgWorkDuration: res.AvgWorkDuration,
	}
}

func toDeliveryMlErrorTransaction(res usecase.MlErrorTransaction) MlErrorTransaction {
	return MlErrorTransaction{
		TransactionID:  res.TransactionID,
		SourceImageUrl: res.SourceImageUrl,
		DebugImageUrl:  res.DebugImageUrl,
	}
}

func toArrDeliveryMlErrorTransaction(res []usecase.MlErrorTransaction) []MlErrorTransaction {
	result := make([]MlErrorTransaction, len(res))
	for i := range res {
		result[i] = toDeliveryMlErrorTransaction(res[i])
	}

	return result
}

func toDeliveryModelOrHumanStatsRes(res *usecase.ModelOrHumanStatsRes) *ModelOrHumanStatsRes {
	return &ModelOrHumanStatsRes{
		MlErrors:    res.MlErrors,
		HumanErrors: res.HumanErrors,
	}
}
