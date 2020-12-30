package cmccloud

import (
	"log"

	"github.com/cmccloud/gocmcapi"
)

// Config object
type Config struct {
	APIKey           string
	APIEndpoint      string
	TerraformVersion string
}

// CombinedConfig struct
type CombinedConfig struct {
	client *gocmcapi.Client
}

func (c *CombinedConfig) goCMCClient() *gocmcapi.Client { return c.client }

// Client config
func (c *Config) Client() (*CombinedConfig, error) {
	client, err := gocmcapi.NewClient(c.APIKey)
	if err != nil {
		log.Fatal(err)
	}
	return &CombinedConfig{
		client: client,
	}, nil
}
