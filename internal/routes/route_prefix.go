package routes

import (
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
)

func PathPrefix() string {
	return fmt.Sprintf("/api/%s/%s", config.Application.Name, config.Application.Version)
}
