# 코드 정리

## 개요

현재 디렉토리 구성은 프로젝트 루트에 난잡하게 흩어져 있어 정리하고자 하.ㅁ

## 클린 아키텍처 스타일 Go 프로젝트 예시

엄격하게 의존성 방향을 관리하는 방식입니다.

핵심 원칙은 이것입니다.

외부 -> 내부로 의존 가능
내부 -> 외부를 알면 안 됨

즉:

- domain은 Echo를 몰라야 함
- usecase는 DB 구현을 몰라야 함
- interface adapter가 변환 역할을 맡음

---

추천 구조

```txt
myapp/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   └── errors.go
│   ├── usecase/
│   │   └── user/
│   │       ├── create.go
│   │       ├── get.go
│   │       └── interface.go
│   ├── adapter/
│   │   ├── inbound/
│   │   │   └── http/
│   │   │       ├── handler/
│   │   │       │   └── user_handler.go
│   │   │       ├── dto/
│   │   │       │   ├── request.go
│   │   │       │   └── response.go
│   │   │       └── router.go
│   │   └── outbound/
│   │       └── persistence/
│   │           └── mysql/
│   │               └── user_repository.go
│   ├── infrastructure/
│   │   ├── db/
│   │   │   └── mysql.go
│   │   └── logger/
│   │       └── logger.go
│   └── bootstrap/
│       └── wire.go
├── go.mod
└── README.md
```

의존성 방향

handler -> usecase -> repository interface
repository implementation -> repository interface 구현

즉 실제 DB 구현체는 바깥쪽에 있고,
usecase는 인터페이스만 의존합니다.

예시

`internal/domain/user.go`

```go
package domain

type User struct {
	ID    uint64
	Name  string
	Email string
}
```

`internal/usecase/user/interface.go`

```go
package user

import (
	"context"

	"myapp/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uint64) (*domain.User, error)
}
```

`internal/usecase/user/create.go`

```go
package user

import (
	"context"

	"myapp/internal/domain"
)

type CreateUseCase struct {
	repo Repository
}

func NewCreateUseCase(repo Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

func (uc *CreateUseCase) Execute(ctx context.Context, name, email string) (*domain.User, error) {
	user := &domain.User{
		Name:  name,
		Email: email,
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
```

`internal/usecase/user/get.go`

```go
package user

import (
	"context"

	"myapp/internal/domain"
)

type GetUseCase struct {
	repo Repository
}

func NewGetUseCase(repo Repository) *GetUseCase {
	return &GetUseCase{repo: repo}
}

func (uc *GetUseCase) Execute(ctx context.Context, id uint64) (*domain.User, error) {
	return uc.repo.FindByID(ctx, id)
}
```

`internal/adapter/inbound/http/handler/user_handler.go`

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	useruc "myapp/internal/usecase/user"
)

type UserHandler struct {
	createUC *useruc.CreateUseCase
	getUC    *useruc.GetUseCase
}

func NewUserHandler(createUC *useruc.CreateUseCase, getUC *useruc.GetUseCase) *UserHandler {
	return &UserHandler{
		createUC: createUC,
		getUC:    getUC,
	}
}

func (h *UserHandler) Create(c echo.Context) error {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	user, err := h.createUC.Execute(c.Request().Context(), req.Name, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid id"})
	}

	user, err := h.getUC.Execute(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "not found"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
```

`internal/adapter/outbound/persistence/mysql/user_repository.go`

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
		"INSERT INTO users(name, email) VALUES(?, ?)",
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

`internal/bootstrap/wire.go`

```go
package bootstrap

import (
	"database/sql"

	"myapp/internal/adapter/inbound/http/handler"
	"myapp/internal/adapter/outbound/persistence/mysql"
	useruc "myapp/internal/usecase/user"
)

type App struct {
	UserHandler *handler.UserHandler
}

func NewApp(db *sql.DB) *App {
	userRepo := mysql.NewUserRepository(db)

	createUC := useruc.NewCreateUseCase(userRepo)
	getUC := useruc.NewGetUseCase(userRepo)

	userHandler := handler.NewUserHandler(createUC, getUC)

	return &App{
		UserHandler: userHandler,
	}
}
```

## 3. 테스트를 포함한 구성

보통은 테스트 대상 코드 가까이에 둔다가 기본 원칙입니다.
Go에서는 이 방식이 가장 자연스럽습니다.

핵심부터 말하면:

- unit test: 해당 패키지 안에 같이 둠
- mock: 그 mock을 사용하는 테스트와 최대한 가까이 둠
- integration test: 별도 디렉토리나 패키지로 분리 가능
- e2e test: 루트 하위 test/ 또는 tests/로 분리

### 클린 아키텍처 구조에서의 test/mock 배치

클린 아키텍처에서는 인터페이스 경계가 분명해서 mock 배치도 조금 더 체계적으로 갈 수 있습니다.

예를 들어 구조가 이렇다면:

internal/
├── domain/
├── usecase/
├── adapter/
│   ├── inbound/http/
│   └── outbound/persistence/mysql/
├── infrastructure/
└── bootstrap/

테스트 배치는 보통 이렇게 갑니다.

myapp/
├── internal/
│   ├── domain/
│   │   └── user.go
│   ├── usecase/
│   │   └── user/
│   │       ├── create.go
│   │       ├── get.go
│   │       ├── interface.go
│   │       ├── create_test.go
│   │       ├── get_test.go
│   │       └── mock_repository_test.go
│   ├── adapter/
│   │   ├── inbound/
│   │   │   └── http/
│   │   │       ├── handler/
│   │   │       │   ├── user_handler.go
│   │   │       │   ├── user_handler_test.go
│   │   │       │   └── mock_usecase_test.go
│   │   │       └── dto/
│   │   └── outbound/
│   │       └── persistence/
│   │           └── mysql/
│   │               ├── user_repository.go
│   │               ├── user_repository_test.go
│   │               └── testdata/
│   │                   └── users.sql
│   └── bootstrap/
│       └── wire.go
├── test/
│   ├── integration/
│   │   ├── create_user_flow_test.go
│   │   └── fixture/
│   └── e2e/
│       └── http_e2e_test.go
└── go.mod

---

클린 아키텍처에서의 포인트

usecase 테스트

usecase는 repository interface에 의존하므로,
mock repository를 두고 단위 테스트하기 아주 좋습니다.

internal/usecase/user/
├── create.go
├── create_test.go
└── mock_repository_test.go

이 배치가 가장 흔합니다.

inbound adapter 테스트

HTTP handler는 usecase를 mock으로 대체해서 테스트합니다.

internal/adapter/inbound/http/handler/
├── user_handler.go
├── user_handler_test.go
└── mock_usecase_test.go

outbound adapter 테스트

DB 구현체는 mock보다 실제 DB와 붙여보는 테스트 비중이 큽니다.

internal/adapter/outbound/persistence/mysql/
├── user_repository.go
└── user_repository_test.go

---

### 3.3. mock을 어디까지 공용화할까?

이 부분이 가장 많이 헷갈립니다.

결론은:

- 작은 프로젝트: mock을 테스트 파일 근처에 둠
- 큰 프로젝트: 반복되는 mock만 별도 패키지로 분리

---

작은 프로젝트 권장

가장 추천하는 방식입니다.

internal/service/
├── user_service.go
├── user_service_test.go
└── mock_user_repository_test.go

장점:

- 테스트 읽기 쉬움
- mock 수정 영향 범위 작음
- 쓸데없는 공용화 방지

---

큰 프로젝트에서만 shared mock 고려

예:

internal/testutil/
├── mock/
│   ├── user_repository.go
│   └── clock.go
└── fixture/
    └── user.go

하지만 이건 정말 중복이 심할 때만 추천합니다.

이유:
	•	공용 mock이 점점 비대해짐
	•	사용하지 않는 메서드까지 다 구현됨
	•	테스트 의도가 흐려짐

즉, testutil은 최소화하는 게 좋습니다.

---

### 3.4. testify/mock, mockgen 같은 자동 생성 mock은 어디에 둘까?

패턴 B: mocks 디렉토리 분리

예:

internal/usecase/user/
├── interface.go
└── mocks/
    └── repository.gen.go

혹은

internal/mocks/
├── user_repository.gen.go
└── create_user_usecase.gen.go

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

### 3.5. testdata는 어떻게 두나?

Go는 testdata/ 디렉토리를 특별 취급하므로 적극 추천합니다.

예:

internal/handler/http/
├── user_handler_test.go
└── testdata/
    ├── create_user_request.json
    └── create_user_response.json

또는 repository 테스트용:

internal/adapter/outbound/persistence/mysql/
├── user_repository_test.go
└── testdata/
    ├── schema.sql
    ├── seed.sql
    └── expected_users.json

testdata/에는 보통 이런 걸 둡니다.
	•	JSON 요청/응답 샘플
	•	SQL schema/seed
	•	fixture 파일
	•	golden file

### 3.6. 내가 추천하는 실제 배치

클린 아키텍처용 추천

internal/
├── usecase/
│   └── user/
│       ├── create.go
│       ├── create_test.go
│       ├── get.go
│       ├── get_test.go
│       └── mock_repository_test.go
├── adapter/
│   ├── inbound/http/handler/
│   │   ├── user_handler.go
│   │   ├── user_handler_test.go
│   │   └── mock_usecase_test.go
│   └── outbound/persistence/mysql/
│       ├── user_repository.go
│       ├── user_repository_test.go
│       └── testdata/
└── domain/

### 3.7. 피하는 게 좋은 배치

이런 건 보통 나중에 불편해집니다.

internal/
└── mocks/
    ├── everything.go
    ├── all_repository_mocks.go
    └── all_service_mocks.go

이유:
	•	mock이 거대해짐
	•	테스트 의도 파악이 어려움
	•	변경 영향 범위가 넓음

그리고 이것도 자주 안 좋습니다.

pkg/testutil/

프로젝트 내부 테스트 전용이면 internal/testutil/ 쪽이 더 자연스럽습니다.
