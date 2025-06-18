package clientratelimits_test

import (
	integrationsuite "balancer/internal/integration-suite"
	"balancer/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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
	integrationsuite.RepoTestSuite
}

func (s *repoTestSuite) TestCreateCustomer() {
	t := s.T()
	err := s.Repo.CreateClientLimits(s.Ctx, client)
	s.Suite.NoError(err)

	isCliclientExcists, err := s.Repo.IsClientExists(s.Ctx, client.ClientId)
	s.Suite.NoError(err)
	assert.True(t, isCliclientExcists)

	getedClient, err := s.Repo.GetClientLimits(s.Ctx, client.ClientId)
	s.Suite.NoError(err)
	assert.Equal(t, client, getedClient)
}

func (s *repoTestSuite) TestGetCustomer() {
	t := s.T()
	getedClient, err := s.Repo.GetClientLimits(s.Ctx, client.ClientId)
	s.Suite.Error(err)
	assert.Equal(t, model.ClientLimits{}, getedClient)

	err = s.Repo.CreateClientLimits(s.Ctx, client)
	s.Suite.NoError(err)

	getedClient, err = s.Repo.GetClientLimits(s.Ctx, client.ClientId)
	s.Suite.NoError(err)
	assert.Equal(t, client, getedClient)
}

func (s *repoTestSuite) TestUpdateCustomer() {
	t := s.T()

	err := s.Repo.CreateClientLimits(s.Ctx, client)
	s.Suite.NoError(err)

	err = s.Repo.UpdateClientLimits(s.Ctx, updatedClient)
	s.Suite.NoError(err)

	getedClient, err := s.Repo.GetClientLimits(s.Ctx, client.ClientId)
	s.Suite.NoError(err)
	assert.Equal(t, updatedClient, getedClient)
}

func (s *repoTestSuite) TestDeleteCustomer() {
	t := s.T()

	err := s.Repo.CreateClientLimits(s.Ctx, client)
	s.Suite.NoError(err)

	err = s.Repo.DeleteClientLimits(s.Ctx, client.ClientId)
	s.Suite.NoError(err)

	isCliclientExcists, err := s.Repo.IsClientExists(s.Ctx, client.ClientId)
	s.Suite.NoError(err)
	assert.False(t, isCliclientExcists)
}

func (s *repoTestSuite) TestIsClientExistsCustomer() {
	t := s.T()
	isCliclientExcists, err := s.Repo.IsClientExists(s.Ctx, client.ClientId)
	s.Suite.NoError(err)
	assert.False(t, isCliclientExcists)

	err = s.Repo.CreateClientLimits(s.Ctx, client)
	s.Suite.NoError(err)

	isCliclientExcists, err = s.Repo.IsClientExists(s.Ctx, client.ClientId)
	s.Suite.NoError(err)
	assert.True(t, isCliclientExcists)

	err = s.Repo.DeleteClientLimits(s.Ctx, client.ClientId)
	s.Suite.NoError(err)

	isCliclientExcists, err = s.Repo.IsClientExists(s.Ctx, client.ClientId)
	s.Suite.NoError(err)
	assert.False(t, isCliclientExcists)
}

func TestClientRepository(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}
