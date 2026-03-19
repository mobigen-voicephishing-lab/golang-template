package sample_test

import (
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
	"github.com/stretchr/testify/assert"
)

func TestCreateUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewCreateUseCase(mockRepo)

	input := &domain.Sample{Name: "new sample", Desc: "new desc"}
	expected := &domain.Sample{ID: 10, Name: "new sample", Desc: "new desc"}
	mockRepo.EXPECT().Create(input).Return(expected, nil).Times(1)

	result, err := uc.Execute(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
