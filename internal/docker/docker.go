package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Dockers interface {
	PullAndStart(string, string) error
	Create(string, string, string, []string, []string) error
	Start(string) error
	Stop(string) error
	Restart(string) error
	Remove(string) error
	ListNetworks() []types.NetworkResource
	ListContainers(types.ContainerListOptions) []types.Container
	RunDatabaseDump() error
	RunDatabaseRestore()
}

// Docker represents a docker client
type Docker struct {
	cli        *client.Client
	ctx        context.Context
	auth       types.AuthConfig
	logger     *structs.AppLogger
	authString string
}

func NewDocker(username, password, serverAddr string, logger *structs.AppLogger) (*Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	log.Info().Str("module", "docker").Msg("connected to the docker daemon")

	authConfig := types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: serverAddr,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	return &Docker{cli: cli, ctx: ctx, auth: authConfig, logger: logger, authString: authStr}, nil
}

func (d *Docker) Ping() {
	ping, err := d.cli.Ping(d.ctx)

	if err != nil {
		log.Error().Str("module", "docker").Err(err).Msg("error pinging docker")
	}
	log.Info().Str("module", "docker").Str("Api-version", ping.APIVersion).Msg("docker daemon health check")
}

func (d *Docker) RunDatabaseDump() (err error) {
	args, err := filters.ParseFlag("label=com.docker.compose.service=mariadb", filters.NewArgs())

	if err != nil {
		d.logger.SystemLogger.Error(err, "error prepping args")
		return
	}

	runningContainers := d.ListContainers(types.ContainerListOptions{All: true, Filters: args})

	containerId := runningContainers[0].ID

	c := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd: []string{"/bin/sh",
			"-c",
			fmt.Sprintf(
				os.ExpandEnv(
					"mysqldump -u ${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE} %s > ${DB_BACKUP_DIR}${DB_BACKUP_FILE}"),
				persistence.ElementsTable,
			),
		},
	}

	id, err := d.cli.ContainerExecCreate(d.ctx, containerId, c)

	if err != nil {
		d.logger.SystemLogger.Error(err, fmt.Sprintf("error creating exec to container: %s", containerId))
		return
	}

	resp, err := d.cli.ContainerExecAttach(d.ctx, id.ID, types.ExecConfig{})

	if err != nil {
		d.logger.SystemLogger.Error(err, fmt.Sprintf("error attaching exec to container: %s", containerId))
		return
	}

	defer resp.Close()

	err = d.cli.ContainerExecStart(d.ctx, id.ID, types.ExecStartCheck{})
	if err != nil {
		d.logger.SystemLogger.Error(err, fmt.Sprintf("error starting exec to container: %s", containerId))
		return
	}

	return nil
}

func (d *Docker) RunDatabaseRestore() {
	args, err := filters.ParseFlag("label=com.docker.compose.service=mariadb", filters.NewArgs())

	if err != nil {
		d.logger.SystemLogger.Error(err, "error prepping args")
	}

	runningContainers := d.ListContainers(types.ContainerListOptions{Filters: args})

	containerId := runningContainers[0].ID

	c := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd: []string{"/bin/sh",
			"-c",
			os.ExpandEnv("mysql -u ${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE} < ${DB_BACKUP_DIR}${DB_BACKUP_FILE}"),
		}}

	id, err := d.cli.ContainerExecCreate(d.ctx, containerId, c)

	if err != nil {
		d.logger.SystemLogger.Error(err, fmt.Sprintf("error creating exec to container: %s", containerId))
	}

	resp, err := d.cli.ContainerExecAttach(d.ctx, id.ID, types.ExecConfig{})

	if err != nil {
		d.logger.SystemLogger.Error(err, fmt.Sprintf("error attaching exec to container: %s", containerId))
	}

	defer resp.Close()

	err = d.cli.ContainerExecStart(d.ctx, id.ID, types.ExecStartCheck{})
	if err != nil {
		d.logger.SystemLogger.Error(err, fmt.Sprintf("error starting exec to container: %s", containerId))
	}
}

func (d *Docker) ListContainers(options types.ContainerListOptions) []types.Container {
	containers, err := d.cli.ContainerList(d.ctx, options)
	if err != nil {
		log.Error().Str("module", "docker").Err(err).Msg("error listing containers")
	}
	return containers
}

func (d *Docker) ListNetworks() []types.NetworkResource {
	networks, err := d.cli.NetworkList(d.ctx, types.NetworkListOptions{})
	if err != nil {
		log.Error().Str("module", "docker").Err(err).Msg("error listing networks")
	}
	return networks
}

func (d *Docker) Stop(containerID string) error {
	timeout := time.Until(time.Now().Add(30 * time.Second))
	if err := d.cli.ContainerStop(d.ctx, containerID, &timeout); err != nil {
		log.Error().Str("module", "docker").Str("containerID", containerID).Err(err).Msg("error stopping container")
		return err
	}
	return nil
}

func (d *Docker) Start(containerID string) error {
	if err := d.cli.ContainerStart(d.ctx, containerID, types.ContainerStartOptions{}); err != nil {
		log.Error().Str("module", "docker").Str("containerID", containerID).Err(err).Msg("error starting container")
		return err
	}
	return nil
}

func (d *Docker) Restart(containerID string) error {
	if err := d.cli.ContainerRestart(d.ctx, containerID, nil); err != nil {
		log.Error().Str("module", "docker").Str("containerID", containerID).Err(err).Msg("error restarting container")
		return err
	}
	return nil
}

func (d *Docker) Remove(containerID string) error {
	if err := d.Stop(containerID); err != nil {
		log.Error().Str("module", "docker").Str("containerID", containerID).Err(err).Msg("error stopping container")
		return err
	}
	if err := d.cli.ContainerRemove(d.ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		log.Error().Str("module", "docker").Str("containerID", containerID).Err(err).Msg("error removing container")
		return err
	}
	return nil
}

func (d *Docker) Create(imageRef, containerName, containerNetwork string, volumes []string, envVars []string) error {
	out, err := d.cli.ImagePull(d.ctx, imageRef, types.ImagePullOptions{RegistryAuth: d.authString})
	if err != nil {
		return err
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	networkConfig := &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{}}

	gatewayConfig := &network.EndpointSettings{
		NetworkID: containerNetwork,
		Aliases:   []string{containerName},
	}

	networkConfig.EndpointsConfig[containerNetwork] = gatewayConfig

	resp, err := d.cli.ContainerCreate(d.ctx, &container.Config{
		Image: imageRef,
		Tty:   false,
		Env:   envVars,
	}, &container.HostConfig{
		Binds: volumes,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}, networkConfig, containerName)

	if err != nil {
		return err
	}

	if err := d.cli.ContainerStart(d.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (d *Docker) Login() (registry.AuthenticateOKBody, error) {
	authStatus, err := d.cli.RegistryLogin(d.ctx, types.AuthConfig{
		Username:      d.auth.Username,
		Password:      d.auth.Password,
		ServerAddress: d.auth.ServerAddress,
	})

	return authStatus, err
}

func (d *Docker) PullAndStart(imageRef, createdContainerName string) error {
	out, err := d.cli.ImagePull(d.ctx, imageRef, types.ImagePullOptions{RegistryAuth: d.authString})
	if err != nil {
		return err
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := d.cli.ContainerCreate(d.ctx, &container.Config{
		Image: imageRef,
		Tty:   false,
	}, nil, nil, createdContainerName)
	if err != nil {
		return err
	}

	if err := d.cli.ContainerStart(d.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (d *Docker) Inspect(containerID string) (*types.ContainerJSON, error) {
	ctr, err := d.cli.ContainerInspect(d.ctx, containerID)
	if err != nil {
		log.Error().Str("module", "docker").Str("containerID", containerID).Err(err).Msg("error inspecting container")
		return nil, err
	}
	return &ctr, nil
}

func (d *Docker) Logs(containerID string, tail string) ([]string, error) {
	var (
		bytes          = make([]byte, 3000) // Telegram message length limit
		logs  []string = nil
		err   error    = nil
	)
	if tail != "all" && !isNumber(tail) {
		tail = "10"
	}
	logsReader, err := d.cli.ContainerLogs(d.ctx, containerID, types.ContainerLogsOptions{Tail: tail, ShowStderr: true, ShowStdout: true})
	if err != nil {
		log.Error().Str("module", "docker").Str("containerID", containerID).Err(err).Msg("error getting container logs")
		return nil, err
	}
	defer func() {
		err := logsReader.Close()
		if err != nil {
			log.Error().Str("module", "docker").Str("containerID", containerID).Err(err).Msg("error closing io.Reader")
		}
	}()

	for {
		numBytes, err := logsReader.Read(bytes)
		logs = append(logs, string(bytes[:numBytes]))
		if err == io.EOF {
			break
		}
	}
	return logs, nil
}

func (d *Docker) isValidID(containerID string) bool {
	re := regexp.MustCompile(`(?m)^[A-Fa-f0-9]{10,12}$`)
	return re.MatchString(containerID)
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
