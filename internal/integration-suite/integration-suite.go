package integrationsuite

import (
	"balancer/internal/integration-suite/testhelpers"
	"balancer/internal/repository"
	clientratelimits "balancer/internal/repository/client-rate-limits"
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
)

const (
	migrationsDir = "./migrations"
)

type RepoTestSuite struct {
	PgContainer *testhelpers.PostgresContainer
	Repo        repository.LimitsRepository
	DbForGoose  *sql.DB
	Ctx         context.Context
	suite.Suite
}

func (s *RepoTestSuite) SetupSuite() {
	s.Ctx = context.Background()
	pgContainer, err := testhelpers.NewPostgresContainer(s.Ctx)
	if err != nil {
		log.Fatal(err)
	}
	s.PgContainer = pgContainer

	pool, err := pgxpool.New(s.Ctx, pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	s.Repo = clientratelimits.NewClientRateLimitsRepo(pool)

	s.DbForGoose = stdlib.OpenDBFromPool(pool)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	
	err = os.Chdir("../../..")
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := goose.UpContext(s.Ctx, s.DbForGoose, migrationsDir); err != nil {
		log.Fatalf("goose up failed: %v", err)
	}
}

func (s *RepoTestSuite) SetupTest() {
	if err := goose.Up(s.DbForGoose, migrationsDir); err != nil {
		log.Fatal(err)
	}
}

func (s *RepoTestSuite) TearDownTest() {
	if err := goose.Down(s.DbForGoose, migrationsDir); err != nil {
		log.Fatal(err)
	}
}

func (s *RepoTestSuite) TearDownSuite() {
	if err := s.PgContainer.Terminate(s.Ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}
