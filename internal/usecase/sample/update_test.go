package sample_test

import (
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewUpdateUseCase(mockRepo)

	input := &domain.Sample{ID: 1, Name: "updated", Desc: "updated desc"}
	mockRepo.EXPECT().Update(input).Return(input, nil).Times(1)

	result, err := uc.Execute(input)

	assert.NoError(t, err)
	assert.Equal(t, input, result)
}
