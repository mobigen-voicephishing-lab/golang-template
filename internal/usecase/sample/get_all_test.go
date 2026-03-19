package sample_test

import (
	"errors"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewGetAllUseCase(mockRepo)

	expected := &[]domain.Sample{
		{ID: 1, Name: "foo", Desc: "bar"},
		{ID: 2, Name: "baz", Desc: "qux"},
	}
	mockRepo.EXPECT().GetAll().Return(expected, nil).Times(1)

	result, err := uc.Execute()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetAllUseCase_Execute_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewGetAllUseCase(mockRepo)

	mockRepo.EXPECT().GetAll().Return(nil, errors.New("db error")).Times(1)

	result, err := uc.Execute()

	assert.Error(t, err)
	assert.Nil(t, result)
}
