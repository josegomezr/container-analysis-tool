package downloader

import (
	"os"
	"github.com/containers/image/v5/types"
)

func getDockerAuthConfig() types.DockerAuthConfig {
	return types.DockerAuthConfig{
		Username: os.Getenv("REGISTRY_USERNAME"),
		Password: os.Getenv("REGISTRY_PASSWORD"),
	}
}
