package config

// 빌드 타임에 설정되는 변수 들
// ldflags -X 옵션으로 설정 됨
var (
	Name      string = "web-platform"
	Version   string = "-"
	BuildHash string = "-"
)

// VersionInfo app version info
// swagger 모델 이름 단축: 닫는 괄호에 // @name <alias> 를 추가한다.
type VersionInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildHash string `json:"buildHash"`
} // @name VersionInfo
