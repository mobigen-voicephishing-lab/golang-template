package domain

// Sample 도메인 엔티티 (순수 구조체, 프레임워크 의존 없음)
// swagger:model 이름을 단축하려면 닫는 괄호에 // @name <alias> 를 추가한다.
type Sample struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	CreateAt int64  `json:"createAt"`
} // @name Sample
