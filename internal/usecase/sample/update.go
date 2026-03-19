package sample

import (
	"github.com/mobigen/golang-web-template/internal/domain"
)

// UpdateUseCase Sample 수정 유스케이스
type UpdateUseCase struct {
	repo Repository
}

// NewUpdateUseCase UpdateUseCase 생성자
func NewUpdateUseCase(repo Repository) *UpdateUseCase {
	return &UpdateUseCase{repo: repo}
}

// Execute Sample 수정 실행
func (uc *UpdateUseCase) Execute(sample *domain.Sample) (*domain.Sample, error) {
	return uc.repo.Update(sample)
}
