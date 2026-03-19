package sample

import (
	"github.com/mobigen/golang-web-template/internal/domain"
)

// CreateUseCase Sample 생성 유스케이스
type CreateUseCase struct {
	repo Repository
}

// NewCreateUseCase CreateUseCase 생성자
func NewCreateUseCase(repo Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

// Execute Sample 생성 실행
func (uc *CreateUseCase) Execute(sample *domain.Sample) (*domain.Sample, error) {
	return uc.repo.Create(sample)
}
