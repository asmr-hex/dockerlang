package dockerlang

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	uuid "github.com/satori/go.uuid"
)

const (
	DOCKERLANG_IMAGE = "dockerlang"
	MEMORY_PORT      = "6666"

	COMPUTATION_TYPE_ENV_VAR  = "COMPUTATION_TYPE_ENV_VAR"
	COMPUTATION_VALUE_ENV_VAR = "COMPUTATION_VALUE_ENV_VAR"
	COMPUTATION_DEPS_ENV_VAR  = "COMPUTATION_DEPS_ENV_VAR"
)

var (
	executer *ExecutionEngine
)

type ExecutionEngine struct {
	Docker     *client.Client
	Guillotine string
	Network    string
}

type ExecutionData struct {
	ComputationType string
	Value           string
	Operands        []string
}

// constructs an ExecutionEngine and binds to the globally scoped executer.
func NewExecutionEngine() error {
	// set the API version to use in an environment variable
	// TODO it would be nice to configure based on the docker version
	// a user currently has.... not enough time right now so skipping that.
	err := os.Setenv("DOCKER_API_VERSION", "1.35")
	if err != nil {
		return err
	}

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		// this is probably because the person who is using dockerlang
		// hasn't installed or started docker on their system -____-
		// unclear why anyone would *not* have docker in their life.
		return err
	}

	// define unique network name
	networkName := fmt.Sprintf("dockerlang.%s", uuid.NewV4().String())

	// bind new ExecutionEngine to globally scoped variable
	executer = &ExecutionEngine{
		Docker:     dockerClient,
		Guillotine: "robespierre",
		Network:    networkName,
	}

	// setup container bridge network if one doesn't already exist.
	_, err = executer.Docker.NetworkCreate(
		context.TODO(),
		networkName,
		types.NetworkCreate{},
	)
	if err != nil {
		return err
	}

	// create docker image from Dockerfile
	// NOTE: this didn't work because of relative directories etc.
	// _, err = executer.Docker.ImageBuild(
	// 	context.TODO(),
	// 	nil, // TODO what is this io.Reader about?
	// 	types.ImageBuildOptions{
	// 		Dockerfile: "../Dockerfile",
	// 		Tags:       []string{DOCKERLANG_IMAGE},
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }

	return nil
}

func ShutdownExecutionEngine() error {
	err := executer.Docker.NetworkRemove(context.TODO(), executer.Network)
	if err != nil && !client.IsErrNotFound(err) {
		// something is very wrong here
		panic(err)
	}

	return nil
}

// construct the arguments to the computation about to be run and then create/start
// a new docker container to perform the actual computation!
func (e *ExecutionEngine) Run(d *ExecutionData) (string, error) {
	var ctx = context.Background()

	// construct a comma delimited list of dockerlang container ids
	// which we will pass to the container as an environment variable
	dependencies := strings.Join(d.Operands, ",")

	// create the DockerLang Container Id for this computation
	dlci := uuid.NewV4().String()

	// create docker container for this computation
	_, err := e.Docker.ContainerCreate(
		ctx, // i have no idea what this is or should be
		&container.Config{
			ExposedPorts: nat.PortSet{MEMORY_PORT: struct{}{}},
			Image:        DOCKERLANG_IMAGE,
			Env: []string{
				fmt.Sprintf("%s=%s", COMPUTATION_TYPE_ENV_VAR, d.ComputationType),
				fmt.Sprintf("%s=%s", COMPUTATION_VALUE_ENV_VAR, d.Value),
				fmt.Sprintf("%s=%s", COMPUTATION_DEPS_ENV_VAR, dependencies),
			},
			AttachStdout: true,
			Tty:          true,
		},
		nil,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				e.Network: &network.EndpointSettings{
					NetworkID: e.Network,
				},
			},
		},
		dlci,
	)
	if err != nil {
		return "", err
	}

	// setup stdout stream from container
	hijackedResp, err := e.Docker.ContainerAttach(
		ctx,
		dlci,
		types.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		},
	)
	if err != nil {
		return "", err
	}

	go func() {
		defer hijackedResp.Close()

		io.Copy(os.Stdout, hijackedResp.Reader)
	}()

	// okay lets start the container...
	err = e.Docker.ContainerStart(
		ctx,
		dlci,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return "", err
	}

	// lets return the DockerLang Container Id for this computation
	return dlci, nil
}
