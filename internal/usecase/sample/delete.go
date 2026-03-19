package sample

import (
	"github.com/mobigen/golang-web-template/internal/domain"
)

// DeleteUseCase Sample 삭제 유스케이스
type DeleteUseCase struct {
	repo Repository
}

// NewDeleteUseCase DeleteUseCase 생성자
func NewDeleteUseCase(repo Repository) *DeleteUseCase {
	return &DeleteUseCase{repo: repo}
}

// Execute Sample 삭제 실행
func (uc *DeleteUseCase) Execute(id int) (*domain.Sample, error) {
	return uc.repo.Delete(id)
}
