package sample

import (
	"github.com/mobigen/golang-web-template/internal/domain"
)

// GetByIDUseCase ID로 Sample 조회 유스케이스
type GetByIDUseCase struct {
	repo Repository
}

// NewGetByIDUseCase GetByIDUseCase 생성자
func NewGetByIDUseCase(repo Repository) *GetByIDUseCase {
	return &GetByIDUseCase{repo: repo}
}

// Execute ID로 Sample 조회 실행
func (uc *GetByIDUseCase) Execute(id int) (*domain.Sample, error) {
	return uc.repo.GetByID(id)
}
