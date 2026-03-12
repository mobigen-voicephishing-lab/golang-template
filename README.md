# Web-Server Platform Sample (Golang 1.26)

Go 언어를 이용한 웹서버 샘플

## 1. 요구사항 및 의존성

1. Go 버전

    최소 **Go 1.17** 이상이 필요합니다. (`go install pkg@version` 지원 버전)
    현재 프로젝트는 **Go 1.26** 기준으로 작성되었습니다.

2. 주요 의존성 패키지

    | 패키지                           | 버전                  | 사용처                                        |
    | -------------------------------- | --------------------- | --------------------------------------------- |
    | `github.com/labstack/echo/v5`    | v5.0.4                | HTTP 웹 프레임워크 (라우팅, 미들웨어)         |
    | `gorm.io/gorm`                   | v1.31.1               | ORM 프레임워크                                |
    | `gorm.io/driver/mysql`           | v1.6.0                | MySQL 데이터베이스 드라이버                   |
    | `gorm.io/driver/postgres`        | v1.6.0                | PostgreSQL 데이터베이스 드라이버              |
    | `gorm.io/driver/sqlite`          | v1.6.0                | SQLite 데이터베이스 드라이버                  |
    | `github.com/spf13/viper`         | v1.21.0               | 환경변수 및 설정 파일(YAML) 관리              |
    | `github.com/fsnotify/fsnotify`   | v1.9.0                | 설정 파일 변경 감지 (hot-reload)              |
    | `github.com/mobigen/gologger`    | v1.1.1                | 커스텀 로거                                   |
    | `github.com/sirupsen/logrus`     | v1.9.4                | 구조화된 로깅                                 |
    | `github.com/swaggo/swag`         | v1.16.6               | Swagger 문서 자동 생성 (코드 어노테이션 파싱) |
    | `github.com/swaggo/echo-swagger` | v1.5.0                | Echo에서 Swagger UI 제공                      |
    | `github.com/alecthomas/template` | v0.0.0-20190718012654 | Swagger HTML 템플릿 렌더링                    |
    | `github.com/stretchr/testify`    | v1.11.1               | 테스트 assertion 유틸리티                     |

## 2. 프로젝트 구조

설계 바탕이 된 레이어 구조

- Layered Architecture

레이어라는 말에서 알 수 있듯 각 레이어(계층)는 동일한 관심사(역할)의 집합으로 이루어져 있고,
겹겹이 쌓아 올린 구조라는 것을 알 수 있을 것 입니다. 그리고 이 계층 구조가 의미하는
다른 부분을 생각해 본다면 계층 간 접근/제어에 대한 부분으로 하위(종속) 계층으로만
컨트롤될 수 있도록 하는 설계 단계에서부터의 제한이라고 볼 수 있습니다.
예를 들어 상위 계층에서 하위 계층으로 접근은 가능하지만 하위 계층에서 상위 계층으로의
조작 코드는 작성하면 안 됩니다.

각 레이어를 분리하고 각 레이어의 역할을 구체화해 규칙을 준수한다면,
기능/사람의 변경에도 유연하게 대처할 수 있어 어려움(비용)이 줄어들 수 있을 것 입니다.

- Software Architecture is The Art Of Drawing Lines That I Call Boundaries. - 로버트C. 마틴

소프트웨어 아키텍처는 경계라는 선을 그리는 예술이다. 이러한 경계들은 소프트웨어 요소들을 서로
분리하고 디펜던시 의존성을 제한한다. 아키텍트의 목표는 필요한 시스템을 구축하고 유지하는 데
필요한 인적 리소스를 최소화하는 것이다. 예를 들어서, 비즈니스 유스케이스와 데이터베이스 사이의
경계선을 그릴 수 있다. 그 선은 비즈니스 규칙이 데이터베이스에 대해 전혀 알지 못하도록 막았다.
그 결정은 데이터베이스의 선택과 실행을 뒤로 늦출 수 있었고, 데이터베이스에 의존한 문제가
발생하지 않았다. 중요한 것과 중요하지 않은 것 사이에 선을 긋는다. UI는 비즈니스 규칙에 영향을
미치지 않아야 하고, 데이터베이스는 UI에 영향을 미치지 않아야 한다. 물론, 대부분의 우리들은
데이터베이스는 비즈니스 규칙과 불가피하게 연결되어 있다고 믿고 있다. 하지만, 데이터베이스는
비즈니스 규칙이 간접적으로 사용할 수 있는 도구일 뿐이다. 그래서, 우리는 인터페이스 뒤에
데이터베이스를 놓을 수 있도록 설계를 해야 한다. 실제로 소프트웨어 개발, 기술의 역사는 확장
가능하고 유지 관리 가능한 시스템 아키텍처를 구축하기 위해 플러그인을 만드는 방법에 관한
이야기이다. 핵심 비즈니스 규칙은 다른 컴포넌트들과 독립적으로 유지된다.  - 로버트 C. 마틴

### layered Architecture의 적용

- 3tier Architecture : Clean architecture
  전자정부 프레임워크로 국내에서 가장 많이 사용하는 구조
    - presentation layer : controllers 폴더
    - business layer : services 폴더
    - data access layer : repositories 폴더
- Presentation layer
  HTTP Framework(Echo, Gin?)로부터 최초로 호출되는 API엔드 포인트
    - 클라이언트에서 보내온 데이터의 변환(Param Data)
    - 기본적인 인증과 요청 내용 검증
    - 수행 결과를 클라이언트에 반환
- Business layer
    - 비즈니스 로직을 작성한다.
- Data access layer
    - 데이터(데이터 정책(Unique, Max, Min))과 관련된 비즈니스 로직(메서드)
- Entities => model 폴더
    - 레이어 전체에서 사용되는 데이터의 구조
    - 메서드를 포함하는 객체일 수도 있고, 단순 데이터 구조일 수도 있다.
- Infrastructure ( Framework Drivers )
    - 이 영역은 일반적으로 데이터베이스 및 웹 프레임워크와 같은 도구로 구성된다.
    - HTTP Framework : Server
    - Database(ORM) : datastore
    - ETC...

참고용 그림( 이 프로젝트 코드와 관계는 없다. )
![Arch Img](./docs/img/layer_arch.png)

### 클린 아키텍처, 디펜던시 의존성

클린 아키텍처에서 가장 중요한 개념은 디펜던시(의존성) 규칙이다.
우리는 이 문서에서 설명한 계층 외 추가 계층을 필요로 할 수 있다. 그러나, 디펜던시 규칙은 항상
적용이 된다. 규칙을 준수하는 것은 어렵지 않으며, 앞으로 많은 고민들을 해결해 줄 것이다.
소프트웨어를 계층으로 분리하고, 디펜던시 규칙을 준수함으로써, 데이터베이스나 웹 프레임워크와 같은
시스템이 외부 부분들이 쓸모없게 될 때, 그러한 쓸모없는 요소들을 최소한의 작업으로 대체할 수 있을 것이다.

## 3. 개발

개발 시 다음 내용을 참고하여 작성한다.
문서상에서는 사용자 입력(요청)을 받아 데이터 베이스 제어하는 부분까지의 내용을 다룬다.

### 3.1. Change Module Name

코드 작성에 앞서 프로젝트(모듈)의 이름을 변경한다.
다음 명령은 소스 코드에서 `{search}`를 찾아 `{replace}`로 변경한다.

```sh
# foo를 찾아 bar로 변경한다.
$ find . -type f -name "*.go" -print0 | xargs -0 sed -i 's/foo/bar/g'

# The easier and much more readable option is to use another delimiter
# character. Most people use the vertical bar (|) or colon (:) but you
# can use any other character:

# github.com/mobigen/test => github.com/mobigen/blahblah 로 변경
$ find . -type f -name "*.go" -print0 | xargs -0 sed -i 's|github.com/mobigen/test|github.com/mobigen/blahblah|g'
```

### 3.2. Create Data(DTO, Entity) Structure

데이터베이스, Client와 주고 받을 데이터 정의

- models/sample.go

```go
package models

// Sample ....
type Sample struct {
    ID        int       `json:"id", gorm:"column:id;primaryKey;autoIncrement"`
    Name      string    `json:"name" gorm:"column:name;not null"`
    Desc      string    `json:"desc" gorm:"column:desc;size:256"`
    CreateAt  int64     `json:"createAt gorm:"column:createAt;autoCreateTime:milli"`
}
```

### 3.3. 테이블 생성(Migration)

- infrastructures/datastore/gorm.go - Migrate 함수

```go
func (ds *DataStore) Migrate() error {
    ...
    ds.Orm.AutoMigrate(&models.Sample{})
    ...
}
```

### 3.4. Usecase 작성

- controllers/sample.go

    ```go
    type SampleUsecase interface {
        GetAll()(*[]models.Sample, error)
        GetByID(int)(*models.Sample, error)
        Create(*models.Sample)(*models.Sample, error)
        Update(*models.Sample)(*models.Sample, error)
        Delete(int)(*models.Sample, error)
    }
    ```

- services/sample.go

    ```go
    type SampleRepository interface {
        GetAll()(*[]models.Sample, error)
        GetByID(int)(*models.Sample, error)
        Create(*models.Sample)(*models.Sample, error)
        Update(*models.Sample)(*models.Sample, error)
        Delete(int)(*models.Sample, error)
    }
    ```

### 3.5. 작성 순서

기능 정의, 모델 작성, Usecase까지 정의 되어야 하는 부분이 완료되었다면,
Controller, Service, Repositories 중 어떤 것을 먼저 작성하더라도 상관 없다.
각 레이어간 필요로 하는 인터페이스(모델)이 모두 정의되어 있으므로, 필요로 하는
부분을 처리하도록 구현을 진행하면 되기 때문이다.

### 3.6. Repositories 작성

- repositories/sample.go

    ```go
    package repositories

    import (
        "fmt"

        "github.com/mobigen/golang-web-template/infrastructures/datastore"
        "github.com/mobigen/golang-web-template/infrastructures/tools/util"
        "github.com/mobigen/golang-web-template/models"
    )

    // Sample is struct of todo.
    type Sample struct {
        *datastore.DataStore
    }

    // New is constructor that creates SampleRepository
    func (Sample) New(handler *datastore.DataStore) *Sample {
        return &Sample{handler}
    }

    // GetAll get all sample from database(store)
    func (repo *Sample) GetAll() (*[]models.Sample, error) {
        dst := new([]models.Sample)
        result := repo.Orm.Find(dst)
        if result.Error != nil {
            return nil, result.Error
        }
        if result.RowsAffected <= 0 {
            return nil, fmt.Errorf("no have result")
        }
        return dst, nil
    }

    // GetByID get sample whoes id match
    func (repo *Sample) GetByID(id int) (*models.Sample, error) {
        var dst *models.Sample
        result := repo.Orm.Find(dst).Where(&models.Sample{ID: id})
        if result.Error != nil {
            return nil, result.Error
        }
        if result.RowsAffected <= 0 {
            return nil, fmt.Errorf("no have result")
        }
        return dst, nil
    }

    // Create create sample
    func (repo *Sample) Create(input *models.Sample) (*models.Sample, error) {
        input.CreateAt = util.GetMillis()
        result := repo.Orm.Create(input)
        if result.Error != nil {
            return nil, result.Error
        }
        return input, nil
    }

    // Update update sample
    func (repo *Sample) Update(input *models.Sample) (*models.Sample, error) {
        // Save/Update All Fields
        // repo.Orm.Save(input)

        // ID       int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
        // Name     string `json:"name" gorm:"column:name;not null"`
        // Desc     string `json:"desc" gorm:"column:desc;size:256"`
        // CreateAt int64  `json:"createAt" gorm:"column:createAt;autoCreateTime:milli"`
        result := repo.Orm.Model(input).
            Where(&models.Sample{ID: input.ID}).
            Updates(
                map[string]interface{}{
                    "name": input.Name,
                    "desc": input.Desc,
                })
        if result.Error != nil {
            return nil, result.Error
        }
        return input, nil

    }

    // Delete delete sample from id(primaryKey)
    func (repo *Sample) Delete(id int) (*models.Sample, error) {
        dst := new(models.Sample)
        result := repo.Orm.Find(dst).Where(&models.Sample{ID: id})
        if result.Error != nil {
            return nil, result.Error
        }
        if result.RowsAffected <= 0 {
            return nil, fmt.Errorf("no have result")
        }
        // Delete with additional conditions
        result = repo.Orm.Delete(&models.Sample{}, id)
        if result.Error != nil {
            return nil, result.Error
        }
        return dst, nil
    }
    ```

### 3.7. Services 작성

- services/sample.go

    ```go
    package services

    import (
        "github.com/mobigen/golang-web-template/models"
    )

    // Sample service - repository - interactor for Sample entity.
    type Sample struct {
        Repo SampleRepository
    }

    // New is constructor that creates Sample service
    func (Sample) New(repo SampleRepository) *Sample {
        return &Sample{repo}
    }

    // GetAll returns All of samples.
    func (service *Sample) GetAll() ([]*models.Sample, error) {
        return service.Repo.GetAll()
    }

    // GetByID returns sample whoes that ID mathces.
    func (service *Sample) GetByID(id int) (*models.Sample, error) {
        return service.Repo.GetByID(id)
    }

    // Create create a new sample.
    func (service *Sample) Create(sample *models.Sample) (*models.Sample, error) {
        return service.Repo.Create(sample)
    }

    // Update update a sample.
    func (service *Sample) Update(sample *models.Sample) (*models.Sample, error) {
        return service.Repo.Update(sample)
    }

    // Delete delete sample from id.
    func (service *Sample) Delete(id int) (*models.Sample, error) {
        return service.Repo.Delete(id)
    }
    ```

### 3.8. Controller 작성

- controllers/sample.go

    ```go
    package controllers

    import (
        "net/http"
        "strconv"

        "github.com/mobigen/golang-web-template/models"
        "github.com/labstack/echo/v4"
    )

    // Sample Controller
    type Sample struct {
        Usecase SampleUsecase
    }

    // SampleUsecase usecase define
    type SampleUsecase interface {
        GetAll()(*[]models.Sample, error)
        GetByID(int)(*models.Sample, error)
        Create(*models.Sample)(*models.Sample, error)
        Update(*models.Sample)(*models.Sample, error)
        Delete(int)(*models.Sample, error)
    }

    // New create Sample instance.
    func (Sample) New(usecase SampleUsecase) *Sample {
        return &Sample{usecase}
    }

    // GetAll returns all of sample as JSON object.
    func (controller *Sample) GetAll(c echo.Context) error {
        samples, err := controller.Usecase.GetAll()
        if err != nil {
            return c.JSON(http.StatusBadRequest, samples)
        }
        return c.JSON(http.StatusOK, samples)
    }

    // GetByID return sample whoes ID mathces
    func (controller *Sample) GetByID(c echo.Context) error {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            return c.JSON(http.StatusBadRequest, err)
        }
        sample, err := controller.Usecase.GetByID(id)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }
        return c.JSON(http.StatusOK, sample)
    }

    // Create create a new ...
    func (controller *Sample) Create(c echo.Context) error {
        input := new(models.Sample)
        c.Bind(input)
        sample, err := controller.Usecase.Create(input)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }
        return c.JSON(http.StatusCreated, sample)
    }

    // Update update from input
    func (controller *Sample) Update(c echo.Context) error {
        input := new(models.Sample)
        c.Bind(input)
        sample, err := controller.Usecase.Update(input)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }
        return c.JSON(http.StatusOK, sample)
    }

    // Delete delete sample from id
    func (controller *Sample) Delete(c echo.Context) error {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            return c.JSON(http.StatusBadRequest, err)
        }
        sample, err := controller.Usecase.Delete(id)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }
        return c.JSON(http.StatusOK, sample)
    }
    ```

### 3.9. Injector 작성

- injectors/sample.go : 신규 추가 코드들에 대한 초기화

    ```go
    package injectors

    import (
        "github.com/mobigen/golang-web-template/controllers"
        "github.com/mobigen/golang-web-template/repositories"
        "github.com/mobigen/golang-web-template/services"
    )

    // Sample sample injector
    type Sample struct{}

    // Init for interconnection [ controller(App) - Service(Repository) - repository - datastore ] : Dependency Injection
    func (Sample) Init(in *Injector) *controllers.Sample {
        repo := repositories.Sample{}.New(in.Datastore)
        svc := services.Sample{}.New(repo)
        return controllers.Sample{}.New(svc)
    }
    ```

- PATH 등록 : injectors/injector-core.go - Init 함수에서 path와 controller.func 를 연결해 준다.

    ```go
    // Init ...
    func (h *Injector) Init() error {
        // path grouping
        apiv1 := h.Router.Group("/api/v1")

        // Sample
        h.Log.Errorf("[ PATH ] /api/v1/sample ........................................................... [ OK ]")
        sample := Sample{}.Init(h)
        apiv1.GET("/samples", sample.GetAll)
        apiv1.GET("/sample/:id", sample.GetByID)
        apiv1.POST("/sample", sample.Create)
        apiv1.POST("/sample/update", sample.Update)
        apiv1.DELETE("/sample/:id", sample.Delete)
        return nil
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

`TARGET`, `VERSION`, `BUILD_HASH`는 빌드 타임에 `common/appdata` 패키지 변수로 주입된다(`ldflags -X` 옵션).

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

**golang/mock 프레임워크** 를 이용한 mock 생성과 유닛 테스트에 대해서 설명한다.

Golang은 투명한 동작과 엄격한 타입이 특징인 만큼 Go에서는 Mocking을
이용하는 것은 다른 언어에 비해 쉬운 편은 아닌 것처럼 느껴지기도 합니다.

- interface로 선언된 변수에만 mock type을 할당할 수 있다.
  ( 그래서 우리는 레이어 구조를 만들면서 각 레이거간 의존성을 인터페이스로 만들었습니다. )
- mock type을 직접 정의하거나 mocking framework을 이용해 코드를 생성한다.
  ( framework를 사용하더라도 내부는 직접 작성하고 실행해야 한다. )

### 5.1. 다른 프레임워크와의 비교

GoMock vs. Testify: Mocking frameworks for Go
[gomock-vs-testify](https://blog.codecentric.de/2019/07/gomock-vs-testify/)
위 링크에 내용이 잘 정리되어있으니 참고해보세요.( 영문사이트... )

두 진영 비교
[Star](https://umi0410.github.io/blog/golang/how-to-backend-in-go-testcode/star-comparison.png)
우리는 어떤 진영을 선택할지는 중요하지 않습니다. 테스트를 만드는 것은 우리의 몫이기 때문입니다.

### 5.2. 사용 예시

이 프로젝트의 `sample` 코드를 기준으로 mock 생성부터 유닛 테스트 실행까지의 전체 흐름을 설명합니다.

#### 레이어 구조 및 인터페이스 위치

```
controllers/sample.go  → SampleUsecase 인터페이스 정의 (services 레이어를 추상화)
services/sample.go     → SampleRepository 인터페이스 정의 (repositories 레이어를 추상화)
repositories/sample.go → 실제 DB 접근 구현체
```

각 레이어는 아래 레이어를 인터페이스로만 의존하기 때문에, mock 으로 교체하여 독립적인 유닛 테스트가 가능합니다.

---

#### Step 1. `//go:generate` 지시어 추가

mock 을 생성할 인터페이스가 있는 파일에 `//go:generate` 주석을 추가합니다.

**`services/sample.go`** — `SampleRepository` 인터페이스에 추가

```go
package services

//go:generate mockgen -destination=../mocks/mock_sample_repository.go -package=mocks github.com/mobigen/golang-web-template/services SampleRepository

type SampleRepository interface {
    GetAll() (*[]models.Sample, error)
    GetByID(int) (*models.Sample, error)
    Create(*models.Sample) (*models.Sample, error)
    Update(*models.Sample) (*models.Sample, error)
    Delete(int) (*models.Sample, error)
}
```

**`controllers/sample.go`** — `SampleUsecase` 인터페이스에 추가

```go
package controllers

//go:generate mockgen -destination=../mocks/mock_sample_usecase.go -package=mocks github.com/mobigen/golang-web-template/controllers SampleUsecase

type SampleUsecase interface {
    GetAll() (*[]models.Sample, error)
    // ...
}
```

`mockgen` 옵션 설명:

| 옵션                          | 설명                                                      |
| ----------------------------- | --------------------------------------------------------- |
| `-destination`                | 생성될 mock 파일 경로                                     |
| `-package`                    | 생성될 mock 의 패키지 이름                                |
| 마지막 인자 (패키지 경로)     | mock 을 생성할 인터페이스가 있는 패키지                   |
| 마지막 인자 (인터페이스 이름) | mock 을 생성할 인터페이스 이름 (콤마로 여러 개 지정 가능) |

---

#### Step 2. mock 생성 (`make mocks`)

```sh
make mocks
```

이 명령은 내부적으로 `go generate ./...` 를 실행하여 `//go:generate` 지시어가 있는 모든 파일에서 mock 을 생성합니다.

실행 결과로 `mocks/` 디렉토리에 파일이 생성됩니다:

```txt
mocks/
  mock_sample_repository.go   ← services.SampleRepository mock
  mock_sample_usecase.go      ← controllers.SampleUsecase mock
```

생성된 mock 의 구조 (예: `mock_sample_repository.go`):

```go
// Code generated by MockGen. DO NOT EDIT.
package mocks

type MockSampleRepository struct {
    ctrl     *gomock.Controller
    recorder *MockSampleRepositoryMockRecorder
}

// NewMockSampleRepository creates a new mock instance.
func NewMockSampleRepository(ctrl *gomock.Controller) *MockSampleRepository { ... }

// EXPECT 를 통해 호출 기대값을 설정한다.
func (m *MockSampleRepository) EXPECT() *MockSampleRepositoryMockRecorder { ... }

func (m *MockSampleRepository) GetAll() (*[]models.Sample, error) { ... }
// ... 나머지 메서드
```

---

#### Step 3. 유닛 테스트 작성

**Service 레이어 테스트** (`services/sample_test.go`)

Service 는 `SampleRepository` 인터페이스에 의존하므로, `MockSampleRepository` 로 교체하여 테스트합니다.

```go
package services_test

import (
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

    // 1. mock 생성
    mockRepo := mocks.NewMockSampleRepository(ctrl)
    // 2. mock 을 주입하여 service 생성
    svc := services.Sample{}.New(mockRepo)

    expected := &[]models.Sample{
        {ID: 1, Name: "foo", Desc: "bar"},
    }
    // 3. mock 기대값 설정: GetAll() 이 1회 호출되면 expected 를 반환
    mockRepo.EXPECT().GetAll().Return(expected, nil).Times(1)

    // 4. 테스트 실행
    result, err := svc.GetAll()

    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

**Controller 레이어 테스트** (`controllers/sample_test.go`)

Controller 는 `SampleUsecase` 인터페이스에 의존하므로, `MockSampleUsecase` 로 교체합니다.
HTTP 요청/응답은 `net/http/httptest` 와 echo 를 사용합니다.

```go
package controllers_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    gomock "github.com/golang/mock/gomock"
    "github.com/labstack/echo/v5"
    "github.com/mobigen/golang-web-template/controllers"
    "github.com/mobigen/golang-web-template/mocks"
    "github.com/mobigen/golang-web-template/models"
    "github.com/stretchr/testify/assert"
)

func TestSampleController_GetByID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // 1. mock 생성 및 controller 주입
    mockUsecase := mocks.NewMockSampleUsecase(ctrl)
    controller := controllers.Sample{}.New(mockUsecase)

    expected := &models.Sample{ID: 1, Name: "foo", Desc: "bar"}
    mockUsecase.EXPECT().GetByID(1).Return(expected, nil).Times(1)

    // 2. echo Context 구성 (echo v5 에서 path parameter 설정 방법)
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/samples/1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetPathValues(echo.PathValues{{Name: "id", Value: "1"}})

    // 3. 테스트 실행 및 검증
    err := controller.GetByID(c)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)

    var result models.Sample
    json.Unmarshal(rec.Body.Bytes(), &result)
    assert.Equal(t, 1, result.ID)
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

---

#### Step 4. 테스트 실행 (`make test`)

```sh
make test
```

출력 예시:

```sh
make test
=== RUN   TestSampleService_GetAll
--- PASS: TestSampleService_GetAll (0.00s)
=== RUN   TestSampleService_GetByID
--- PASS: TestSampleService_GetByID (0.00s)
=== RUN   TestSampleService_Create
--- PASS: TestSampleService_Create (0.00s)
=== RUN   TestSampleController_GetAll
--- PASS: TestSampleController_GetAll (0.00s)
=== RUN   TestSampleController_GetByID
--- PASS: TestSampleController_GetByID (0.00s)
--- PASS: TestSampleController_GetByID_InvalidID (0.00s)
PASS
ok  github.com/mobigen/golang-web-template/services      0.52s
ok  github.com/mobigen/golang-web-template/controllers   0.89s
```

커버리지 리포트 생성:

```sh
make coverage
# build/cov-out.html 파일로 브라우저에서 확인 가능
```

## 6. 문서작성(Swagger)

API 문서를 코드 개발과 분리하여 처리하는 것이 아닌, 코드 개발 단계에서 `// @` 형식의 주석으로 작성하고
`make swag`(`swag init` 명령을 makefile 에 추가) 명령으로 자동 생성하는 방식입니다. **코드 주석이 곧 API 문서**가 됩니다.

자세한 내용은 [swaggo/swag](https://github.com/swaggo/swag) 를 참고하세요.

### 6.1. 서버 등록

`infrastructures/router/server.go` 에서 Swagger 경로가 이미 등록되어 있습니다.

```go
import (
    _ "github.com/mobigen/golang-web-template/docs/swagger"  // swag init 으로 생성된 docs 패키지
    echoSwagger "github.com/swaggo/echo-swagger"
)

// Swagger
r.GET("/swagger/*", echoSwagger.WrapHandler)
```

### 6.2. 작성 방법

**① 서버 전체 정보** — `main.go` 의 `main()` 함수 위에 작성

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

**② 개별 API** — 각 handler 함수 위에 작성 (`controllers/version.go` 예시)

```go
// GetVersion return app version
// @Summary Get Server Version
// @Description get server version info
// @Tags version
// @Accept  json
// @Produce  json
// @success 200 {object} controllers.HTTPResponse{data=appdata.VersionInfo} "app info(name, version, hash)"
// @Router /version [get]
func (controller *Version) GetVersion(c *echo.Context) error {
    ...
}
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

### 6.3. 문서 생성

```sh
make swag
```

외부 패키지 타입을 참조하는 경우 파싱 오류가 발생할 수 있습니다.
다음과 같이 `Makefile`을 수정하여 `swag init`이 동작할 수 있도록 합니다.

```sh
# ParseComment error: cannot find type definition: ...
swag init --parseDependency --parseInternal --output docs/swagger
```

생성 결과로 `docs/` 디렉토리에 다음 파일들이 생성됩니다.

```txt
docs/swagger/
  docs.go        ← Go 패키지 (서버 import 용)
  swagger.json   ← Swagger 명세
  swagger.yaml   ← Swagger 명세 (YAML 형식)
```

### 6.4. 문서 확인

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
  output: "stdout"       # stdout, file, both
  level: "debug"         # debug, info, warn, error, silent
  # 파일 출력 시 추가 옵션 (output: "file" 일 때 유효)
  savePath: "logs"  # $APP_HOME 기준 상대경로 또는 절대경로
  fileName: "app.log"  # 로그 파일의 이름, 백업 시 app-{2026-01-01T14:15:12.000}.log
  sizePerFileMb: 100        # 로그 파일 최대 크기 (MB)
  maxOfDay: 10              # 보관할 백업 파일 수
  maxAge: 7                 # 로그 파일 보관 기간 (일)
  compress: false           # 오래된 로그 파일 gzip 압축 여부
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
go run main.go

# 환경 변수 지정
APP_HOME=$(pwd) PROFILE=prod LOG_LEVEL=debug go run main.go
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

**종료**

`Ctrl+C` 또는 `SIGTERM` 시그널을 보내면 graceful shutdown 됩니다.

```log
[ SIGNAL ] Receive [ terminated ]
[ DataStore ] Shutdown ............................................................ [ OK ]
[ Router ] Shutdown ............................................................... [ OK ]
```
