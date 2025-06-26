package config

import (
	"balancer/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testTable := []struct {
		name           string
		configPath     string
		errParseConfig error
		expectedConfig MainConfig
	}{
		{
			name:           "valid config",
			configPath:     "testdata/test_config.yaml",
			errParseConfig: nil,
			expectedConfig: MainConfig{
				TickerRateSec: 123,
				DefaultLimits: defaultLimits{
					Capacity:   10,
					RatePerSec: 1,
				},
				BackendConfig: []*backendConfig{
					{
						BackendURL: "http://testhost_backend1:3",
						Config: struct {
							Health struct {
								Method string `yaml:"method"`
								URL    string `yaml:"url"`
							} `yaml:"health"`
						}{
							Health: struct {
								Method string `yaml:"method"`
								URL    string `yaml:"url"`
							}{
								Method: "TEST_METHOD",
								URL:    "/helth_test_url",
							},
						},
					},
					{
						BackendURL: "http://testhost_backend1:4",
						Config: struct {
							Health struct {
								Method string `yaml:"method"`
								URL    string `yaml:"url"`
							} `yaml:"health"`
						}{
							Health: struct {
								Method string `yaml:"method"`
								URL    string `yaml:"url"`
							}{
								Method: "TEST_METHOD",
								URL:    "/helth_test_url",
							},
						},
					},
				},
				ServerConfig: serverConfig{
					Host: "testhost_server",
					Port: "1",
				},
				DbConfig: dbConfig{
					DbHost:     "testhost_db",
					DbPort:     "2",
					DbName:     "test_db_name",
					DbUser:     "test_db_user",
					DbPassword: "1234",
					DbSSL:      "test_db_sslmode",
				},
			},
		},
		{
			name:           "invalid config path",
			configPath:     "invalid_path/config.yaml",
			errParseConfig: model.ErrParseConfig,
			expectedConfig: MainConfig{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, err := InitMainConfig(testCase.configPath, true)

			if testCase.errParseConfig != nil {
				assert.ErrorIs(t, err, testCase.errParseConfig)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedConfig.GetTickerRateSec(), cfg.TickerRateSec)
			assert.Equal(t, testCase.expectedConfig.DefaultLimits, cfg.DefaultLimits)
			assert.Equal(t, testCase.expectedConfig.ServerConfig, cfg.ServerConfig)
			assert.Equal(t, testCase.expectedConfig.DbConfig, cfg.DbConfig)

			assert.Equal(t, len(testCase.expectedConfig.BackendConfig), len(cfg.BackendConfig))
			for i, b := range testCase.expectedConfig.LoadBackendConfig() {
				assert.Equal(t, b.BckndUrl, cfg.BackendConfig[i].BackendURL)
				assert.Equal(t, b.Method, cfg.BackendConfig[i].Config.Health.Method)
				assert.Equal(t, b.HelthUrl, cfg.BackendConfig[i].Config.Health.URL)
			}
		})
	}
}
