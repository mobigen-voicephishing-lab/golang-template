package dto

// VersionResponse 버전 정보 HTTP 응답 DTO
// swagger 모델 이름 단축: 닫는 괄호에 // @name <alias> 를 추가한다.
type VersionResponse struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildHash string `json:"buildHash"`
} // @name VersionResponse
