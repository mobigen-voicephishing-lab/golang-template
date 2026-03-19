package dto

import "github.com/mobigen/golang-web-template/internal/domain"

// SampleCreateRequest Sample 생성 HTTP 요청 바디
type SampleCreateRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
} // @name SampleCreateRequest

// ToDomain HTTP 요청을 도메인 엔티티로 변환한다.
func (r *SampleCreateRequest) ToDomain() *domain.Sample {
	return &domain.Sample{
		Name: r.Name,
		Desc: r.Desc,
	}
}

// SampleUpdateRequest Sample 수정 HTTP 요청 바디
type SampleUpdateRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
} // @name SampleUpdateRequest

// ToDomain HTTP 요청을 도메인 엔티티로 변환한다.
func (r *SampleUpdateRequest) ToDomain() *domain.Sample {
	return &domain.Sample{
		ID:   r.ID,
		Name: r.Name,
		Desc: r.Desc,
	}
}
