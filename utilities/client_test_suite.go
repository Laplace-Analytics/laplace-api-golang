package utilities

import (
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

const testConfig = "./utilities/test.env"

type ClientTestSuite struct {
	suite.Suite
	Config LaplaceConfiguration
}

func NewClientTestSuite() *ClientTestSuite {
	return &ClientTestSuite{}
}

func (s *ClientTestSuite) SetupTest() {
	repoRoot, err := FindModuleRoot()

	if err != nil {
		s.T().Fatalf("Could not find module root: %v", err)
	}

	configPath := filepath.Join(repoRoot, testConfig)
	godotenv.Load(configPath)

	config, err := LoadGlobal(configPath)
	if err != nil {
		s.T().Fatalf("Error loading config: %v", err)
	}

	if config.APIKey == "" {
		s.T().Fatalf("API key is not set")
	}

	if config.BaseURL == "" {
		s.T().Fatalf("API base URL is not set")
	}

	s.Config = *config
}
