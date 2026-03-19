package dto

// HTTPResponse 공통 API 응답 포맷 (제네릭)
// swagger 모델 이름 단축: 닫는 괄호에 // @name <alias> 를 추가한다.
// 예) HTTPResponse[domain.Sample] → "HTTPResponse-Sample"
type HTTPResponse[T any] struct {
	IsSuccess bool   `json:"isSuccess" yaml:"isSuccess"`
	Code      int    `json:"code"      yaml:"code"`
	Message   string `json:"message"   yaml:"message"`
	Data      T      `json:"data,omitempty" yaml:"data,omitempty"`
} // @name HTTPResponse

// Error Codes
const (
	// 0 - Success
	Success = 0

	// 1000–1999 - Business Logic: 도메인 규칙 위반, 상태 불일치, 중복 등
	ErrAlreadyExists           = 1000
	ErrNotFound                = 1001
	ErrAlreadyProcessed        = 1002
	ErrInvalidStatusTransition = 1003

	// 2000–2999 - Client Request: 요청 형식, 파라미터, 바인딩 오류
	ErrInvalidRequestBody = 2000
	ErrValidationFailed   = 2001
	ErrInvalidParameter   = 2002
	ErrInvalidDateFormat  = 2010

	// 3000–3999 - Resource / Auth: 인증/인가, 리소스 접근 오류
	ErrUnauthorized     = 3000
	ErrForbidden        = 3001
	ErrTokenExpired     = 3002
	ErrInvalidToken     = 3003
	ErrRouteNotFound    = 3010
	ErrMethodNotAllowed = 3011

	// 4000–4999 - External Dependency: 외부 시스템 연동 실패
	ErrDatabaseConnection = 4010
	ErrDatabaseQuery      = 4011

	// 5000–5999 - Internal Server: 예상치 못한 예외, 초기화 실패
	ErrInternalServer       = 5000
	ErrComponentInitFailed  = 5001
	ErrUnexpectedNilPointer = 5002
)

// ErrMessages 에러 코드별 기본 메시지
var ErrMessages = map[int]string{
	// 1000 Business Logic
	ErrAlreadyExists:           "already exists",
	ErrNotFound:                "not found",
	ErrAlreadyProcessed:        "already processed",
	ErrInvalidStatusTransition: "invalid status transition",

	// 2000 Client Request
	ErrInvalidRequestBody: "invalid request body",
	ErrValidationFailed:   "validation failed",
	ErrInvalidParameter:   "invalid parameter",
	ErrInvalidDateFormat:  "invalid date format",

	// 3000 Resource / Auth
	ErrUnauthorized:     "unauthorized",
	ErrForbidden:        "forbidden",
	ErrTokenExpired:     "token expired",
	ErrInvalidToken:     "invalid token",
	ErrRouteNotFound:    "route not found",
	ErrMethodNotAllowed: "method not allowed",

	// 4000 External Dependency
	ErrDatabaseConnection: "database connection failed",
	ErrDatabaseQuery:      "database query failed",

	// 5000 Internal Server
	ErrInternalServer:       "internal server error",
	ErrComponentInitFailed:  "failed to initialize component",
	ErrUnexpectedNilPointer: "unexpected nil pointer",
}
