package sample_test

import (
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	uc := sample.NewDeleteUseCase(mockRepo)

	expected := &domain.Sample{ID: 1, Name: "foo", Desc: "bar"}
	mockRepo.EXPECT().Delete(1).Return(expected, nil).Times(1)

	result, err := uc.Execute(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
