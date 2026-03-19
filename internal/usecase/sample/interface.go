package sample

//go:generate go run go.uber.org/mock/mockgen -destination=mock_repository_test.go -package=sample_test github.com/mobigen/golang-web-template/internal/usecase/sample Repository

import (
	"github.com/mobigen/golang-web-template/internal/domain"
)

// Repository sample repository 인터페이스 (usecase가 의존하는 포트)
type Repository interface {
	GetAll() (*[]domain.Sample, error)
	GetByID(int) (*domain.Sample, error)
	Create(*domain.Sample) (*domain.Sample, error)
	Update(*domain.Sample) (*domain.Sample, error)
	Delete(int) (*domain.Sample, error)
}
