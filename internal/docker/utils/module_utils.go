package utils

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func BuildModuleEnvVars(vars *[]string) {
	tokenVar := fmt.Sprintf("INTERNAL_TOKEN=%s", viper.GetString("internaltoken"))
	*vars = append(*vars, tokenVar)

	controllerVar := fmt.Sprintf("CONTROLLER_SVC_NAME=%s", os.Getenv("CONTROLLER_SVC_NAME"))
	*vars = append(*vars, controllerVar)

	controllerPortVar := fmt.Sprintf("CONTROLLER_PORT=%s", os.Getenv("CONTROLLER_PORT"))
	*vars = append(*vars, controllerPortVar)

	hostnameVar := fmt.Sprintf("HOST_DOMAIN=%s", os.Getenv("HOST_DOMAIN"))
	*vars = append(*vars, hostnameVar)
}

func AddModuleBindPaths(binds *[]string, moduleId string) {
	projectRoot := os.Getenv("PROJECT_ROOT")

	for i, bind := range *binds {
		(*binds)[i] = fmt.Sprintf("%s/%s/%s", projectRoot, moduleId, bind)
	}
}

func AddModuleNetwork(networks []types.NetworkResource) (string, error) {
	for _, val := range networks {
		if strings.HasSuffix(val.Name, os.Getenv("MODULE_NETWORK_NAME")) {
			return val.Name, nil
		}
	}

	return "", errors.New("couldn't extract name of network")
}

func ValidImageRef(imageRef string) bool {
	return strings.HasPrefix(imageRef, os.Getenv("DOCKER_PREFIX"))
}
