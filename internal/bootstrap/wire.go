package bootstrap

import (
	apphttp "github.com/mobigen/golang-web-template/internal/adapter/inbound/http"
	"github.com/mobigen/golang-web-template/internal/adapter/inbound/http/handler"
	persistence "github.com/mobigen/golang-web-template/internal/adapter/outbound/persistence/gorm"
	"github.com/mobigen/golang-web-template/internal/infrastructure/db"
	"github.com/mobigen/golang-web-template/internal/infrastructure/logger"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
)

// Injector web-server layer initializer : Dependency Injection
type Injector struct {
	Router    *apphttp.Router
	Datastore *db.DataStore
	Log       *logger.LogrusLogger
}

// New create Injector
func (Injector) New(r *apphttp.Router, d *db.DataStore, l *logger.LogrusLogger) *Injector {
	return &Injector{
		Router:    r,
		Datastore: d,
		Log:       l,
	}
}

// Init init web-server layer interconnection (Dependency Injection)
func (in *Injector) Init() error {
	// ── Version Handler ──
	ver := handler.NewVersionHandler()
	in.Router.GET("/version", ver.GetVersion)

	// path grouping
	apiv1 := in.Router.Group("/api/v1")

	// ── Sample: 메서드별 UseCase 방식 (실제 사용) ──
	// Repository → 개별 UseCase → Handler 순서로 wiring
	in.Log.Errorf("[ PATH ] /api/v1/sample ........................................................... [ OK ]")
	repo := persistence.NewSampleRepository(in.Datastore)

	getAllUC := sample.NewGetAllUseCase(repo)
	getByIDUC := sample.NewGetByIDUseCase(repo)
	createUC := sample.NewCreateUseCase(repo)
	updateUC := sample.NewUpdateUseCase(repo)
	deleteUC := sample.NewDeleteUseCase(repo)

	sampleHandler := handler.NewSampleHandler(in.Log, getAllUC, getByIDUC, createUC, updateUC, deleteUC)

	// ── 묶음 방식 (대안) ──
	// 모든 메서드를 하나의 SampleUseCase 구조체에 담는 방식입니다.
	// SampleUseCase는 Repository 인터페이스에 의존하며, 모든 CRUD 메서드를 포함합니다.
	//
	// uc := sample.NewSampleUseCase(repo)
	// sampleHandler := handler.NewSampleHandlerWithUsecase(uc)

	apiv1.GET("/samples", sampleHandler.GetAll)
	apiv1.GET("/sample/:id", sampleHandler.GetByID)
	apiv1.POST("/sample", sampleHandler.Create)
	apiv1.POST("/sample/update", sampleHandler.Update)
	apiv1.DELETE("/sample/:id", sampleHandler.Delete)

	return nil
}
