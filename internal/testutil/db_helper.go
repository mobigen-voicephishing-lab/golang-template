package testutil

import (
	"testing"

	"github.com/mobigen/golang-web-template/internal/infrastructure/config"
	"github.com/mobigen/golang-web-template/internal/infrastructure/db"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MakeTestDataStore creates an in-memory SQLite DataStore and migrates the given models.
func MakeTestDataStore(tb testing.TB, log *logrus.Logger, models ...interface{}) *db.DataStore {
	tb.Helper()

	tmpDir := tb.TempDir()
	ds, err := db.DataStore{}.New(tmpDir, log)
	require.NoError(tb, err)

	conf := &config.DatastoreConfiguration{
		Database: config.Sqlite,
		Endpoint: config.EndpointInfo{
			Path:   "tmp.db",
			Option: "mode=memory&cache=shared",
		},
		Debug: config.DatastoreDebug{
			LogLevel:      "info",
			SlowThreshold: "1sec",
		},
	}
	err = ds.Connect(conf)
	require.NoError(tb, err)

	err = ds.Migrate(models...)
	require.NoError(tb, err)

	return ds
}

// CloseTestDataStore closes the database connection after testing.
func CloseTestDataStore(tb testing.TB, ds *db.DataStore) {
	tb.Helper()

	d, err := ds.Orm.DB()
	assert.NoError(tb, err)
	err = d.Close()
	assert.NoError(tb, err)
}
