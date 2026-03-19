# 코드 정리

## 개요

현재 디렉토리 구성은 프로젝트 루트에 난잡하게 흩어져 있어 정리하고자 하.ㅁ

## Echo 기반 REST API용 추천 구조

```txt
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   └── user.go
│   ├── dto/
│   │   ├── user_request.go
│   │   └── user_response.go
│   ├── handler/
│   │   └── http/
│   │       ├── router.go
│   │       ├── user_handler.go
│   │       └── health_handler.go
│   ├── service/
│   │   └── user_service.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── mysql/
│   │       └── user_repository.go
│   ├── platform/
│   │   ├── db/
│   │   │   └── mysql.go
│   │   └── logger/
│   │       └── logger.go
│   └── middleware/
│       ├── recover.go
│       ├── logger.go
│       └── auth.go
├── migrations/
├── scripts/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

각 계층 역할

`cmd/server/main.go`

앱 시작점입니다.

- 설정 로드
- DB 연결
- Echo 생성
- 미들웨어 등록
- 라우터 등록
- 서버 실행

`internal/config`

환경변수, yaml, json 등 설정을 로드합니다.

`internal/domain`

핵심 모델을 둡니다.

예:

```go
package domain

import "time"

type User struct {
	ID        uint64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

`internal/dto`

HTTP 요청/응답 전용 구조체입니다.
domain과 분리하는 걸 추천합니다.

예:

```go
package dto

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
```

`internal/handler/http`

Echo 핸들러입니다.

역할:

- 요청 바인딩
- validation
- service 호출
- HTTP 응답 생성

`internal/service`

비즈니스 로직 담당입니다.

역할:

- 유스케이스 수행
- 트랜잭션 경계 제어
- 도메인 규칙 적용

`internal/repository`

DB 접근 인터페이스와 구현체입니다.

- 상위 레벨에 interface
- 하위 디렉토리에 mysql/postgres 구현체

`internal/platform`

DB, logger, redis, 외부 SDK 같은 인프라 코드입니다.

`internal/middleware`

Echo 미들웨어 정의입니다.

---

Echo REST API 흐름

보통 이런 식입니다.

HTTP Request
  -> middleware
  -> handler
  -> service
  -> repository
  -> DB
  -> response

---

예시 코드

`cmd/server/main.go`

```go
package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"myapp/internal/config"
	httpHandler "myapp/internal/handler/http"
	"myapp/internal/platform/db"
	"myapp/internal/repository/mysql"
	"myapp/internal/service"
)

func main() {
	cfg := config.Load()

	database, err := db.NewMySQL(cfg)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	userRepo := mysql.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := httpHandler.NewUserHandler(userService)

	httpHandler.RegisterRoutes(e, userHandler)

	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
```

`internal/handler/http/router.go`

```go
package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, userHandler *UserHandler) {
	api := e.Group("/api/v1")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	api.POST("/users", userHandler.Create)
	api.GET("/users/:id", userHandler.GetByID)
}
```

`internal/handler/http/user_handler.go`

```go
package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"myapp/internal/dto"
	"myapp/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid request body",
		})
	}

	user, err := h.userService.Create(c.Request().Context(), req.Name, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid user id",
		})
	}

	user, err := h.userService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "user not found",
		})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
```

`internal/service/user_service.go`

```go
package service

import (
	"context"

	"myapp/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uint64) (*domain.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Create(ctx context.Context, name, email string) (*domain.User, error) {
	user := &domain.User{
		Name:  name,
		Email: email,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}
```

`internal/repository/user_repository.go`

```go
package repository
```

이 파일은 비워두기보다, 보통 인터페이스를 service 쪽에 둘지 repository 쪽에 둘지 기준을 정해서 갑니다.

실무에서는 대체로:

- service가 필요한 인터페이스를 정의
- repository는 구현체만 제공

이 방식이 더 깔끔합니다.

`internal/repository/mysql/user_repository.go`

```go
package mysql

import (
	"context"
	"database/sql"

	"myapp/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	result, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (name, email) VALUES (?, ?)",
		user.Name,
		user.Email,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint64(id)
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, email FROM users WHERE id = ?",
		id,
	)

	user := &domain.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, err
	}

	return user, nil
}
```

이 구조의 장점

- Echo 사용이 단순함
- handler/service/repository 역할이 명확함
- 테스트 작성이 쉬움
- 초반 생산성이 좋음
- 팀원이 봐도 바로 이해 가능

이 구조의 단점

- 프로젝트가 아주 커지면 기능별 응집도가 떨어질 수 있음
- service가 비대해질 수 있음
- domain이 단순 DTO처럼 변질될 가능성이 있음

## 2. 테스트를 포함한 구성

보통은 테스트 대상 코드 가까이에 둔다가 기본 원칙입니다.
Go에서는 이 방식이 가장 자연스럽습니다.

핵심부터 말하면:

- unit test: 해당 패키지 안에 같이 둠
- mock: 그 mock을 사용하는 테스트와 최대한 가까이 둠
- integration test: 별도 디렉토리나 패키지로 분리 가능
- e2e test: 루트 하위 test/ 또는 tests/로 분리

### 2.1. Echo 기반 REST API 구조에서의 test/mock 배치

기존 구조가 이런 형태라고 하면:

```txt
internal/
├── config/
├── domain/
├── dto/
├── handler/http/
├── service/
├── repository/
├── middleware/
└── platform/
```

테스트는 보통 이렇게 둡니다.

```txt
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   │   └── http/
│   │       ├── router.go
│   │       ├── user_handler.go
│   │       ├── user_handler_test.go
│   │       └── testdata/
│   │           └── user_request.json
│   ├── service/
│   │   ├── user_service.go
│   │   ├── user_service_test.go
│   │   └── mock_user_repository_test.go
│   ├── repository/
│   │   ├── mysql/
│   │   │   ├── user_repository.go
│   │   │   └── user_repository_test.go
│   │   └── repository_test.go
│   ├── domain/
│   │   └── user.go
│   └── platform/
│       └── db/
│           └── mysql.go
├── test/
│   ├── integration/
│   │   ├── user_api_test.go
│   │   └── docker-compose.yml
│   └── e2e/
│       └── api_e2e_test.go
└── go.mod
```

권장 원칙

1. unit test는 같은 디렉토리

    예:

    - user_service.go
    - user_service_test.go

    이렇게 같은 패키지에 둡니다.

    이유:

    - 파일 찾기 쉬움
    - 패키지 구조와 함께 읽기 좋음
    - Go 기본 스타일과 잘 맞음

2. mock은 “그 테스트가 있는 곳” 가까이에

    예를 들어 service 테스트에서 repository mock이 필요하면:

    ```txt
    internal/service/
    ├── user_service.go
    ├── user_service_test.go
    └── mock_user_repository_test.go
    ```

    이런 식이 가장 실용적입니다.

    파일명에 _test.go를 붙이면:

    - 테스트 빌드에서만 포함됨
    - 실제 프로덕션 바이너리에 섞이지 않음

    이 방식이 특히 좋습니다.

    예:

    ```go
    package service

    import (
        "context"

        "myapp/internal/domain"
    )

    type mockUserRepository struct {
        createFn   func(ctx context.Context, user *domain.User) error
        findByIDFn func(ctx context.Context, id uint64) (*domain.User, error)
    }

    func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
        return m.createFn(ctx, user)
    }

    func (m *mockUserRepository) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
        return m.findByIDFn(ctx, id)
    }
    ```

3. handler test용 mock도 handler 가까이에

    handler/http 테스트에서 service mock이 필요하면:

    ```txt
    internal/handler/http/
    ├── user_handler.go
    ├── user_handler_test.go
    └── mock_user_service_test.go
    ```

    이렇게 둡니다.

    즉, mock을 공용 디렉토리 하나에 몰아넣기보다, 사용하는 계층 근처에 둔다가 좋습니다.

4. repository test는 mock보다 실제 DB 테스트가 많음

    repository는 보통 인터페이스 mock보다:

    - sqlmock
    - testcontainers
    - docker-compose 기반 테스트 DB

    를 더 자주 씁니다.

    그래서 예:

    ```txt
    internal/repository/mysql/
    ├── user_repository.go
    ├── user_repository_test.go
    └── testdata/
        └── schema.sql
    ```

    혹은 통합 테스트로 빼서:

    ```txt
    test/integration/
    └── user_repository_integration_test.go
    ```

### 2.2. mock을 어디까지 공용화할까?

이 부분이 가장 많이 헷갈립니다.

결론은:

- 작은 프로젝트: mock을 테스트 파일 근처에 둠
- 큰 프로젝트: 반복되는 mock만 별도 패키지로 분리

---

작은 프로젝트 권장

가장 추천하는 방식입니다.

```txt
internal/service/
├── user_service.go
├── user_service_test.go
└── mock_user_repository_test.go
```

장점:

- 테스트 읽기 쉬움
- mock 수정 영향 범위 작음
- 쓸데없는 공용화 방지

큰 프로젝트에서만 shared mock 고려

예:

```txt
internal/testutil/
├── mock/
│   ├── user_repository.go
│   └── clock.go
└── fixture/
    └── user.go
```

하지만 이건 정말 중복이 심할 때만 추천합니다.

이유:

- 공용 mock이 점점 비대해짐
- 사용하지 않는 메서드까지 다 구현됨
- 테스트 의도가 흐려짐

즉, testutil은 최소화하는 게 좋습니다.

### 2.3. testify/mock, mockgen 같은 자동 생성 mock은 어디에 둘까?

mocks 디렉토리 분리

예:

```txt
internal/usecase/user/
├── interface.go
└── mocks/
    └── repository.gen.go
```

혹은

```txt
internal/mocks/
├── user_repository.gen.go
└── create_user_usecase.gen.go
```

테스트 전용이면 _test.go를 붙이는 편이 더 안전합니다.
mock 파일은 생성 파일이라 이름을 명확히 하는 게 좋습니다.

---

실무 추천

자동 생성 mock을 쓰더라도 가능하면 범위를 좁히세요.

예:

internal/usecase/user/
├── interface.go
├── create.go
├── create_test.go
└── mocks/
    └── repository.gen.go

이게 전역 mocks/보다 낫습니다.

### 2.4. testdata는 어떻게 두나?

Go는 testdata/ 디렉토리를 특별 취급하므로 적극 추천합니다.

예:

```txt
internal/handler/http/
├── user_handler_test.go
└── testdata/
    ├── create_user_request.json
    └── create_user_response.json
```

또는 repository 테스트용:

```txt
internal/adapter/outbound/persistence/mysql/
├── user_repository_test.go
└── testdata/
    ├── schema.sql
    ├── seed.sql
    └── expected_users.json
```

testdata/에는 보통 이런 걸 둡니다.

- JSON 요청/응답 샘플
- SQL schema/seed
- fixture 파일
- golden file

---

추천하는 실제 배치

Echo 계층형 구조용 추천

internal/
├── handler/
│   └── http/
│       ├── user_handler.go
│       ├── user_handler_test.go
│       ├── mock_user_service_test.go
│       └── testdata/
├── service/
│   ├── user_service.go
│   ├── user_service_test.go
│   └── mock_user_repository_test.go
├── repository/
│   └── mysql/
│       ├── user_repository.go
│       ├── user_repository_test.go
│       └── testdata/
└── domain/

### 2.5. 피하는 게 좋은 배치

이런 건 보통 나중에 불편해집니다.

internal/
└── mocks/
    ├── everything.go
    ├── all_repository_mocks.go
    └── all_service_mocks.go

이유:

- mock이 거대해짐
- 테스트 의도 파악이 어려움
- 변경 영향 범위가 넓음

그리고 이것도 자주 안 좋습니다.

`pkg/testutil/`

프로젝트 내부 테스트 전용이면 `internal/testutil/` 쪽이 더 자연스럽습니다.
