package usecase

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/repository"
	"airport-tools-backend/pkg/e"
	"airport-tools-backend/pkg/logger"
	"context"
	"errors"
	"log"
	"strings"
	"time"
)

// TODO: заменить на реальные данные
const (
	SourceImages string = "source_images"
	DefaultSetId int64  = 1
	Checkin      string = "Checkin"
	Checkout     string = "Checkout"
)

type Service struct {
	userRepo          repository.UserRepository
	cvScanRepo        repository.CvScanRepository
	cvScanDetailRepo  repository.CvScanDetailRepository
	toolTypeRepo      repository.ToolTypeRepository
	transactionRepo   repository.TransactionRepository
	mlGateway         MLGateway
	toolSetRepo       repository.ToolSetRepository
	imageStorage      ImageStorage
	ConfidenceCompare float32
	CosineSimCompare  float32
	trResolution      repository.TransactionResolutionsRepository
	logger            logger.Logger
	roleRepo          repository.RoleRepository
}

func NewService(
	u repository.UserRepository, c repository.CvScanRepository, cd repository.CvScanDetailRepository,
	tt repository.ToolTypeRepository, t repository.TransactionRepository, ml MLGateway, s3 ImageStorage,
	ts repository.ToolSetRepository, condfidence, cosineSim float32, tr repository.TransactionResolutionsRepository,
	logger logger.Logger, roleRepo repository.RoleRepository,
) *Service {
	return &Service{
		userRepo:          u,
		cvScanRepo:        c,
		cvScanDetailRepo:  cd,
		toolTypeRepo:      tt,
		transactionRepo:   t,
		mlGateway:         ml,
		imageStorage:      s3,
		toolSetRepo:       ts,
		ConfidenceCompare: condfidence,
		CosineSimCompare:  cosineSim,
		trResolution:      tr,
		logger:            logger,
		roleRepo:          roleRepo,
	}
}

func (s *Service) Check(ctx context.Context, req *CheckReq) (*CheckRes, error) {
	const op = "usecase.Check"
	var res *CheckRes

	user, err := s.userRepo.GetByEmployeeIdWithTransactions(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	transactionProcess := NewTransactionProcess(user.Id, req.Data, req.ToolSetId)

	if err := user.CanCheckout(); err != nil {
		if err := user.CanCheckin(); err != nil {
			return nil, e.Wrap(op, err)
		}

		err := s.logger.Track("usecase.Checkin", func() error {
			res, err = s.Checkin(ctx, transactionProcess)
			return err
		})

		if err != nil {
			return nil, e.Wrap(op, err)
		}

		return res, nil
	}

	err = s.logger.Track("usecase.Checkout", func() error {
		res, err = s.Checkout(ctx, transactionProcess)
		return err
	})
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return res, nil
}

// Checkout обрабатывает выдачу инструментов инженеру
func (s *Service) Checkout(ctx context.Context, req *TransactionProcess) (res *CheckRes, err error) {
	const op = "usecase.Checkout"

	toolSetId := req.ToolSetId
	if toolSetId == 0 {
		toolSetId = DefaultSetId
	}

	referenceSet, err := s.toolSetRepo.GetByIdWithTools(ctx, toolSetId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var uploadImageRes *UploadImageRes
	uplImageReq := NewUploadImageReq(req.Data, SourceImages)
	err = s.logger.Track("usecase.Checkout.imageStorage.UploadImage", func() error {
		uploadImageRes, err = s.imageStorage.UploadImage(ctx, uplImageReq)
		return err
	})
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var scanResult *ScanResult
	scanReq := NewScanReq(uploadImageRes.Key, uploadImageRes.ImageUrl, s.ConfidenceCompare)
	err = s.logger.Track("usecase.Checkout.mlGateway.ScanTools", func() error {
		scanResult, err = s.mlGateway.ScanTools(ctx, scanReq)
		return err
	})
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	filterReq := NewFilterReq(s.ConfidenceCompare, s.CosineSimCompare, scanResult.Tools, referenceSet.Tools)
	filterRes, err := filterRecognizedTools(filterReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	hasLowConfidence := false
	for _, tool := range filterRes.ManualCheckTools {
		if tool.Confidence <= 50 {
			hasLowConfidence = true
			break
		}
	}

	var status domain.Status
	if len(filterRes.MissingTools) > 0 || len(filterRes.UnknownTools) > 0 || ((len(filterRes.AccessTools) + len(filterRes.ManualCheckTools)) != len(referenceSet.Tools)) || hasLowConfidence {
		status = domain.FAILED
	} else {
		status = domain.OPEN
	}

	var transaction *domain.Transaction
	existing, err := s.transactionRepo.GetLastFailedByUserId(ctx, req.UserId)
	if err != nil && !errors.Is(err, e.ErrTransactionNotFound) {
		return nil, e.Wrap(op, err)
	}

	if existing != nil {
		existing.Status = status
		transaction, err = s.transactionRepo.Update(ctx, existing)
		if err != nil {
			return nil, e.Wrap(op, err)
		}
	} else {
		newTransaction := domain.NewTransaction(req.UserId, referenceSet.Id, status)
		transaction, err = s.transactionRepo.Create(ctx, newTransaction)
		if err != nil {
			return nil, e.Wrap(op, err)
		}
	}

	createScanReq := NewCreateScanReq(transaction.Id, domain.Checkout, uploadImageRes.ImageUrl, scanResult.DebugImageUrl, scanResult.Tools)
	if err := s.CreateScan(ctx, createScanReq); err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewCheckinRes(uploadImageRes.ImageUrl, scanResult.DebugImageUrl, filterRes.AccessTools, filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools, Checkout, string(transaction.Status)), nil
}

// Checkin обрабатывает возврат инструментов инженером
func (s *Service) Checkin(ctx context.Context, req *TransactionProcess) (res *CheckRes, err error) {
	const op = "usecase.Checkin"

	transaction, err := s.transactionRepo.GetByUserIdWhereStatusIsOpenOrQA(ctx, req.UserId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// проверка на 3 и более попыток
	if err := transaction.CheckCountOfChecks(); err != nil {
		return nil, e.Wrap(op, err)
	}

	referenceSet, err := s.toolSetRepo.GetByIdWithTools(ctx, transaction.ToolSetId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	uplImageReq := NewUploadImageReq(req.Data, SourceImages)

	var uploadImage *UploadImageRes
	err = s.logger.Track("usecase.Checkin.imageStorage.UploadImage", func() error {
		uploadImage, err = s.imageStorage.UploadImage(ctx, uplImageReq)
		return err
	})
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var scanResult *ScanResult
	scanReq := NewScanReq(uploadImage.Key, uploadImage.ImageUrl, s.ConfidenceCompare)
	err = s.logger.Track("usecase.Checkin.mlGateway.ScanTools", func() error {
		scanResult, err = s.mlGateway.ScanTools(ctx, scanReq)
		return err
	})
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	createScanReq := NewCreateScanReq(transaction.Id, domain.Checkin, uploadImage.ImageUrl, scanResult.DebugImageUrl, scanResult.Tools)
	if err := s.CreateScan(ctx, createScanReq); err != nil {
		return nil, e.Wrap(op, err)
	}

	filterReq := NewFilterReq(s.ConfidenceCompare, s.CosineSimCompare, scanResult.Tools, referenceSet.Tools)
	filterRes, err := filterRecognizedTools(filterReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	transaction.CountOfChecks++
	transaction.EvaluateStatus(len(filterRes.ManualCheckTools), len(filterRes.UnknownTools), len(filterRes.MissingTools))
	transaction.UpdatedAt = time.Now()

	if _, err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewCheckinRes(uploadImage.ImageUrl, scanResult.DebugImageUrl, filterRes.AccessTools, filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools, Checkin, string(transaction.Status)), nil
}

// CreateScan создает записи в таблицы cv_scans, cv_scan_details
func (s *Service) CreateScan(ctx context.Context, req *CreateScanReq) error {
	const op = "usecase.CreateScan"

	tools, err := s.toolTypeRepo.GetAll(ctx)
	if err != nil {
		return e.Wrap(op, err)
	}

	toolMap := make(map[int64]*domain.ToolType)
	for _, t := range tools {
		toolMap[t.Id] = t
	}

	newScan := domain.NewCvScan(req.TransactionId, req.ScanType, req.ImageUrl, req.DebugImageUrl)
	scan, err := s.cvScanRepo.Create(ctx, newScan)
	if err != nil {
		return e.Wrap(op, err)
	}

	for _, recognized := range req.Tools {
		if _, exists := toolMap[recognized.ToolTypeId]; exists {
			if len(recognized.Embedding) == 0 {
				recognized.Embedding = make([]float32, 1280)
			}
			scanDetail := domain.NewCvScanDetail(scan.Id, recognized.ToolTypeId, recognized.Confidence, recognized.Embedding, recognized.Bbox)
			_, err := s.cvScanDetailRepo.Create(ctx, scanDetail)
			if err != nil {
				return e.Wrap(op, err)
			}
		} else {
			log.Printf("unknown tool type: %v", recognized.ToolTypeId)
		}
	}

	return nil
}

// List возвращает список транзакций, возможна фильтрация по статусу
func (s *Service) List(ctx context.Context, statusStr string) (*ListTransactionsRes, error) {
	const op = "usecase.List"

	if statusStr == "" {
		transactions, err := s.transactionRepo.GetAllWithUser(ctx)
		if err != nil {
			return nil, e.Wrap(op, err)
		}

		result := NewListTransactionsRes(toListTransactionsRes(transactions))
		return result, nil
	} else {
		statusStr = strings.ToUpper(statusStr)
		status, err := domain.ValidateStatus(statusStr)
		if err != nil {
			return nil, e.Wrap(op, err)
		}

		transactions, err := s.transactionRepo.GetAllWithStatusAndUser(ctx, status)
		if err != nil {
			return nil, e.Wrap(op, err)
		}

		result := NewListTransactionsRes(toListTransactionsRes(transactions))
		return result, nil
	}

	// return nil, e.ErrRequestNotSupported
}

// Login возвращает роль пользователя для дальнейшей работы. MVP вариант, небезопасно
func (s *Service) Login(ctx context.Context, req *LoginReq) (*LoginRes, error) {
	const op = "usecase.Login"

	user, err := s.userRepo.GetByEmployeeId(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewLoginRes(user.Role.Name), nil
}

// TODO: добавить таблицу ролей чтобы не хардкодить
// GetRoles возвращает список ролей
func (s *Service) GetRoles(ctx context.Context) (*GetRolesRes, error) {
	const op = "usecase.GetRoles"

	roles, err := s.roleRepo.GetAll(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	result := make([]string, len(roles))
	for i, role := range roles {
		result[i] = role.Name
	}

	return NewGetRolesRes(result), nil
}

// Register регистрирует пользователя
func (s *Service) Register(ctx context.Context, req *RegisterReq) (*RegisterRes, error) {
	const op = "usecase.Register"

	role, err := s.roleRepo.GetByName(ctx, req.Role)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	newUser := domain.NewUser(req.FullName, req.EmployeeId, role.Id)
	user, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewRegisterRes(user.Id), nil
}

// Verification отвечает за QA-проверку и завершение проблемной транзакции
func (s *Service) Verification(ctx context.Context, req *Verification) (*VerificationRes, error) {
	const op = "usecase.postVerification"

	if err := domain.ValidateReason(req.Reason); err != nil {
		return nil, e.Wrap(op, err)
	}

	user, err := s.userRepo.GetByEmployeeId(ctx, req.QAEmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	new_resolution := domain.NewTransactionResolution(req.TransactionID, user.Id, req.Reason, req.Notes)
	resolution, err := s.trResolution.Create(ctx, new_resolution, req.ToolsIds)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	transaction, err := s.transactionRepo.GetById(ctx, resolution.TransactionId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	transaction.Status = domain.CLOSED
	updTransaction, err := s.transactionRepo.Update(ctx, transaction)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewVerificationRes(updTransaction.Id, string(updTransaction.Status), user.EmployeeId, resolution.CreatedAt), nil
}

// GetQATransaction возвращает структурированное описание об инструментах с привязкой к конкретной транзакции и изображению.
func (s *Service) GetQATransaction(ctx context.Context, transactionId int64) (*GetQAVerificationRes, error) {
	const op = "usecase.GetQATransaction"

	scan, err := s.cvScanRepo.GetByTransactionIdWithDetectedToolsAndTransaction(ctx, transactionId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	toolSet, err := s.toolSetRepo.GetByIdWithTools(ctx, scan.TransactionObj.ToolSetId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	detectedTools := make([]*domain.RecognizedTool, len(scan.DetectedTools))
	for i, tool := range scan.DetectedTools {
		detectedTools[i] = domain.NewRecognizedTool(tool.DetectedToolTypeId, tool.Confidence, tool.Embedding, tool.Bbox)
	}

	filterReq := NewFilterReq(s.ConfidenceCompare, s.CosineSimCompare, detectedTools, toolSet.Tools)
	filterRes, err := filterRecognizedTools(filterReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	problematicTools := NewProblematicTools(filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools)
	userDto := NewUserDto(scan.TransactionObj.User.FullName, scan.TransactionObj.User.EmployeeId)
	res := NewGetQAVerificationRes(scan.TransactionId, toolSet.Id, scan.TransactionObj.CreatedAt, userDto, filterRes.AccessTools, problematicTools, scan.ImageUrl, string(scan.TransactionObj.Status))

	return res, nil
}

// UserTransactions возвращает список транзакций конкретного пользователя
func (s *Service) UserTransactions(ctx context.Context, req *UserTransactionsReq) (*GetUsersListTransactionsRes, error) {
	const op = "usecase.UserTransactions"

	user, err := s.userRepo.GetByEmployeeId(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	transactions, err := s.transactionRepo.GetAllByUserId(ctx, user.Id, req.StartDate, req.EndDate, req.Limit)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var avgHours float64
	if req.Avg == true {
		var totalHours float64
		for _, tr := range transactions {
			totalHours += tr.UpdatedAt.Sub(tr.CreatedAt).Hours()
		}

		if len(transactions) > 0 {
			avgHours = totalHours / float64(len(transactions))
		}
	}

	result := NewGetUsersListTransactionsRes(toListTransactionsRes(transactions), avgHours)

	return result, nil
}

// GetUsersQAStats возвращает инженеров, чьи транзакции попадали на QA по причине HUMAN_ERR,
// по убыванию кол-ва проверок
func (s *Service) GetUsersQAStats(ctx context.Context) ([]HumanErrorStats, error) {
	const op = "usecase.GetUsersQAStats"

	users, err := s.trResolution.GetTopHumanErrorUsers(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var result []HumanErrorStats
	for _, user := range users {
		result = append(result, NewHumanErrorStats(user.FullName, user.EmployeeId, user.QAHitsCount))
	}

	return result, nil
}

// GetQAChecks возвращает какие проверки делал сотрудник QA
func (s *Service) GetQAChecks(ctx context.Context, qaId string) (*QaTransactionsRes, error) {
	const op = "usecase.GetQAChecks"

	qa, err := s.userRepo.GetByEmployeeId(ctx, qaId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	result, err := s.trResolution.GetByQAId(ctx, qa.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	res := NewQaTransactionsRes(NewUserDto(qa.FullName, qa.EmployeeId), ToListTransactionResolutionDTO(result))
	return res, nil
}

// GetMlVsHuman возвращает два числа "Model vs Human errors"
func (s *Service) GetMlVsHuman(ctx context.Context) (*ModelOrHumanStatsRes, error) {
	const op = "usecase.GetMlVsHuman"

	humansErrors, err := s.trResolution.GetAllHumanError(ctx)
	if err != nil && !errors.Is(err, e.ErrTransactionResolutionsNotFound) {
		return nil, e.Wrap(op, err)
	}

	modelErrors, err := s.trResolution.GetAllModelError(ctx)
	if err != nil && !errors.Is(err, e.ErrTransactionResolutionsNotFound) {
		return nil, e.Wrap(op, err)
	}

	return NewModelOrHumanStatsRes(len(modelErrors), len(humansErrors)), nil
}

// GetAllQaEmployers возваращает всех QA проверяющих
func (s *Service) GetAllQaEmployers(ctx context.Context) ([]UserDto, error) {
	const op = "usecase.GetAllQaEmployers"

	qaEmployers, err := s.userRepo.GetAllQa(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var result []UserDto
	for _, qa := range qaEmployers {
		result = append(result, NewUserDto(qa.FullName, qa.EmployeeId))
	}

	return result, nil
}

// GetTransactionStatistics возвращает список транзакций по типам
func (s *Service) GetTransactionStatistics(ctx context.Context) (*GetTransactionStatisticsRes, error) {
	const op = "usecase.GetTransactionStatistics"

	transactions, err := s.transactionRepo.GetAll(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	opened, err := s.transactionRepo.GetAllWithStatus(ctx, domain.OPEN)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	closed, err := s.transactionRepo.GetAllWithStatus(ctx, domain.CLOSED)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	qa, err := s.transactionRepo.GetAllWithStatus(ctx, domain.QA)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	failed, err := s.transactionRepo.GetAllWithStatus(ctx, domain.FAILED)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewGetTransactionStatisticsRes(len(transactions), len(opened), len(closed), len(qa), len(failed)), nil
}

// GetAvgWorkDuration возвращает среднее время работы каждого инженера по всем его транзакциям
func (s *Service) GetAvgWorkDuration(ctx context.Context) (*GetAvgWorkDurationRes, error) {
	const op = "usecase.GetAvgWorkDuration"

	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// Только инженеры
	engineers := make([]*domain.User, 0, len(users))
	userIds := make([]int64, 0, len(users))
	for _, u := range users {
		if u.Role.Name == domain.Engineer {
			engineers = append(engineers, u)
			userIds = append(userIds, u.Id)
		}
	}

	transactions, err := s.transactionRepo.GetByUserIds(ctx, userIds)
	if err != nil && !errors.Is(err, e.ErrTransactionNotFound) {
		return nil, e.Wrap(op, err)
	}

	userTxMap := make(map[int64][]*domain.Transaction)
	for _, tx := range transactions {
		userTxMap[tx.UserId] = append(userTxMap[tx.UserId], tx)
	}

	result := make([]GetAvgWorkDuration, len(engineers))
	for i, user := range engineers {
		txs := userTxMap[user.Id]
		var totalHours float64
		for _, tx := range txs {
			totalHours += tx.UpdatedAt.Sub(tx.CreatedAt).Hours()
		}

		var avgHours float64
		if len(txs) > 0 {
			avgHours = totalHours / float64(len(txs))
		}

		result[i] = NewGetAvgWorkDuration(NewUserDto(user.FullName, user.EmployeeId), avgHours)
	}

	return NewGetAvgWorkDurationRes(result), nil
}

// GetMlErrorTransactions выводит список транзакций, в которых модель ошиблась
func (s *Service) GetMlErrorTransactions(ctx context.Context) ([]MlErrorTransaction, error) {
	const op = "usecase.GetMlErrorTransactions"

	qaTransactions, err := s.trResolution.GetMlErrorTransactions(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var result []MlErrorTransaction
	for _, qa := range qaTransactions {
		tx := qa.Transaction
		if tx == nil || len(tx.CvScans) == 0 {
			continue
		}

		var lastCheckin *domain.CvScan
		for _, scan := range tx.CvScans {
			switch scan.ScanType {
			case domain.Checkin:
				if lastCheckin == nil || scan.CreatedAt.After(lastCheckin.CreatedAt) {
					lastCheckin = scan
				}
			}
		}
		if lastCheckin != nil {
			result = append(result, *NewMlErrorTransaction(tx.Id, lastCheckin.ImageUrl, lastCheckin.DebugImageUrl))
		}
	}

	return result, nil
}

// Получить транзакции всех юзеров
func (s *Service) GetAllTransactions(ctx context.Context) ([]*GetAllTransactions, error) {
	const op = "usecase.GetAllTransactions"

	users, err := s.userRepo.GetAllEngineersWithTransactions(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	result := make([]*GetAllTransactions, len(users))
	for i, u := range users {
		result[i] = NewGetAllTransactions(u, NewArrLightTransaction(u.Transactions))
	}

	return result, nil
}

// Создать новый сет инструментов
func (s *Service) AddToolSet(ctx context.Context, req AddToolSetReq) (*AddToolSetRes, error) {
	const op = "usecase.AddToolSet"

	newSet := domain.NewToolSet(req.ToolSetName)

	res, err := s.toolSetRepo.CreateWithTools(ctx, newSet, req.ToolsIds)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewAddToolSetRes(res.Id, res.Name, toArrToolTypeDTO(res.Tools)), nil
}

// GetMlErrorTools возвращает наборы инструментов вместе с инструментами,
// где для каждого инструмента указано,
// сколько раз на нём была зарегистрирована ошибка MODEL_ERR
func (s *Service) GetMlErrorTools(ctx context.Context) ([]*repository.ToolSetWithErrors, error) {
	const op = "usecase.GetMlErrorTools"

	res, err := s.trResolution.GetMlErrorTools(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return res, nil
}
