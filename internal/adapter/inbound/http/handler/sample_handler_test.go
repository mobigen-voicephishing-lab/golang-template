package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/labstack/echo/v5"
	"github.com/mobigen/golang-web-template/internal/adapter/inbound/http/dto"
	"github.com/mobigen/golang-web-template/internal/adapter/inbound/http/handler"
	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/stretchr/testify/assert"
)

func newTestHandler(mockUsecase handler.SampleUsecase) *handler.SampleHandler {
	return handler.NewSampleHandlerFromUsecase(mockUsecase)
}

func TestSampleHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockSampleUsecase(ctrl)
	h := newTestHandler(mockUsecase)

	expected := &[]domain.Sample{
		{ID: 1, Name: "foo", Desc: "bar"},
	}
	mockUsecase.EXPECT().GetAll().Return(expected, nil).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/samples", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		IsSuccess bool            `json:"isSuccess"`
		Code      int             `json:"code"`
		Data      []domain.Sample `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.True(t, resp.IsSuccess)
	assert.Equal(t, dto.Success, resp.Code)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "foo", resp.Data[0].Name)
}

func TestSampleHandler_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockSampleUsecase(ctrl)
	h := newTestHandler(mockUsecase)

	mockUsecase.EXPECT().GetAll().Return(nil, errors.New("db error")).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/samples", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp struct {
		IsSuccess bool `json:"isSuccess"`
		Code      int  `json:"code"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.False(t, resp.IsSuccess)
	assert.Equal(t, dto.ErrInternalServer, resp.Code)
}

func TestSampleHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockSampleUsecase(ctrl)
	h := newTestHandler(mockUsecase)

	expected := &domain.Sample{ID: 1, Name: "foo", Desc: "bar"}
	mockUsecase.EXPECT().GetByID(1).Return(expected, nil).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/samples/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "1"}})

	err := h.GetByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		IsSuccess bool          `json:"isSuccess"`
		Data      domain.Sample `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.True(t, resp.IsSuccess)
	assert.Equal(t, 1, resp.Data.ID)
	assert.Equal(t, "foo", resp.Data.Name)
}

func TestSampleHandler_GetByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockSampleUsecase(ctrl)
	h := newTestHandler(mockUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/samples/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "abc"}})

	err := h.GetByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp struct {
		IsSuccess bool `json:"isSuccess"`
		Code      int  `json:"code"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.False(t, resp.IsSuccess)
	assert.Equal(t, dto.ErrInvalidParameter, resp.Code)
}

func TestSampleHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockSampleUsecase(ctrl)
	h := newTestHandler(mockUsecase)

	input := &domain.Sample{Name: "new sample", Desc: "new desc"}
	expected := &domain.Sample{ID: 10, Name: "new sample", Desc: "new desc"}

	mockUsecase.EXPECT().Create(gomock.Any()).Return(expected, nil).Times(1)

	body, _ := json.Marshal(input)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/samples", strings.NewReader(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		IsSuccess bool          `json:"isSuccess"`
		Data      domain.Sample `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.True(t, resp.IsSuccess)
	assert.Equal(t, 10, resp.Data.ID)
}

func TestSampleHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockSampleUsecase(ctrl)
	h := newTestHandler(mockUsecase)

	expected := &domain.Sample{ID: 1, Name: "foo"}
	mockUsecase.EXPECT().Delete(1).Return(expected, nil).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/samples/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "1"}})

	err := h.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		IsSuccess bool          `json:"isSuccess"`
		Data      domain.Sample `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.True(t, resp.IsSuccess)
	assert.Equal(t, 1, resp.Data.ID)
}
