package services_test

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/mobigen/golang-web-template/mocks"
	"github.com/mobigen/golang-web-template/models"
	"github.com/mobigen/golang-web-template/services"
	"github.com/stretchr/testify/assert"
)

func TestSampleService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSampleRepository(ctrl)
	svc := services.Sample{}.New(mockRepo)

	expected := &[]models.Sample{
		{ID: 1, Name: "foo", Desc: "bar"},
		{ID: 2, Name: "baz", Desc: "qux"},
	}

	// mock 설정: GetAll 이 한 번 호출되면 expected 결과를 반환
	mockRepo.EXPECT().GetAll().Return(expected, nil).Times(1)

	result, err := svc.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSampleService_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSampleRepository(ctrl)
	svc := services.Sample{}.New(mockRepo)

	mockRepo.EXPECT().GetAll().Return(nil, errors.New("db error")).Times(1)

	result, err := svc.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSampleService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSampleRepository(ctrl)
	svc := services.Sample{}.New(mockRepo)

	expected := &models.Sample{ID: 1, Name: "foo", Desc: "bar"}

	// mock 설정: GetByID(1) 이 한 번 호출되면 expected 결과를 반환
	mockRepo.EXPECT().GetByID(1).Return(expected, nil).Times(1)

	result, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSampleService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSampleRepository(ctrl)
	svc := services.Sample{}.New(mockRepo)

	input := &models.Sample{Name: "new sample", Desc: "new desc"}
	expected := &models.Sample{ID: 10, Name: "new sample", Desc: "new desc"}

	// mock 설정: Create 가 input 과 동일한 인자로 호출되면 expected 결과를 반환
	mockRepo.EXPECT().Create(input).Return(expected, nil).Times(1)

	result, err := svc.Create(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSampleService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSampleRepository(ctrl)
	svc := services.Sample{}.New(mockRepo)

	input := &models.Sample{ID: 1, Name: "updated", Desc: "updated desc"}

	mockRepo.EXPECT().Update(input).Return(input, nil).Times(1)

	result, err := svc.Update(input)

	assert.NoError(t, err)
	assert.Equal(t, input, result)
}

func TestSampleService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSampleRepository(ctrl)
	svc := services.Sample{}.New(mockRepo)

	expected := &models.Sample{ID: 1, Name: "foo", Desc: "bar"}

	mockRepo.EXPECT().Delete(1).Return(expected, nil).Times(1)

	result, err := svc.Delete(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
