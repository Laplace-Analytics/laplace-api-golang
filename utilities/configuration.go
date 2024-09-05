package utilities

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type LaplaceConfiguration struct {
	BaseURL string `split_words:"true"`
	APIKey  string `split_words:"true"`
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

func LoadGlobal(filename string) (*LaplaceConfiguration, error) {
	cfg, err := loadGlobal(filename, validationFuncRegular)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadGlobal(filename string, fn validationFunc) (*LaplaceConfiguration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(LaplaceConfiguration)

	// although the package is called "auth" it used to be called "gotrue"
	// so environment configs will remain to be called "GOTRUE"
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	if err := config.ApplyDefaults(); err != nil {
		return nil, err
	}

	if fn != nil {
		if err := fn(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

type validationFunc func(*LaplaceConfiguration) error

func validationFuncRegular(config *LaplaceConfiguration) error {
	return config.Validate()
}

func (c *LaplaceConfiguration) Validate() error {
	return nil
}

func (c *LaplaceConfiguration) ApplyDefaults() error {
	return nil
}
