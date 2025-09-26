package usecase

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"context"
)

// Checkout отвечает за выдачу инструментов инженеру
func (s *Service) MockCheckout(ctx context.Context, req *MockCheckReq) (res *CheckRes, err error) {
	const op = "usecase.MockCheckout"

	// проверка что фото не пустое, потом вынести можно в хэндлер
	//if req.Image.Filename == "" || req.Image.ContentType == "" || len(req.Image.Data) == 0 {
	//	return nil, e.Wrap(op, e.ErrEmptyFields)
	//}

	// проверка что инженер существует
	user, err := s.userRepo.GetByEmployeeId(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// проверка что у инженера нет открытых транзакций
	if _, err := s.transactionRepo.GetByUserIdWhereStatusIsOpenOrManual(ctx, user.Id); err == nil {
		return nil, e.Wrap(op, e.ErrTransactionUnfinished)
	}

	// отправка фото в ML сервис
	scanReq := NewScanReq(req.ImageId, req.ImageUrl)
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
	createScanReq := NewCreateScanReq(transaction.Id, domain.Checkin, req.ImageUrl, scanResult.Tools)
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

	return NewCheckinRes(req.ImageUrl, filterRes.AccessTools, filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools), nil
}

// TODO: допилить reason/status
// Checkin отвечает за возврат инструментов инженером
func (s *Service) MockCheckin(ctx context.Context, req *MockCheckReq) (res *CheckRes, err error) {
	const op = "usecase.MockCheckin"

	// проверка что инженер существует
	//TODO: мб можно сделать будет проверку на то что в транзакции у юзера тот айди???
	user, err := s.userRepo.GetByEmployeeId(ctx, req.EmployeeId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	transaction, err := s.transactionRepo.GetByUserIdWhereStatusIsOpenOrManual(ctx, user.Id)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// отправка фото в ML сервис
	scanReq := NewScanReq(req.ImageId, req.ImageUrl)
	scanResult, err := s.mlGateway.ScanTools(ctx, scanReq)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	// Поиск сета
	referenceSet, err := s.toolSetRepo.GetByIdWithTools(ctx, user.DefaultToolSetId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	createScanReq := NewCreateScanReq(transaction.Id, domain.Checkout, req.ImageUrl, scanResult.Tools)
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

	return NewCheckinRes(req.ImageUrl, filterRes.AccessTools, filterRes.ManualCheckTools, filterRes.UnknownTools, filterRes.MissingTools), nil
}
