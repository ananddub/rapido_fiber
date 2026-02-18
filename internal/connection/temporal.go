package connection

import (
	"fmt"

	"encore.app/internal/config"

	"go.temporal.io/sdk/client"
)

func NewTemporalClient(cfg *config.Config) (client.Client, error) {
	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%d", cfg.TemporalHost, cfg.TemporalPort),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	return c, nil
}
