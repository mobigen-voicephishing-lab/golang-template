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
