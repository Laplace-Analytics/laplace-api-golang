package laplace

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

const testConfig = "./test.env"

type ClientTestSuite struct {
	suite.Suite
	Config LaplaceConfiguration
}

func NewClientTestSuite() *ClientTestSuite {
	return &ClientTestSuite{}
}

func (s *ClientTestSuite) SetupTest() {
	repoRoot, err := findModuleRoot()

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

	s.Config = *config
}

func loadConfig() (*LaplaceConfiguration, error) {
	repoRoot, err := findModuleRoot()

	if err != nil {
		return nil, fmt.Errorf("could not find module root: %v", err)
	}

	configPath := filepath.Join(repoRoot, testConfig)

	return LoadGlobal(configPath)
}
