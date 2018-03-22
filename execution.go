package dockerlang

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	uuid "github.com/satori/go.uuid"
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
	Operands []DLCI
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

func (e *ExecutionEngine) Run(d *ExecutionData) (DLCI, error) {
	// start container with network name

	// pass data structure needed to compute

	return "", nil
}
