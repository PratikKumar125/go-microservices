package cli

import (
	"github.com/PratikKumar125/go-microservices/pkg/logging"
)

type Config struct {
	MigrationsDir string
	Logger        *logging.Logger
}
