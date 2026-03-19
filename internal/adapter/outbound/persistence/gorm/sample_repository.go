package gorm

import (
	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/infrastructure/db"
	"github.com/mobigen/golang-web-template/util"
)

// sampleModel GORM 전용 모델 (domain.Sample과 분리)
type sampleModel struct {
	ID       int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name     string `gorm:"column:name;not null"`
	Desc     string `gorm:"column:desc;size:256"`
	CreateAt int64  `gorm:"column:createAt;autoCreateTime:milli"`
}

func (sampleModel) TableName() string { return "samples" }

func toDomain(m *sampleModel) *domain.Sample {
	return &domain.Sample{
		ID:       m.ID,
		Name:     m.Name,
		Desc:     m.Desc,
		CreateAt: m.CreateAt,
	}
}

func toModel(d *domain.Sample) *sampleModel {
	return &sampleModel{
		ID:       d.ID,
		Name:     d.Name,
		Desc:     d.Desc,
		CreateAt: d.CreateAt,
	}
}

func toDomainSlice(models []sampleModel) *[]domain.Sample {
	result := make([]domain.Sample, len(models))
	for i, m := range models {
		result[i] = *toDomain(&m)
	}
	return &result
}

// SampleRepository GORM 기반 Sample 리포지토리 구현체
type SampleRepository struct {
	ds *db.DataStore
}

// NewSampleRepository SampleRepository 생성자
func NewSampleRepository(ds *db.DataStore) *SampleRepository {
	return &SampleRepository{ds: ds}
}

// SampleModel returns the GORM model for migration registration
func SampleModel() interface{} {
	return &sampleModel{}
}

// GetAll get all sample from database
func (repo *SampleRepository) GetAll() (*[]domain.Sample, error) {
	var dst []sampleModel
	result := repo.ds.Orm.Find(&dst)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainSlice(dst), nil
}

// GetByID get sample whose id matches
func (repo *SampleRepository) GetByID(id int) (*domain.Sample, error) {
	var dst sampleModel
	result := repo.ds.Orm.Where(&sampleModel{ID: id}).Find(&dst)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, domain.NewNotFoundError("sample not found")
	}
	return toDomain(&dst), nil
}

// Create create sample
func (repo *SampleRepository) Create(input *domain.Sample) (*domain.Sample, error) {
	m := toModel(input)
	m.CreateAt = util.GetMillis()
	result := repo.ds.Orm.Create(m)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomain(m), nil
}

// Update update sample
func (repo *SampleRepository) Update(input *domain.Sample) (*domain.Sample, error) {
	m := toModel(input)
	result := repo.ds.Orm.Model(m).
		Where(&sampleModel{ID: m.ID}).
		Updates(map[string]interface{}{
			"name": m.Name,
			"desc": m.Desc,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomain(m), nil
}

// Delete delete sample from id(primaryKey)
func (repo *SampleRepository) Delete(id int) (*domain.Sample, error) {
	var dst sampleModel
	result := repo.ds.Orm.Where(&sampleModel{ID: id}).Find(&dst)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, domain.NewNotFoundError("sample not found")
	}
	result = repo.ds.Orm.Delete(&sampleModel{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomain(&dst), nil
}
