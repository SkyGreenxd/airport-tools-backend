package usecase

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/repository"
	"airport-tools-backend/pkg/e"
	"context"
	"log"
)

// TODO: заменить на реальные данные
const (
	ConfidenceCompare float32 = 0.85
	CosineSimCompare  float64 = 0.70
)

type Service struct {
	userRepo         repository.UserRepository
	cvScanRepo       repository.CvScanRepository
	cvScanDetailRepo repository.CvScanDetailRepository
	toolTypeRepo     repository.ToolTypeRepository
	transactionRepo  repository.TransactionRepository
	mlGateway        MLGateway
	toolSetRepo      repository.ToolSetRepository
	imageStorage     ImageStorage
}

func NewService(
	u repository.UserRepository, c repository.CvScanRepository, cd repository.CvScanDetailRepository,
	tt repository.ToolTypeRepository, t repository.TransactionRepository, ml MLGateway, s3 ImageStorage,
	ts repository.ToolSetRepository,
) *Service {
	return &Service{
		userRepo:         u,
		cvScanRepo:       c,
		cvScanDetailRepo: cd,
		toolTypeRepo:     tt,
		transactionRepo:  t,
		mlGateway:        ml,
		imageStorage:     s3,
		toolSetRepo:      ts,
	}
}

func (s *Service) MlService() (string, error) {
	return "", nil
}

// Checkout обрабатывает выдачу инструментов инженеру
func (s *Service) Checkout(ctx context.Context, req *CheckReq) (res *CheckRes, err error) {
	const op = "usecase.Checkout"

	user, err := s.userRepo.GetByEmployeeIdWithTransactions(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	if err := user.CanCheckout(); err != nil {
		return nil, e.Wrap(op, err)
	}

	// сохранение фотки в s3
	uploadImageRes, err := s.imageStorage.UploadImage(ctx, req.Data)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	scanReq := NewScanReq(uploadImageRes.Key, uploadImageRes.ImageUrl, ConfidenceCompare)
	scanResult, err := s.mlGateway.ScanTools(ctx, scanReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	referenceSet, err := s.toolSetRepo.GetByIdWithTools(ctx, user.DefaultToolSetId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	newTransaction := domain.NewTransaction(user.Id, referenceSet.Id)
	transaction, err := s.transactionRepo.Create(ctx, newTransaction)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	createScanReq := NewCreateScanReq(transaction.Id, domain.Checkin, uploadImageRes.ImageUrl, scanResult.Tools)
	if err := s.CreateScan(ctx, createScanReq); err != nil {
		return nil, e.Wrap(op, err)
	}

	filterReq := NewFilterReq(ConfidenceCompare, CosineSimCompare, scanResult.Tools, referenceSet.Tools)
	filterRes, err := filterRecognizedTools(filterReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewCheckinRes(uploadImageRes.ImageUrl, scanResult.DebugImageUrl, filterRes.AccessTools, filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools), nil
}

// Checkin обрабатывает возврат инструментов инженером
func (s *Service) Checkin(ctx context.Context, req *CheckReq) (res *CheckRes, err error) {
	const op = "usecase.Checkin"

	user, err := s.userRepo.GetByEmployeeIdWithTransactions(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	if err := user.CanCheckin(); err != nil {
		return nil, e.Wrap(op, err)
	}

	transaction, err := s.transactionRepo.GetByUserIdWhereStatusIsOpenOrManual(ctx, user.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	uploadImage, err := s.imageStorage.UploadImage(ctx, req.Data)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	scanReq := NewScanReq(uploadImage.Key, uploadImage.ImageUrl, ConfidenceCompare)
	scanResult, err := s.mlGateway.ScanTools(ctx, scanReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	referenceSet, err := s.toolSetRepo.GetByIdWithTools(ctx, user.DefaultToolSetId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	createScanReq := NewCreateScanReq(transaction.Id, domain.Checkout, uploadImage.ImageUrl, scanResult.Tools)
	if err := s.CreateScan(ctx, createScanReq); err != nil {
		return nil, e.Wrap(op, err)
	}

	filterReq := NewFilterReq(ConfidenceCompare, CosineSimCompare, scanResult.Tools, referenceSet.Tools)
	filterRes, err := filterRecognizedTools(filterReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	transaction.EvaluateStatus(len(filterRes.ManualCheckTools), len(filterRes.UnknownTools), len(filterRes.MissingTools))

	if _, err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewCheckinRes(uploadImage.ImageUrl, scanResult.DebugImageUrl, filterRes.AccessTools, filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools), nil
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

	newScan := domain.NewCvScan(req.TransactionId, req.ScanType, req.ImageUrl)
	scan, err := s.cvScanRepo.Create(ctx, newScan)
	if err != nil {
		return e.Wrap(op, err)
	}

	for _, recognized := range req.Tools {
		if _, exists := toolMap[recognized.ToolTypeId]; exists {
			if len(recognized.Embedding) == 0 {
				recognized.Embedding = make([]float32, 1280)
			}
			scanDetail := domain.NewCvScanDetail(scan.Id, recognized.ToolTypeId, recognized.Confidence, recognized.Embedding)
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
