# uber-go/fx 도입 가이드

## 개요

`uber-go/fx`는 Uber에서 만든 Go용 의존성 주입(DI) 프레임워크입니다.
Spring의 `@Component` / `@Autowired` / `@PostConstruct` 개념을 Go에서 런타임 리플렉션으로 구현합니다.

현재 프로젝트의 `main.go Context` + `bootstrap/wire.go Injector` 구조를 fx로 대체하면
순서 관리, 라이프사이클 제어, 의존성 연결을 자동화할 수 있습니다.

---

## Spring vs fx 개념 대응표

| Spring | uber-go/fx | 설명 |
|--------|-----------|------|
| `@Component` / `@Service` | `fx.Provide()` | 의존성 생성자 등록 |
| `@Autowired` | 타입 기반 자동 주입 | 함수 파라미터 타입으로 자동 매핑 |
| `@PostConstruct` | `fx.Lifecycle.Append(OnStart)` | 초기화 훅 |
| `@PreDestroy` | `fx.Lifecycle.Append(OnStop)` | 종료 훅 |
| `@Configuration` | `fx.Module()` | 모듈 단위 묶음 |
| `ApplicationContext.run()` | `fx.New(...).Run()` | 앱 실행 |
| `@Bean` | `fx.Provide()` 내 팩토리 함수 | 인스턴스 제공 함수 |

---

## 설치

```bash
go get go.uber.org/fx
```

---

## 핵심 API

### fx.Provide

의존성 생성자(Constructor)를 컨테이너에 등록합니다.
함수의 **반환 타입**을 보고 자동으로 어떤 의존성인지 판단합니다.

```go
fx.Provide(
    NewLogger,      // *Logger 를 제공
    NewDatabase,    // *Database 를 제공
    NewRepository,  // *Repository 를 제공
)
```

생성자 함수 시그니처:

```go
// 파라미터 = 필요한 의존성 (자동 주입)
// 반환값   = 이 함수가 제공하는 의존성
func NewRepository(db *Database, log *Logger) *Repository {
    return &Repository{db: db, log: log}
}
```

### fx.Invoke

의존성을 소비하는 함수를 등록합니다.
`Provide`와 달리 반환값을 컨테이너에 저장하지 않고 **실행 자체가 목적**입니다.
라우트 등록, 서버 시작 등에 사용합니다.

```go
fx.Invoke(
    RegisterRoutes,  // 라우트 등록 (반환값 없음)
)
```

### fx.Lifecycle

컴포넌트의 시작/종료 훅을 등록합니다.
`OnStart`는 `app.Run()` 시점에, `OnStop`은 시그널 수신 또는 `app.Stop()` 시점에 호출됩니다.

```go
func NewDatabase(lc fx.Lifecycle, conf *Config) (*Database, error) {
    db := &Database{}
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            return db.Connect(conf.DSN)
        },
        OnStop: func(ctx context.Context) error {
            return db.Close()
        },
    })
    return db, nil
}
```

### fx.Module

관련 Provider를 모듈로 묶어 관리합니다.
Spring의 `@Configuration` 클래스와 유사합니다.

```go
var InfraModule = fx.Module("infra",
    fx.Provide(
        NewLogger,
        NewDatabase,
    ),
)

var DomainModule = fx.Module("domain",
    fx.Provide(
        NewSampleRepository,
        NewSampleHandler,
    ),
)
```

### fx.Annotate / fx.Tag

같은 타입의 의존성이 여러 개일 때 이름 태그로 구분합니다.
Spring의 `@Qualifier`와 유사합니다.

```go
// 제공 시 이름 지정
fx.Provide(
    fx.Annotate(NewMySQLDB,    fx.ResultTags(`name:"mysql"`)),
    fx.Annotate(NewPostgresDB, fx.ResultTags(`name:"postgres"`)),
)

// 주입 시 이름 지정
func NewRepository(
    db *gorm.DB `name:"mysql"`,
) *Repository { ... }
```

---

## 현재 프로젝트에 fx 적용 예시

### 현재 구조 vs fx 적용 후

**현재 main.go:**

```go
// 순서를 직접 관리, 에러 처리 수동
c := new(Context)
c.InitLog()
c.ReadEnv()
c.ReadConfig()
c.SetLogger()
c.InitDatastore()
c.InitRouter()
c.InitDepencyInjection()
c.StartSubModules()
```

**fx 적용 후 main.go:**

```go
func main() {
    fx.New(
        InfraModule,
        DomainModule,
        fx.Invoke(bootstrap.RegisterRoutes),
    ).Run()
}
```

---

### 변환 예시: Infrastructure 레이어

#### Logger (`internal/infrastructure/logger/logger.go`)

```go
// 기존
func (LogrusLogger) GetInstance() *LogrusLogger { ... }

// fx 적용 후
func NewFxLogger(conf *config.Configuration) (*LogrusLogger, error) {
    log := LogrusLogger{}.GetInstance()
    if err := log.Setting(&conf.Log, conf.Home); err != nil {
        return nil, err
    }
    log.Start()
    return log, nil
}
```

#### DataStore (`internal/infrastructure/db/gorm.go`)

```go
// fx 적용 후 - Lifecycle 훅으로 시작/종료 자동화
func NewFxDataStore(lc fx.Lifecycle, conf *config.Configuration, log *LogrusLogger) (*DataStore, error) {
    ds, err := DataStore{}.New(conf.Home, log.Logger)
    if err != nil {
        return nil, err
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            if err := ds.Connect(&conf.Datastore); err != nil {
                return err
            }
            return ds.Migrate(persistence.SampleModel())
        },
        OnStop: func(ctx context.Context) error {
            return ds.Shutdown()
        },
    })

    return ds, nil
}
```

#### Router (`internal/adapter/inbound/http/router.go`)

```go
// fx 적용 후 - Lifecycle 훅으로 서버 시작/종료
func NewFxRouter(lc fx.Lifecycle, conf *config.Configuration, log *LogrusLogger) (*Router, error) {
    r, err := Init(log.Logger, conf.Server.Debug)
    if err != nil {
        return nil, err
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            addr := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
            go func() {
                if err := r.Run(addr); err != nil {
                    log.Errorf("server error: %s", err)
                }
            }()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            r.Shutdown()
            return nil
        },
    })

    return r, nil
}
```

---

### 변환 예시: Domain 레이어 (Sample)

```go
// bootstrap/routes.go - 라우트 등록 (Invoke 대상)
func RegisterRoutes(
    router *apphttp.Router,
    sample *handler.SampleHandler,
    version *handler.VersionHandler,
) {
    router.GET("/version", version.GetVersion)

    apiv1 := router.Group("/api/v1")
    apiv1.GET("/samples",       sample.GetAll)
    apiv1.GET("/sample/:id",    sample.GetByID)
    apiv1.POST("/sample",       sample.Create)
    apiv1.POST("/sample/update",sample.Update)
    apiv1.DELETE("/sample/:id", sample.Delete)
}
```

---

### 모듈 구성

```go
// internal/bootstrap/modules.go

var InfraModule = fx.Module("infra",
    fx.Provide(
        config.NewEnvironment,
        config.NewConfiguration,
        logger.NewFxLogger,
        db.NewFxDataStore,
        apphttp.NewFxRouter,
    ),
)

var SampleModule = fx.Module("sample",
    fx.Provide(
        persistence.NewSampleRepository,
        usecase.NewGetAllUseCase,
        usecase.NewGetByIDUseCase,
        usecase.NewCreateUseCase,
        usecase.NewUpdateUseCase,
        usecase.NewDeleteUseCase,
        handler.NewSampleHandler,
        handler.NewVersionHandler,
    ),
)

// main.go
func main() {
    fx.New(
        InfraModule,
        SampleModule,
        fx.Invoke(bootstrap.RegisterRoutes),
    ).Run()
}
```

---

## 새 도메인 추가 시 작업량 비교

도메인(예: `User`)을 새로 추가할 때:

| 현재 방식 | fx 방식 |
|----------|--------|
| `bootstrap/wire.go`에 wiring 코드 직접 작성 | `fx.Provide()` 한 줄 추가 |
| 초기화 순서 고려 필요 | 타입 기반 자동 순서 결정 |
| `main.go Context`에 필드 추가 가능성 | main.go 변경 없음 |

```go
// fx 방식: 새 도메인 추가 = Provider 등록만
var UserModule = fx.Module("user",
    fx.Provide(
        persistence.NewUserRepository,
        usecase.NewCreateUserUseCase,
        handler.NewUserHandler,
    ),
)

// main.go에 모듈만 추가
fx.New(
    InfraModule,
    SampleModule,
    UserModule,     // 추가
    fx.Invoke(bootstrap.RegisterRoutes),
).Run()
```

---

## 주의사항

### 런타임 에러

Wire(컴파일 타임)와 달리 fx는 **런타임에 타입 불일치를 탐지**합니다.
`app.Run()` 전에 `app.Err()`로 초기화 에러를 확인할 수 있습니다.

```go
app := fx.New(...)
if err := app.Err(); err != nil {
    log.Fatal(err)
}
app.Run()
```

### 순환 의존성

fx는 순환 의존성을 감지하면 `cycle detected` 에러를 발생시킵니다.
인터페이스를 통해 의존성 방향을 정리해야 합니다.

### 디버깅: 의존성 그래프 확인

```go
fx.New(
    fx.Provide(...),
    fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
        return &fxevent.ZapLogger{Logger: log}
    }),
).Run()
```

또는 `DOT` 형식으로 의존성 그래프를 출력할 수 있습니다.

---

## 도입 시점 권장

| 조건 | 권장 |
|------|------|
| 도메인 1~2개 | 현재 수동 방식 유지 |
| 도메인 3~4개 이상 | fx 도입 검토 |
| 팀 규모가 커지고 모듈 경계가 필요할 때 | fx.Module로 명확한 경계 정의 |

---

## 참고

- [uber-go/fx GitHub](https://github.com/uber-go/fx)
- [fx 공식 문서](https://pkg.go.dev/go.uber.org/fx)
- [dig (fx 내부 엔진)](https://github.com/uber-go/dig)
