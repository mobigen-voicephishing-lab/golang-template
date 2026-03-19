# Web-Server Platform Sample (Golang 1.26)

Go 언어를 이용한 웹서버 샘플

## 1. 요구사항 및 의존성

1. Go 버전

    최소 **Go 1.17** 이상이 필요합니다. (`go install pkg@version` 지원 버전)
    현재 프로젝트는 **Go 1.26** 기준으로 작성되었습니다.

2. 주요 의존성 패키지

    | 패키지                              | 버전       | 사용처                                        |
    | ----------------------------------- | ---------- | --------------------------------------------- |
    | `github.com/labstack/echo/v5`       | v5.0.4     | HTTP 웹 프레임워크 (라우팅, 미들웨어)         |
    | `gorm.io/gorm`                      | v1.31.1    | ORM 프레임워크                                |
    | `gorm.io/driver/mysql`              | v1.6.0     | MySQL 데이터베이스 드라이버                   |
    | `gorm.io/driver/postgres`           | v1.6.0     | PostgreSQL 데이터베이스 드라이버              |
    | `gorm.io/driver/sqlite`             | v1.6.0     | SQLite 데이터베이스 드라이버                  |
    | `github.com/spf13/viper`            | v1.21.0    | 환경변수 및 설정 파일(YAML) 관리              |
    | `github.com/fsnotify/fsnotify`      | v1.9.0     | 설정 파일 변경 감지 (hot-reload)              |
    | `github.com/sirupsen/logrus`        | v1.9.4     | 구조화된 로깅                                 |
    | `gopkg.in/natefinch/lumberjack.v2`  | v2.2.1     | 로그 파일 로테이션                            |
    | `github.com/swaggo/swag/v2`         | v2.0.0-rc5 | Swagger 문서 자동 생성 (코드 어노테이션 파싱) |
    | `github.com/swaggo/echo-swagger/v2` | v2.0.1     | Echo v5에서 Swagger UI 제공 (OAS 2.0)         |
    | `github.com/stretchr/testify`       | v1.11.1    | 테스트 assertion 유틸리티                     |
    | `go.uber.org/mock`                  | v0.6.0     | Mock 코드 생성 (mockgen)                      |

## 2. 프로젝트 구조

### 2.1. 아키텍처: Hexagonal Architecture (Ports & Adapters)

이 프로젝트는 헥사고날 아키텍처(Ports & Adapters)를 기반으로 설계되었습니다.
도메인(비즈니스 로직)이 중심에 위치하고, 외부 시스템(HTTP, DB)은 어댑터를 통해 연결됩니다.

```text
                     ┌─────────────────────────────────────┐
                     │          internal/                  │
  HTTP Request  ───► │  adapter/inbound  ──► usecase  ──► │ ──► DB
                     │  (handler, router)     (use case)   │     (adapter/outbound)
                     │                    domain (entity)  │
                     └─────────────────────────────────────┘
```

- **Domain** (중심): 순수 비즈니스 엔티티와 도메인 에러. 프레임워크 의존성 없음.
- **UseCase** (애플리케이션 레이어): 비즈니스 로직 구현. Repository 인터페이스(포트)에만 의존.
- **Adapter Inbound** (HTTP 어댑터): Echo 핸들러. HTTP 요청을 UseCase로 전달.
- **Adapter Outbound** (Persistence 어댑터): GORM 기반 DB 구현체. Repository 인터페이스 구현.
- **Infrastructure** (기반 서비스): Config, Logger, DB 연결 등 기술적 관심사.
- **Bootstrap** (의존성 주입): 모든 레이어를 조립하고 라우트를 등록.

### 2.2. 디렉토리 구조

```text
.
├── cmd/
│   └── server/
│       └── main.go                  # 애플리케이션 진입점
├── internal/
│   ├── adapter/
│   │   ├── inbound/
│   │   │   └── http/
│   │   │       ├── dto/             # 응답 구조체, 에러 코드 정의
│   │   │       ├── handler/         # HTTP 핸들러 (Echo)
│   │   │       └── router.go        # Echo 초기화, 미들웨어, 전역 에러 핸들러
│   │   └── outbound/
│   │       └── persistence/
│   │           └── gorm/            # GORM 기반 Repository 구현체
│   ├── domain/
│   │   ├── sample.go                # 도메인 엔티티
│   │   └── errors.go                # AppError (비즈니스 에러)
│   ├── usecase/
│   │   └── sample/
│   │       ├── interface.go         # Repository 포트 (인터페이스)
│   │       ├── get.go               # GetByID 유스케이스
│   │       ├── get_all.go           # GetAll 유스케이스
│   │       ├── create.go            # Create 유스케이스
│   │       ├── update.go            # Update 유스케이스
│   │       ├── delete.go            # Delete 유스케이스
│   │       └── usecase.go           # 묶음 SampleUseCase (대안)
│   ├── infrastructure/
│   │   ├── config/                  # 환경변수, 설정 파일, 빌드 정보
│   │   ├── logger/                  # Logrus + Lumberjack 로거
│   │   └── db/                      # GORM 연결, 마이그레이션
│   ├── bootstrap/
│   │   └── wire.go                  # 의존성 주입 (DI wiring)
│   └── testutil/                    # 테스트 유틸리티
├── util/
│   └── util.go                      # 시간 유틸리티
├── test/
│   └── integration/                 # 통합 테스트
├── configs/
│   └── prod.yaml                    # 설정 파일
├── docs/
│   └── swagger/                     # swag init 으로 생성된 Swagger 문서
├── build/
│   ├── Dockerfile
│   └── bin/                         # 빌드 결과물
├── Makefile
└── go.mod
```

## 3. 개발

### 3.1. Change Module Name

코드 작성에 앞서 프로젝트(모듈)의 이름을 변경한다.

```sh
# github.com/mobigen/golang-web-template => github.com/myorg/myapp 로 변경
$ find . -type f -name "*.go" -print0 | xargs -0 sed -i 's|github.com/mobigen/golang-web-template|github.com/myorg/myapp|g'
# go.mod 도 함께 수정
$ sed -i 's|github.com/mobigen/golang-web-template|github.com/myorg/myapp|g' go.mod
```

### 3.2. Domain Entity 작성

데이터베이스, 클라이언트와 주고 받을 도메인 엔티티 정의

- `internal/domain/sample.go`

```go
package domain

type Sample struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Desc     string `json:"desc"`
    CreateAt int64  `json:"createAt"`
}
```

> GORM 태그는 도메인 엔티티에 포함하지 않습니다. GORM 전용 모델은 `internal/adapter/outbound/persistence/gorm/` 에 별도로 정의합니다.

### 3.3. Repository 인터페이스 작성 (포트 정의)

UseCase가 의존하는 Repository 인터페이스를 `internal/usecase/sample/interface.go` 에 정의합니다.

```go
package sample

//go:generate go run go.uber.org/mock/mockgen -destination=mock_repository_test.go -package=sample_test github.com/mobigen/golang-web-template/internal/usecase/sample Repository

type Repository interface {
    GetAll() (*[]domain.Sample, error)
    GetByID(int) (*domain.Sample, error)
    Create(*domain.Sample) (*domain.Sample, error)
    Update(*domain.Sample) (*domain.Sample, error)
    Delete(int) (*domain.Sample, error)
}
```

### 3.4. UseCase 작성

각 CRUD 동작을 개별 UseCase 파일로 작성합니다.

- `internal/usecase/sample/get.go` — GetByID 유스케이스 예시:

```go
package sample

type GetByIDUseCase struct {
    repo Repository
}

func NewGetByIDUseCase(repo Repository) *GetByIDUseCase {
    return &GetByIDUseCase{repo: repo}
}

func (uc *GetByIDUseCase) Execute(id int) (*domain.Sample, error) {
    return uc.repo.GetByID(id)
}
```

### 3.5. Outbound Adapter 작성 (Repository 구현체)

GORM 전용 모델과 도메인 엔티티를 분리하여 구현합니다.

- `internal/adapter/outbound/persistence/gorm/sample_repository.go`

```go
package gorm

// GORM 전용 모델 (domain.Sample과 분리)
type sampleModel struct {
    ID       int    `gorm:"column:id;primaryKey;autoIncrement"`
    Name     string `gorm:"column:name;not null"`
    Desc     string `gorm:"column:desc;size:256"`
    CreateAt int64  `gorm:"column:createAt;autoCreateTime:milli"`
}

func (sampleModel) TableName() string { return "samples" }

// SampleRepository GORM 기반 구현체
type SampleRepository struct {
    ds *db.DataStore
}

func NewSampleRepository(ds *db.DataStore) *SampleRepository {
    return &SampleRepository{ds: ds}
}
```

### 3.6. 테이블 생성 (Migration)

`cmd/server/main.go` 의 `InitDatastore()` 에서 GORM 모델을 파라미터로 전달합니다.

```go
if err := ds.Migrate(persistence.SampleModel()); err != nil {
    return err
}
```

### 3.7. Inbound Handler 작성

HTTP 핸들러는 `SampleUsecase` 인터페이스에 의존합니다.

- `internal/adapter/inbound/http/handler/sample_handler.go`

```go
package handler

//go:generate go run go.uber.org/mock/mockgen -destination=mock_usecase_test.go -package=handler_test github.com/mobigen/golang-web-template/internal/adapter/inbound/http/handler SampleUsecase

type SampleUsecase interface {
    GetAll() (*[]domain.Sample, error)
    GetByID(int) (*domain.Sample, error)
    Create(*domain.Sample) (*domain.Sample, error)
    Update(*domain.Sample) (*domain.Sample, error)
    Delete(int) (*domain.Sample, error)
}

type SampleHandler struct {
    Usecase SampleUsecase
}

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
```

### 3.8. Bootstrap (DI Wiring) 등록

- `internal/bootstrap/wire.go` — 레이어를 조립하고 라우트를 등록합니다.

```go
func (in *Injector) Init() error {
    // Version
    ver := handler.NewVersionHandler()
    in.Router.GET("/version", ver.GetVersion)

    apiv1 := in.Router.Group("/api/v1")

    // Repository → 개별 UseCase → Handler 순으로 wiring
    repo := persistence.NewSampleRepository(in.Datastore)

    getAllUC  := sample.NewGetAllUseCase(repo)
    getByIDUC := sample.NewGetByIDUseCase(repo)
    createUC  := sample.NewCreateUseCase(repo)
    updateUC  := sample.NewUpdateUseCase(repo)
    deleteUC  := sample.NewDeleteUseCase(repo)

    sampleHandler := handler.NewSampleHandler(getAllUC, getByIDUC, createUC, updateUC, deleteUC)

    apiv1.GET("/samples",          sampleHandler.GetAll)
    apiv1.GET("/sample/:id",       sampleHandler.GetByID)
    apiv1.POST("/sample",          sampleHandler.Create)
    apiv1.POST("/sample/update",   sampleHandler.Update)
    apiv1.DELETE("/sample/:id",    sampleHandler.Delete)
    return nil
}
```

> **묶음 방식 (대안)**: 모든 CRUD를 하나의 `SampleUseCase` 구조체에 담는 방식도 지원합니다.
> `wire.go` 내 주석 처리된 코드를 참고하세요.

### 3.9. 에러 처리

도메인 레이어에서 `AppError`를 생성하고 핸들러까지 전파합니다.
라우터의 전역 에러 핸들러가 `AppError`를 자동으로 처리합니다.

```go
// domain/errors.go
func NewNotFoundError(msg string) AppError {
    return AppError{HttpStatus: 404, Code: 1001, Message: msg}
}

// 핸들러에서 처리
if ae, ok := err.(domain.AppError); ok {
    return FailApp(c, ae)  // AppError의 HttpStatus, Code, Message를 응답에 반영
}
```

## 4. 빌드

Makefile을 이용해 코드 정적 분석, 빌드 타임 변수 설정, 테스트, Docker 이미지 빌드를 처리한다.

### 4.1. Makefile

프로젝트 시작 시 수정할 항목

```makefile
# 바이너리 이름이자 이미지 이름
TARGET := test
# 버전 정보
VERSION := v1.0.0
# 컨테이너 레지스트리 경로
REPO := repo.iris.tools/template/
# 최종 이미지 이름 (REPO + TARGET + VERSION)
IMAGE ?= $(REPO)$(TARGET):$(VERSION)
```

`TARGET`, `VERSION`, `BUILD_HASH`는 빌드 타임에 `internal/infrastructure/config` 패키지 변수로 주입된다(`ldflags -X` 옵션).

### 4.2. 명령어

1. Lint

    코드 정적 분석

    ```sh
    make lint
    ```

2. Build

    | 명령어                    | 설명                              |
    | ------------------------- | --------------------------------- |
    | `make build`              | 현재 플랫폼용 빌드                |
    | `make build-darwin-arm64` | macOS Apple Silicon(arm64)용 빌드 |
    | `make build-linux-amd64`  | Linux amd64용 빌드                |

    빌드 결과물은 `build/bin/` 디렉토리에 생성된다.

    ```sh
    make build
    ```

3. Mock 생성

    ```sh
    # mock 생성 (mockgen 자동 설치)
    make mocks

    # mock 최신 상태 검증
    make verify-mocks
    ```

4. Test / Coverage

    ```sh
    # 테스트 실행
    make test

    # 커버리지 리포트 생성 (build/cov-out.html)
    make coverage
    ```

5. Docker 이미지 빌드

    linux/amd64 이미지 빌드 (macOS에서 실행해도 결과 이미지는 Linux amd64용)

    ```sh
    make image
    ```

6. Clean

    ```sh
    make clean
    ```

## 5. 테스트

**go.uber.org/mock** 프레임워크를 이용한 mock 생성과 유닛 테스트에 대해서 설명한다.

- interface로 선언된 변수에만 mock type을 할당할 수 있다.
- `//go:generate` 어노테이션은 인터페이스가 정의된 파일에 직접 추가한다.

### 5.1. mock 위치

이 프로젝트에서 mock은 **테스트 파일과 같은 패키지** 내에 인라인으로 생성됩니다.

```
internal/usecase/sample/
  interface.go              ← //go:generate 어노테이션 포함
  mock_repository_test.go   ← 생성된 mock (sample_test 패키지)

internal/adapter/inbound/http/handler/
  sample_handler.go         ← //go:generate 어노테이션 포함
  mock_usecase_test.go      ← 생성된 mock (handler_test 패키지)
```

### 5.2. `//go:generate` 지시어

mock을 생성할 인터페이스가 있는 파일에 `//go:generate` 주석을 추가합니다.

**`internal/usecase/sample/interface.go`** — `Repository` 인터페이스에 추가

```go
package sample

//go:generate go run go.uber.org/mock/mockgen -destination=mock_repository_test.go -package=sample_test github.com/mobigen/golang-web-template/internal/usecase/sample Repository

type Repository interface {
    GetAll() (*[]domain.Sample, error)
    // ...
}
```

**`internal/adapter/inbound/http/handler/sample_handler.go`** — `SampleUsecase` 인터페이스에 추가

```go
package handler

//go:generate go run go.uber.org/mock/mockgen -destination=mock_usecase_test.go -package=handler_test github.com/mobigen/golang-web-template/internal/adapter/inbound/http/handler SampleUsecase

type SampleUsecase interface {
    GetAll() (*[]domain.Sample, error)
    // ...
}
```

### 5.3. mock 생성 (`make mocks`)

```sh
make mocks
```

이 명령은 내부적으로 `go generate ./...` 를 실행하여 `//go:generate` 지시어가 있는 모든 파일에서 mock을 생성합니다.

### 5.4. 유닛 테스트 작성

**UseCase 레이어 테스트** (`internal/usecase/sample/get_test.go`)

UseCase는 `Repository` 인터페이스에 의존하므로, `MockRepository`로 교체하여 테스트합니다.

```go
package sample_test

import (
    "testing"
    "go.uber.org/mock/gomock"
    "github.com/mobigen/golang-web-template/internal/domain"
    "github.com/mobigen/golang-web-template/internal/usecase/sample"
    "github.com/stretchr/testify/assert"
)

func TestGetByIDUseCase_Execute(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockRepository(ctrl)
    uc := sample.NewGetByIDUseCase(mockRepo)

    expected := &domain.Sample{ID: 1, Name: "foo", Desc: "bar"}
    mockRepo.EXPECT().GetByID(1).Return(expected, nil).Times(1)

    result, err := uc.Execute(1)

    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

**Handler 레이어 테스트** (`internal/adapter/inbound/http/handler/sample_handler_test.go`)

Handler는 `SampleUsecase` 인터페이스에 의존하므로, `MockSampleUsecase`로 교체합니다.
HTTP 요청/응답은 `net/http/httptest`와 echo를 사용합니다.

```go
package handler_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "go.uber.org/mock/gomock"
    "github.com/labstack/echo/v5"
    "github.com/mobigen/golang-web-template/internal/adapter/inbound/http/handler"
    "github.com/mobigen/golang-web-template/internal/domain"
    "github.com/stretchr/testify/assert"
)

func TestSampleHandler_GetByID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUsecase := NewMockSampleUsecase(ctrl)
    h := handler.NewSampleHandlerFromUsecase(mockUsecase)

    expected := &domain.Sample{ID: 1, Name: "foo", Desc: "bar"}
    mockUsecase.EXPECT().GetByID(1).Return(expected, nil).Times(1)

    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/sample/1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetPathValues(echo.PathValues{{Name: "id", Value: "1"}})

    err := h.GetByID(c)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

> **gomock 주요 패턴**
>
> | 패턴                | 설명                     |
> | ------------------- | ------------------------ |
> | `.Times(1)`         | 정확히 1번 호출되어야 함 |
> | `.AnyTimes()`       | 0번 이상 호출 가능       |
> | `gomock.Any()`      | 어떤 인자값이어도 매칭   |
> | `.Return(val, nil)` | 반환값 지정              |

### 5.5. 테스트 실행 (`make test`)

```sh
make test
```

커버리지 리포트 생성:

```sh
make coverage
# build/cov-out.html 파일로 브라우저에서 확인 가능
```

## 6. 문서작성(Swagger)

API 문서를 코드 개발 단계에서 `// @` 형식의 주석으로 작성하고 `make swag` 명령으로 자동 생성합니다.
**코드 주석이 곧 API 문서**입니다.

- 런타임 라이브러리: `github.com/swaggo/swag/v2`, `github.com/swaggo/echo-swagger/v2`
- 문서 생성 CLI: `github.com/swaggo/swag/v2/cmd/swag` (Makefile에서 자동 설치)
- Swagger UI: `http://localhost:8080/swagger/index.html`
- 참고: [swaggo/swag](https://github.com/swaggo/swag), [swaggo/echo-swagger](https://github.com/swaggo/echo-swagger)

### 6.1. 서버 등록

`internal/adapter/inbound/http/router.go` 에서 이미 등록되어 있습니다.

```go
import (
    _ "github.com/mobigen/golang-web-template/docs/swagger"   // swag init으로 생성된 docs 패키지
    echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// Swagger UI (OAS 2.0)
// WrapHandlerV3: swag/v2 레지스트리에서 doc.json을 읽는다.
// WrapHandler:   swag/v1 레지스트리를 읽으므로 사용하지 않는다.
r.GET("/swagger/*", echoSwagger.WrapHandlerV3)
```

### 6.2. 작성 방법

**① 서버 전체 정보** — `cmd/server/main.go` 의 `main()` 함수 위에 작성

```go
// @title Golang Web Template API
// @version 1.0.0
// @description This is a golang web template server.

// @contact.name API Support
// @contact.url http://mobigen.com
// @contact.email jblim@mobigen.com

// @host localhost:8080
// @BasePath /
func main() {
```

**② 개별 API** — 각 handler 함수 위에 작성 (`handler/version_handler.go` 예시)

swag/v2 는 Go 제네릭 문법 `[TypeParam]` 을 어노테이션에서 그대로 사용할 수 있습니다.

```go
// GetVersion return app version
// @Summary Get Server Version
// @Description get server version info
// @Tags version
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.HTTPResponse[config.VersionInfo] "app info(name, version, hash)"
// @Router /version [get]
func (h *VersionHandler) GetVersion(c *echo.Context) error {
    ...
}
```

**③ 응답 타입 어노테이션 패턴**

```go
// 단일 객체
// @Success 200 {object} dto.HTTPResponse[domain.Sample] "샘플"

// 배열
// @Success 200 {object} dto.HTTPResponse[[]domain.Sample] "샘플 목록"

// 에러 (데이터 없음)
// @Failure 400 {object} dto.HTTPResponse[any] "잘못된 요청"
// @Failure 500 {object} dto.HTTPResponse[any] "서버 오류"
```

주요 어노테이션:

| 어노테이션     | 설명                                                 |
| -------------- | ---------------------------------------------------- |
| `@Summary`     | API 요약 (목록에 표시)                               |
| `@Description` | 상세 설명                                            |
| `@Tags`        | API 그룹 분류                                        |
| `@Accept`      | 요청 Content-Type (`json`, `multipart/form-data` 등) |
| `@Produce`     | 응답 Content-Type                                    |
| `@Param`       | 파라미터 정의 (이름, 위치, 타입, 필수여부, 설명)     |
| `@Success`     | 성공 응답 코드 및 스키마                             |
| `@Failure`     | 실패 응답 코드 및 스키마                             |
| `@Router`      | 경로 및 HTTP 메서드                                  |

### 6.3. 모델 이름 단축 (`// @name`)

기본적으로 swag는 패키지 경로 전체를 모델 이름으로 사용합니다.
구조체 정의의 **닫는 괄호** 뒤에 `// @name <별칭>` 을 추가하면 Swagger UI에 표시되는 이름을 단축할 수 있습니다.

```go
// 적용 전: github_com_mobigen_golang-web-template_internal_domain.Sample
// 적용 후: Sample
type Sample struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
} // @name Sample

// 제네릭 타입 예시 — 인스턴스화 시 "HTTPResponse-Sample" 형태로 표시됨
type HTTPResponse[T any] struct {
    IsSuccess bool   `json:"isSuccess"`
    Code      int    `json:"code"`
    Message   string `json:"message"`
    Data      T      `json:"data,omitempty"`
} // @name HTTPResponse
```

> **주의**: `type Foo struct { // @name Foo` (여는 괄호 뒤)나 `// @name Foo\ntype Foo struct` (doc comment)에 작성하면 동작하지 않습니다. 반드시 **닫는 괄호 뒤**에 작성해야 합니다.

### 6.4. 문서 생성

```sh
make swag
```

외부 패키지 타입을 참조할 때 파싱 오류가 발생하면 `--parseDependency --parseInternal` 플래그를 사용합니다 (Makefile에 이미 포함).

생성 결과: `docs/swagger/` 디렉토리

```txt
docs/swagger/
  docs.go        ← Go 패키지 (서버 import 용, swag/v2 레지스트리에 등록)
  swagger.json   ← Swagger 명세 (OAS 2.0)
  swagger.yaml   ← Swagger 명세 (YAML 형식)
```

### 6.5. 문서 확인

서버 실행 후 브라우저에서 접속합니다.

`http://localhost:8080/swagger/index.html`

## 7. 실행

### 7.1. 환경 변수

프로그램 시작 시 아래 환경 변수를 읽어 동작을 제어합니다.

| 환경 변수   | 기본값             | 설명                                                            |
| ----------- | ------------------ | --------------------------------------------------------------- |
| `APP_HOME`  | 바이너리 실행 경로 | 앱 홈 디렉토리. config, DB 파일 경로의 기준이 됩니다.           |
| `PROFILE`   | `prod`             | 로드할 config 파일 이름. `configs/$PROFILE.yaml` 을 읽습니다.   |
| `LOG_LEVEL` | config 파일 설정   | 로그 레벨 강제 지정. `debug`, `info`, `warn`, `error`, `silent` |

### 7.2. 설정 파일

config 파일 경로: `$APP_HOME/configs/$PROFILE.yaml`

```yaml
log:
  output: "both"         # stdout, file, both
  level: "debug"         # debug, info, warn, error, silent
  # 파일 출력 시 추가 옵션 (output: "file" 또는 "both" 일 때 유효)
  logRotate:
    savePath: "logs"       # $APP_HOME 기준 상대경로 또는 절대경로
    fileName: "app.log"    # 로그 파일의 이름, 백업 시 app-{2026-01-01T14:15:12.000}.log
    sizePerFileMb: 100     # 로그 파일 최대 크기 (MB)
    maxOfDay: 10           # 보관할 백업 파일 수
    maxAge: 7              # 로그 파일 보관 기간 (일)
    compress: false        # 오래된 로그 파일 gzip 압축 여부
datastore:
  database: "sqlite3"    # mysql, postgres, sqlite3
  endPoint:
    path: "db/store.db"  # sqlite3: $APP_HOME 기준 상대경로
    # host: "1.2.3.4"    # mysql, postgres
    # port: 3306
    # user: "db_user"
    # pass: "db_pass"
    # dbName: "db_name"
    # option: "charset=utf8mb4&parseTime=True"
  connPool:
    maxIdleConns: 1      # 유휴 커넥션 수 (sqlite3 미사용)
    maxOpenConns: 2      # 최대 커넥션 수 (sqlite3 미사용)
  debug:
    logLevel: "info"     # silent, error, warn, info
    slowThreshold: "1sec"  # 슬로우 쿼리 기준 1min, 1sec, 1ms

server:
  debug: true            # echo 프레임워크 요청 로그 활성화
  host: "0.0.0.0"        # listen address
  port: 8080             # listen port
```

### 7.3. 실행

**개발 환경 (go run)**:

```sh
# 기본 실행 (prod 프로파일, config 파일의 로그 레벨 사용)
go run cmd/server/main.go

# 환경 변수 지정
APP_HOME=$(pwd) PROFILE=prod LOG_LEVEL=debug go run cmd/server/main.go
```

**빌드 후 실행**:

```sh
make build

# 실행
APP_HOME=$(pwd) ./build/bin/test
```

**정상 실행 로그 예시**:

```log
APP_HOME : /home/jblim/workspace/golang-web-template
PROFILE : prod
[ Env ] Read ...................................................................... [ OK ]
[ Configuration ] Read ............................................................ [ OK ]
==========================================================================================

                         START. test:v1.0.0-9ff1e66

                                                    Copyright(C) 2026 Mobigen Corporation.

==========================================================================================
[ DataStore ] Initialze ........................................................... [ OK ]
[ Router ] Initialze .............................................................. [ OK ]
[ ALL ] Initialze ................................................................. [ OK ]
[ Signal ] Listener Start ......................................................... [ OK ]
[ Router ] Listener Start ......................................................... [ OK ]
```

#### 종료

`Ctrl+C` 또는 `SIGTERM` 시그널을 보내면 graceful shutdown 됩니다.

```log
[ SIGNAL ] Receive [ terminated ]
[ DataStore ] Shutdown ............................................................ [ OK ]
[ Router ] Shutdown ............................................................... [ OK ]
```
