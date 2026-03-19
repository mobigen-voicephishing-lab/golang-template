package sample

import (
	"github.com/mobigen/golang-web-template/internal/domain"
)

// GetAllUseCase 전체 Sample 조회 유스케이스
type GetAllUseCase struct {
	repo Repository
}

// NewGetAllUseCase GetAllUseCase 생성자
func NewGetAllUseCase(repo Repository) *GetAllUseCase {
	return &GetAllUseCase{repo: repo}
}

// Execute 전체 Sample 조회 실행
func (uc *GetAllUseCase) Execute() (*[]domain.Sample, error) {
	return uc.repo.GetAll()
}
