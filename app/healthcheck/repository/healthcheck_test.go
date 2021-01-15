package repository

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	HealthCheckInterface "backend-sample-la/app/healthcheck"
	"backend-sample-la/models"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Rosaniline/gorm-ut/pkg/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository HealthCheckInterface.IHealthCheckRepository
	person     *model.Person
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("mysql", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(false)

	s.repository = NewHealthCheckRepository(s.DB)
}

func (s *Suite) TestGetDBTimestampType() {
	dbTimestamp := s.repository.GetDBTimestamp()
	var healthCheck models.HealthCheck
	assert.IsType(s.T(), healthCheck, dbTimestamp)
}

func (s *Suite) TestRepositoryIsSingleton() {
	newRepository := NewHealthCheckRepository(s.DB)
	require.Equal(s.T(), newRepository, s.repository)
}

func (s *Suite) TestTimestampIsExecuted() {

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("SELECT current_timestamp")
	s.mock.ExpectCommit().
		WillReturnError(nil)

	dbTimestamp := s.repository.GetDBTimestamp()
	assert.NotNil(s.T(), dbTimestamp)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHealthCheckRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}
