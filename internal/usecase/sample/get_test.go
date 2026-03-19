package sample_test

import (
	"errors"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
	"github.com/stretchr/testify/assert"
)

func TestGetByIDUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewGetByIDUseCase(mockRepo)

	expected := &domain.Sample{ID: 1, Name: "foo", Desc: "bar"}
	mockRepo.EXPECT().GetByID(1).Return(expected, nil).Times(1)

	result, err := uc.Execute(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetByIDUseCase_Execute_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewGetByIDUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(999).Return(nil, errors.New("not found")).Times(1)

	result, err := uc.Execute(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
