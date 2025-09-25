package usecase

import "context"

type MLGateway interface {
	ScanTools(ctx context.Context, req *ScanRequest) (*ScanResult, error)
}
