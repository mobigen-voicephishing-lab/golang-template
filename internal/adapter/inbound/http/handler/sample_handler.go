package handler

//go:generate go run go.uber.org/mock/mockgen -destination=mock_usecase_test.go -package=handler_test github.com/mobigen/golang-web-template/internal/adapter/inbound/http/handler SampleUsecase

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/sirupsen/logrus"

	"github.com/mobigen/golang-web-template/internal/adapter/inbound/http/dto"
	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/infrastructure/logger"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
)

// SampleUsecase handler가 의존하는 유스케이스 인터페이스
type SampleUsecase interface {
	GetAll() (*[]domain.Sample, error)
	GetByID(int) (*domain.Sample, error)
	Create(*domain.Sample) (*domain.Sample, error)
	Update(*domain.Sample) (*domain.Sample, error)
	Delete(int) (*domain.Sample, error)
}

// sampleUsecaseAdapter 메서드별 UseCase를 SampleUsecase 인터페이스로 변환하는 어댑터
type sampleUsecaseAdapter struct {
	getAll  *sample.GetAllUseCase
	getByID *sample.GetByIDUseCase
	create  *sample.CreateUseCase
	update  *sample.UpdateUseCase
	delete  *sample.DeleteUseCase
}

func (a *sampleUsecaseAdapter) GetAll() (*[]domain.Sample, error) {
	return a.getAll.Execute()
}

func (a *sampleUsecaseAdapter) GetByID(id int) (*domain.Sample, error) {
	return a.getByID.Execute(id)
}

func (a *sampleUsecaseAdapter) Create(s *domain.Sample) (*domain.Sample, error) {
	return a.create.Execute(s)
}

func (a *sampleUsecaseAdapter) Update(s *domain.Sample) (*domain.Sample, error) {
	return a.update.Execute(s)
}

func (a *sampleUsecaseAdapter) Delete(id int) (*domain.Sample, error) {
	return a.delete.Execute(id)
}

// SampleHandler Sample HTTP 핸들러
type SampleHandler struct {
	Log     *logrus.Logger
	Usecase SampleUsecase
}

// NewSampleHandler 메서드별 UseCase를 받아 SampleHandler 생성 (wire에서 사용)
func NewSampleHandler(
	getAll *sample.GetAllUseCase,
	getByID *sample.GetByIDUseCase,
	create *sample.CreateUseCase,
	update *sample.UpdateUseCase,
	delete *sample.DeleteUseCase,
) *SampleHandler {
	return &SampleHandler{
		Log: logger.Logger{}.GetInstance().Logger,
		Usecase: &sampleUsecaseAdapter{
			getAll:  getAll,
			getByID: getByID,
			create:  create,
			update:  update,
			delete:  delete,
		},
	}
}

// NewSampleHandlerFromUsecase 묶음 방식 UseCase를 받아 SampleHandler 생성 (대안)
func NewSampleHandlerFromUsecase(usecase SampleUsecase) *SampleHandler {
	return &SampleHandler{
		Log:     logger.Logger{}.GetInstance().Logger,
		Usecase: usecase,
	}
}

// GetAll returns all of sample as JSON object.
// @Summary Get all samples
// @Description 전체 샘플 목록을 반환한다
// @Tags sample
// @Accept json
// @Produce json
// @Success 200 {object} dto.HTTPResponse[[]domain.Sample] "샘플 목록"
// @Failure 500 {object} dto.HTTPResponse[any] "서버 오류"
// @Router /api/v1/samples [get]
func (h *SampleHandler) GetAll(c *echo.Context) error {
	samples, err := h.Usecase.GetAll()
	if err != nil {
		if ae, ok := err.(domain.AppError); ok {
			return FailApp(c, ae)
		}
		return Fail(c, http.StatusInternalServerError, dto.ErrInternalServer, "")
	}
	return OK(c, samples)
}

// GetByID return sample whose ID matches
// @Summary Get sample by ID
// @Description ID로 샘플을 조회한다
// @Tags sample
// @Accept json
// @Produce json
// @Param id path int true "Sample ID"
// @Success 200 {object} dto.HTTPResponse[domain.Sample] "샘플"
// @Failure 400 {object} dto.HTTPResponse[any] "잘못된 파라미터"
// @Failure 404 {object} dto.HTTPResponse[any] "샘플 없음"
// @Failure 500 {object} dto.HTTPResponse[any] "서버 오류"
// @Router /api/v1/sample/{id} [get]
func (h *SampleHandler) GetByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Fail(c, http.StatusBadRequest, dto.ErrInvalidParameter, "id must be a number")
	}
	s, err := h.Usecase.GetByID(id)
	if err != nil {
		if ae, ok := err.(domain.AppError); ok {
			return FailApp(c, ae)
		}
		return Fail(c, http.StatusInternalServerError, dto.ErrInternalServer, "")
	}
	return OK(c, s)
}

// Create create a new sample
// @Summary Create sample
// @Description 새 샘플을 생성한다
// @Tags sample
// @Accept json
// @Produce json
// @Param sample body domain.Sample true "생성할 샘플"
// @Success 200 {object} dto.HTTPResponse[domain.Sample] "생성된 샘플"
// @Failure 400 {object} dto.HTTPResponse[any] "잘못된 요청"
// @Failure 500 {object} dto.HTTPResponse[any] "서버 오류"
// @Router /api/v1/sample [post]
func (h *SampleHandler) Create(c *echo.Context) error {
	input := new(domain.Sample)
	if err := c.Bind(input); err != nil {
		return Fail(c, http.StatusBadRequest, dto.ErrInvalidRequestBody, "")
	}
	s, err := h.Usecase.Create(input)
	if err != nil {
		if ae, ok := err.(domain.AppError); ok {
			return FailApp(c, ae)
		}
		return Fail(c, http.StatusInternalServerError, dto.ErrInternalServer, "")
	}
	return OK(c, s)
}

// Update update from input
// @Summary Update sample
// @Description 샘플을 수정한다
// @Tags sample
// @Accept json
// @Produce json
// @Param sample body domain.Sample true "수정할 샘플 (ID 필수)"
// @Success 200 {object} dto.HTTPResponse[domain.Sample] "수정된 샘플"
// @Failure 400 {object} dto.HTTPResponse[any] "잘못된 요청"
// @Failure 500 {object} dto.HTTPResponse[any] "서버 오류"
// @Router /api/v1/sample/update [post]
func (h *SampleHandler) Update(c *echo.Context) error {
	input := new(domain.Sample)
	if err := c.Bind(input); err != nil {
		return Fail(c, http.StatusBadRequest, dto.ErrInvalidRequestBody, "")
	}
	s, err := h.Usecase.Update(input)
	if err != nil {
		if ae, ok := err.(domain.AppError); ok {
			return FailApp(c, ae)
		}
		return Fail(c, http.StatusInternalServerError, dto.ErrInternalServer, "")
	}
	return OK(c, s)
}

// Delete delete sample from id
// @Summary Delete sample
// @Description ID로 샘플을 삭제한다
// @Tags sample
// @Accept json
// @Produce json
// @Param id path int true "Sample ID"
// @Success 200 {object} dto.HTTPResponse[domain.Sample] "삭제된 샘플"
// @Failure 400 {object} dto.HTTPResponse[any] "잘못된 파라미터"
// @Failure 404 {object} dto.HTTPResponse[any] "샘플 없음"
// @Failure 500 {object} dto.HTTPResponse[any] "서버 오류"
// @Router /api/v1/sample/{id} [delete]
func (h *SampleHandler) Delete(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Fail(c, http.StatusBadRequest, dto.ErrInvalidParameter, "id must be a number")
	}
	s, err := h.Usecase.Delete(id)
	if err != nil {
		if ae, ok := err.(domain.AppError); ok {
			return FailApp(c, ae)
		}
		return Fail(c, http.StatusInternalServerError, dto.ErrInternalServer, "")
	}
	return OK(c, s)
}
