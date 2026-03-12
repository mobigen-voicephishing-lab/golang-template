package controllers

// HTTPResponse 공통으로 사용할 응답 메시지
type HTTPResponse struct {
	Code    int         `yaml:"code" json:"code"`
	Message string      `yaml:"message" json:"message"`
	Data    interface{} `yaml:"data" json:"data"`
}

// For Code ( messages 디렉토리에 json파일의 code 값과 동일해야 한다.
const (
	HTTPSuccess int = 0
	HTTPErrCode1000 int = 1000 + iota
	HTTPErrCode1001
)

// For Message
var (
	HTTPErrMsg map[int]string
)

// ReturnError  ...
func (HTTPResponse) ReturnError(code int, msg string) *HTTPResponse {
	return &HTTPResponse{
		Code:    code,
		Message: msg,
	}
}

// ReturnSuccess ...
func (HTTPResponse) ReturnSuccess(result interface{}) *HTTPResponse {
	return &HTTPResponse{
		Code:    HTTPSuccess,
		Message: "Success",
		Data:   result,
	}
}
