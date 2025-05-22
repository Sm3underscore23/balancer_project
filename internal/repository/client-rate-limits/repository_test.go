package clientratelimits

import (
	"balancer/internal/model"
	"balancer/internal/repository/client-rate-limits/testhelpers"
	"context"
	"database/sql"
	"embed"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

//go:embed test-migrations/*.sql
var embedMigrations embed.FS

var (
	client = model.ClientLimits{
		ClientId:   "test-client-id",
		Capacity:   10,
		RatePerSec: 1,
	}

	updatedClient = model.ClientLimits{
		ClientId:   "test-client-id",
		Capacity:   100,
		RatePerSec: 10,
	}
)

type repoTestSuite struct {
	pgContainer *testhelpers.PostgresContainer
	repository  repo
	dbForGoose  *sql.DB
	ctx         context.Context
	suite.Suite
}

func (suite *repoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.NewPostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	pool, err := pgxpool.New(suite.ctx, pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repo{db: pool}

	suite.dbForGoose = stdlib.OpenDBFromPool(pool)

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
}

func (suite *repoTestSuite) SetupTest() {
	if err := goose.Up(suite.dbForGoose, "test-migrations"); err != nil {
		log.Fatal(err)
	}
}

func (suite *repoTestSuite) TearDownTest() {
	if err := goose.Down(suite.dbForGoose, "test-migrations"); err != nil {
		log.Fatal(err)
	}
}

func (suite *repoTestSuite) TestCreateCustomer() {
	t := suite.T()
	err := suite.repository.CreateClientLimits(suite.ctx, client)
	assert.NoError(t, err)

	isCliclientExcists, err := suite.repository.IsClientExists(suite.ctx, client.ClientId)
	assert.NoError(t, err)
	assert.True(t, isCliclientExcists)

	getedClient, err := suite.repository.GetClientLimits(suite.ctx, client.ClientId)
	assert.NoError(t, err)
	assert.Equal(t, client, getedClient)
}

func (suite *repoTestSuite) TestGetCustomer() {
	t := suite.T()
	getedClient, err := suite.repository.GetClientLimits(suite.ctx, client.ClientId)
	assert.Error(t, err)
	assert.Equal(t, model.ClientLimits{}, getedClient)

	err = suite.repository.CreateClientLimits(suite.ctx, client)
	assert.NoError(t, err)

	getedClient, err = suite.repository.GetClientLimits(suite.ctx, client.ClientId)
	assert.NoError(t, err)
	assert.Equal(t, client, getedClient)
}

func (suite *repoTestSuite) TestUpdateCustomer() {
	t := suite.T()

	err := suite.repository.CreateClientLimits(suite.ctx, client)
	assert.NoError(t, err)

	err = suite.repository.UpdateClientLimits(suite.ctx, updatedClient)
	assert.NoError(t, err)

	getedClient, err := suite.repository.GetClientLimits(suite.ctx, client.ClientId)
	assert.NoError(t, err)
	assert.Equal(t, updatedClient, getedClient)
}

func (suite *repoTestSuite) TestDeleteCustomer() {
	t := suite.T()

	err := suite.repository.CreateClientLimits(suite.ctx, client)
	assert.NoError(t, err)

	err = suite.repository.DeleteClientLimits(suite.ctx, client.ClientId)
	assert.NoError(t, err)

	isCliclientExcists, err := suite.repository.IsClientExists(suite.ctx, client.ClientId)
	assert.NoError(t, err)
	assert.False(t, isCliclientExcists)
}

func (suite *repoTestSuite) TestIsClientExistsCustomer() {
	t := suite.T()
	isCliclientExcists, err := suite.repository.IsClientExists(suite.ctx, client.ClientId)
	assert.NoError(t, err)
	assert.False(t, isCliclientExcists)

	err = suite.repository.CreateClientLimits(suite.ctx, client)
	assert.NoError(t, err)

	isCliclientExcists, err = suite.repository.IsClientExists(suite.ctx, client.ClientId)
	assert.NoError(t, err)
	assert.True(t, isCliclientExcists)

	err = suite.repository.DeleteClientLimits(suite.ctx, client.ClientId)
	assert.NoError(t, err)

	isCliclientExcists, err = suite.repository.IsClientExists(suite.ctx, client.ClientId)
	assert.NoError(t, err)
	assert.False(t, isCliclientExcists)
}

func (suite *repoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestClientRepository(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}
