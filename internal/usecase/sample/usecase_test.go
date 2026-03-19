package sample_test

import (
	"errors"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
	"github.com/stretchr/testify/assert"
)

// ── 묶음 방식 (SampleUseCase) 테스트 ──

func TestSampleUseCase_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewSampleUseCase(mockRepo)

	expected := &[]domain.Sample{
		{ID: 1, Name: "foo", Desc: "bar"},
		{ID: 2, Name: "baz", Desc: "qux"},
	}
	mockRepo.EXPECT().GetAll().Return(expected, nil).Times(1)

	result, err := uc.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSampleUseCase_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewSampleUseCase(mockRepo)

	mockRepo.EXPECT().GetAll().Return(nil, errors.New("db error")).Times(1)

	result, err := uc.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSampleUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewSampleUseCase(mockRepo)

	expected := &domain.Sample{ID: 1, Name: "foo", Desc: "bar"}
	mockRepo.EXPECT().GetByID(1).Return(expected, nil).Times(1)

	result, err := uc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSampleUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewSampleUseCase(mockRepo)

	input := &domain.Sample{Name: "new sample", Desc: "new desc"}
	expected := &domain.Sample{ID: 10, Name: "new sample", Desc: "new desc"}
	mockRepo.EXPECT().Create(input).Return(expected, nil).Times(1)

	result, err := uc.Create(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSampleUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewSampleUseCase(mockRepo)

	input := &domain.Sample{ID: 1, Name: "updated", Desc: "updated desc"}
	mockRepo.EXPECT().Update(input).Return(input, nil).Times(1)

	result, err := uc.Update(input)

	assert.NoError(t, err)
	assert.Equal(t, input, result)
}

func TestSampleUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewSampleUseCase(mockRepo)

	expected := &domain.Sample{ID: 1, Name: "foo", Desc: "bar"}
	mockRepo.EXPECT().Delete(1).Return(expected, nil).Times(1)

	result, err := uc.Delete(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
