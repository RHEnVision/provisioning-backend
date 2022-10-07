package routes

import (
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/version"
)

func PathPrefix() string {
	return fmt.Sprintf("/api/%s/%s", version.APIPathName, version.APIPathVersion)
}
