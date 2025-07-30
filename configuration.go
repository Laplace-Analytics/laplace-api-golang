package laplace

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type LaplaceConfiguration struct {
	APIKey  string `split_words:"true"`
	BaseURL string `split_words:"true"`
}

func loadEnvironment(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Overload(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exist, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

// LoadGlobal loads configuration from environment variables and optionally from a .env file.
func LoadGlobal(filename string) (*LaplaceConfiguration, error) {
	cfg, err := loadGlobal(filename)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadGlobal(filename string) (*LaplaceConfiguration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(LaplaceConfiguration)

	// although the package is called "auth" it used to be called "gotrue"
	// so environment configs will remain to be called "GOTRUE"
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	config.ApplyDefaults()

	return config, nil
}

// Validate performs validation checks on the configuration.
func (c *LaplaceConfiguration) Validate() error {
	if c.APIKey == "" {
		return errors.New("API key is required")
	}

	return nil
}

// ApplyDefaults sets default values for configuration fields that are not provided.
func (c *LaplaceConfiguration) ApplyDefaults() {
	if c.BaseURL == "" {
		c.BaseURL = BaseURL
	}

}
