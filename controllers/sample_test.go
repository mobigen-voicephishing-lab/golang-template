package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/labstack/echo/v5"
	"github.com/mobigen/golang-web-template/controllers"
	"github.com/mobigen/golang-web-template/mocks"
	"github.com/mobigen/golang-web-template/models"
	"github.com/stretchr/testify/assert"
)

func TestSampleController_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSampleUsecase(ctrl)
	controller := controllers.Sample{}.New(mockUsecase)

	expected := &[]models.Sample{
		{ID: 1, Name: "foo", Desc: "bar"},
	}
	mockUsecase.EXPECT().GetAll().Return(expected, nil).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/samples", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GetAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result []models.Sample
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "foo", result[0].Name)
}

func TestSampleController_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSampleUsecase(ctrl)
	controller := controllers.Sample{}.New(mockUsecase)

	mockUsecase.EXPECT().GetAll().Return(nil, errors.New("db error")).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/samples", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.GetAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSampleController_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSampleUsecase(ctrl)
	controller := controllers.Sample{}.New(mockUsecase)

	expected := &models.Sample{ID: 1, Name: "foo", Desc: "bar"}
	mockUsecase.EXPECT().GetByID(1).Return(expected, nil).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/samples/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// echo v5 에서 path parameter 설정
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "1"}})

	err := controller.GetByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var result models.Sample
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "foo", result.Name)
}

func TestSampleController_GetByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSampleUsecase(ctrl)
	controller := controllers.Sample{}.New(mockUsecase)

	// 숫자가 아닌 id → Usecase 호출 없이 400 반환
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/samples/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "abc"}})

	err := controller.GetByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSampleController_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSampleUsecase(ctrl)
	controller := controllers.Sample{}.New(mockUsecase)

	input := &models.Sample{Name: "new sample", Desc: "new desc"}
	expected := &models.Sample{ID: 10, Name: "new sample", Desc: "new desc"}

	// Create 는 어떤 *models.Sample 을 받아도 expected 를 반환하도록 설정
	mockUsecase.EXPECT().Create(gomock.Any()).Return(expected, nil).Times(1)

	body, _ := json.Marshal(input)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/samples", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var result models.Sample
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, 10, result.ID)
}

func TestSampleController_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSampleUsecase(ctrl)
	controller := controllers.Sample{}.New(mockUsecase)

	expected := &models.Sample{ID: 1, Name: "foo"}
	mockUsecase.EXPECT().Delete(1).Return(expected, nil).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/samples/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "1"}})

	err := controller.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
