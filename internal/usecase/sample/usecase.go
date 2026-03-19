package sample

import (
	"github.com/mobigen/golang-web-template/internal/domain"
)

// SampleUseCase 묶음 방식 — 모든 메서드를 하나의 구조체에 포함
type SampleUseCase struct {
	Repo Repository
}

// NewSampleUseCase SampleUseCase 생성자
func NewSampleUseCase(repo Repository) *SampleUseCase {
	return &SampleUseCase{Repo: repo}
}

// GetAll returns all samples.
func (uc *SampleUseCase) GetAll() (*[]domain.Sample, error) {
	return uc.Repo.GetAll()
}

// GetByID returns sample whose ID matches.
func (uc *SampleUseCase) GetByID(id int) (*domain.Sample, error) {
	return uc.Repo.GetByID(id)
}

// Create creates a new sample.
func (uc *SampleUseCase) Create(sample *domain.Sample) (*domain.Sample, error) {
	return uc.Repo.Create(sample)
}

// Update updates a sample.
func (uc *SampleUseCase) Update(sample *domain.Sample) (*domain.Sample, error) {
	return uc.Repo.Update(sample)
}

// Delete deletes sample by id.
func (uc *SampleUseCase) Delete(id int) (*domain.Sample, error) {
	return uc.Repo.Delete(id)
}
