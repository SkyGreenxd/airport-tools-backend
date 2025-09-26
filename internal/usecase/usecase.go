package usecase

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/repository"
	"airport-tools-backend/pkg/e"
	"context"
	"fmt"
	"log"
	"math"
)

// TODO: заменить на реальные данные
const (
	ConfidenceCompare float32 = 0.90
	CosineSimCompare  float64 = 0.6
)

type Service struct {
	userRepo         repository.UserRepository
	cvScanRepo       repository.CvScanRepository
	cvScanDetailRepo repository.CvScanDetailRepository
	toolTypeRepo     repository.ToolTypeRepository
	transactionRepo  repository.TransactionRepository
	mlGateway        MLGateway
	yandexS3         repository.ImageRepository
	toolSetRepo      repository.ToolSetRepository
}

func NewService(
	u repository.UserRepository, c repository.CvScanRepository, cd repository.CvScanDetailRepository,
	tt repository.ToolTypeRepository, t repository.TransactionRepository, ml MLGateway, y repository.ImageRepository,
	ts repository.ToolSetRepository,
) *Service {
	return &Service{
		userRepo:         u,
		cvScanRepo:       c,
		cvScanDetailRepo: cd,
		toolTypeRepo:     tt,
		transactionRepo:  t,
		mlGateway:        ml,
		yandexS3:         y,
		toolSetRepo:      ts,
	}
}

func (s *Service) MlService() (string, error) {
	return "", nil
}

// TODO: доделать, реализовать возвращение айди (бакет+имя)
func (s *Service) UploadImage(ctx context.Context, req ImageReq) (*UploadImageRes, error) {
	newImage := domain.NewImage(req.Filename, int64(len(req.Data)))
	image, err := s.yandexS3.Save(ctx, newImage)
	if err != nil {
		return nil, err
	}

	return NewUploadImageRes(image.ImageId, image.ImageUrl), nil
}

// Checkout отвечает за выдачу инструментов инженеру
func (s *Service) Checkout(ctx context.Context, req *CheckReq) (res *CheckRes, err error) {
	const op = "usecase.Checkout"

	// проверка что фото не пустое, потом вынести можно в хэндлер
	//if req.Image.Filename == "" || req.Image.ContentType == "" || len(req.Image.Data) == 0 {
	//	return nil, e.Wrap(op, e.ErrEmptyFields)
	//}

	// проверка что инженер существует
	user, err := s.userRepo.GetByEmployeeId(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// сохранение фотки в s3
	uploadImageRes, err := s.UploadImage(ctx, req.Image)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// отправка фото в ML сервис
	scanReq := NewScanReq(uploadImageRes.ImageId, uploadImageRes.ImageUrl)
	scanResult, err := s.mlGateway.ScanTools(ctx, scanReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// Поиск сета
	referenceSet, err := s.toolSetRepo.GetByIdWithTools(ctx, user.DefaultToolSetId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// создание транзакции
	// TODO: нельзя сделать новую выдачу инженеру если не закрыта прошлая транзакция
	// TODO: мб не оч хорошо, что если в дальнейшем будут ошибки, то транзакция открыта
	newTransaction := domain.NewTransaction(user.Id, referenceSet.Id)
	transaction, err := s.transactionRepo.Create(ctx, newTransaction)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// добавление записей в бд cv_scan, cv_scan_detail
	createScanReq := NewCreateScanReq(transaction.Id, domain.Checkin, uploadImageRes.ImageUrl, scanResult.Tools)
	if err := s.CreateScan(ctx, createScanReq); err != nil {
		return nil, e.Wrap(op, err)
	}

	// проверка скана
	// TODO: 0.98 и 0.90 лучше вынести в конфигурацию.
	filterReq := NewFilterReq(ConfidenceCompare, CosineSimCompare, scanResult.Tools, referenceSet.Tools)
	filterRes, err := filterRecognizedTools(filterReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewCheckinRes(uploadImageRes.ImageUrl, filterRes.AccessTools, filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools), nil
}

// TODO: допилить reason/status
// Checkin отвечает за возврат инструментов инженером
func (s *Service) Checkin(ctx context.Context, req *CheckReq) (res *CheckRes, err error) {
	const op = "usecase.Checkin"

	// проверка что инженер существует
	//TODO: мб можно сделать будет проверку на то что в транзакции у юзера тот айди???
	user, err := s.userRepo.GetByEmployeeId(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// сохранение фотки в s3
	uploadImage, err := s.UploadImage(ctx, req.Image)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// отправка фото в ML сервис
	scanReq := NewScanReq(uploadImage.ImageId, uploadImage.ImageUrl)
	scanResult, err := s.mlGateway.ScanTools(ctx, scanReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// Поиск сета
	referenceSet, err := s.toolSetRepo.GetByIdWithTools(ctx, user.DefaultToolSetId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	transaction, err := s.transactionRepo.GetByUserIdWhereStatusIsOpenOrManual(ctx, user.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	createScanReq := NewCreateScanReq(transaction.Id, domain.Checkout, uploadImage.ImageUrl, scanResult.Tools)
	if err := s.CreateScan(ctx, createScanReq); err != nil {
		return nil, e.Wrap(op, err)
	}

	// проверка скана
	filterReq := NewFilterReq(ConfidenceCompare, CosineSimCompare, scanResult.Tools, referenceSet.Tools)
	filterRes, err := filterRecognizedTools(filterReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// TODO: при ошибках тоже надо как то менять статус
	var status domain.Status
	var reason domain.Reason
	if len(filterRes.ManualCheckTools) == 0 && len(filterRes.UnknownTools) == 0 && len(filterRes.MissingTools) == 0 {
		status = domain.CLOSED
		reason = domain.RETURNED
	} else {
		status = domain.MANUAL
		reason = domain.PROBLEMS
	}

	transaction.Status = status
	transaction.Reason = &reason
	if _, err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewCheckinRes(uploadImage.ImageUrl, filterRes.AccessTools, filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools), nil
}

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
			scanDetail := domain.NewCvScanDetail(scan.Id, recognized.ToolTypeId, recognized.Confidence, recognized.HashTool, recognized.Embedding)
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

func filterRecognizedTools(req *FilterReq) (*FilterRes, error) {
	accessTools := make([]*domain.RecognizedTool, 0, len(req.Tools))
	manualCheckTools := make([]*domain.RecognizedTool, 0, len(req.Tools))
	unknownTools := make([]*domain.RecognizedTool, 0, len(req.Tools))
	missingTools := make([]*ToolTypeDTO, 0)

	// создаём мапу ReferenceTools для быстрого поиска
	// создаём мапу ReferenceTools для быстрого поиска
	refMap := make(map[int64]*domain.ToolType)
	for _, r := range req.ReferenceTools {
		refMap[r.Id] = r
	}

	recognizedMap := make(map[int64]*domain.RecognizedTool)
	for _, recognized := range req.Tools {
		ref, exists := refMap[recognized.ToolTypeId]
		if !exists {
			// инструмент неизвестен
			unknownTools = append(unknownTools, recognized)
			continue
		}

		cosSim := cosineSimilarity(ref.ReferenceEmbedding, recognized.Embedding)
		fmt.Printf("DEBUG: toolId: %d, Confidence: %f, cosSim: %f\n", recognized.ToolTypeId, recognized.Confidence, cosSim)
		if cosSim >= req.CosineSimCompare && recognized.Confidence >= req.ConfidenceCompare {
			accessTools = append(accessTools, recognized)
		} else {
			manualCheckTools = append(manualCheckTools, recognized)
		}

		recognizedMap[recognized.ToolTypeId] = recognized
	}

	// отсутствующие инструменты из референсного набора
	for _, ref := range req.ReferenceTools {
		if _, ok := recognizedMap[ref.Id]; !ok {
			missingTools = append(missingTools, ToToolTypeDTO(ref))
		}
	}

	return NewFilterRes(accessTools, manualCheckTools, unknownTools, missingTools), nil
}

func cosineSimilarity(reference, recognized []float32) float64 {
	var dot, normReference, normRecognized float64
	for i := range reference {
		dot += float64(reference[i] * recognized[i])
		normReference += float64(reference[i] * reference[i])
		normRecognized += float64(recognized[i] * recognized[i])
	}
	return dot / (math.Sqrt(normReference) * math.Sqrt(normRecognized))
}
