package integration_test

import (
	"testing"

	persistence "github.com/mobigen/golang-web-template/internal/adapter/outbound/persistence/gorm"
	"github.com/mobigen/golang-web-template/internal/domain"
	"github.com/mobigen/golang-web-template/internal/testutil"
	"github.com/mobigen/golang-web-template/internal/usecase/sample"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSampleCRUDFlow(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	ds := testutil.MakeTestDataStore(t, log, persistence.SampleModel())
	defer testutil.CloseTestDataStore(t, ds)

	repo := persistence.NewSampleRepository(ds)

	createUC := sample.NewCreateUseCase(repo)
	getAllUC := sample.NewGetAllUseCase(repo)
	getByIDUC := sample.NewGetByIDUseCase(repo)
	updateUC := sample.NewUpdateUseCase(repo)
	deleteUC := sample.NewDeleteUseCase(repo)

	// Create
	created, err := createUC.Execute(&domain.Sample{Name: "test", Desc: "test desc"})
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "test", created.Name)
	assert.Equal(t, "test desc", created.Desc)
	assert.NotZero(t, created.CreateAt)

	// GetAll
	all, err := getAllUC.Execute()
	require.NoError(t, err)
	assert.Equal(t, 1, len(*all))

	// GetByID
	found, err := getByIDUC.Execute(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "test", found.Name)

	// Update
	updated, err := updateUC.Execute(&domain.Sample{ID: created.ID, Name: "updated", Desc: "updated desc"})
	require.NoError(t, err)
	assert.Equal(t, "updated", updated.Name)

	// Verify update persisted
	found2, err := getByIDUC.Execute(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated", found2.Name)
	assert.Equal(t, "updated desc", found2.Desc)

	// Delete
	deleted, err := deleteUC.Execute(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, deleted.ID)

	// Verify not found after delete
	_, err = getByIDUC.Execute(created.ID)
	assert.Error(t, err)

	// GetAll empty
	all2, err := getAllUC.Execute()
	require.NoError(t, err)
	assert.Equal(t, 0, len(*all2))
}
