package chief_control

import (
	"context"
	"log/slog"
	"slices"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Default docker engine API version
var DEFAULT_API_VERSION = "1.24"

// Indicate that we can control this container
var CONTROL_LABEL = "chief.enable=true"

var LABEL_LIST = []string{"chief.enable", "chief.description", "chief.name", "chief.id"}

type ChiefService struct {
	Cli *client.Client
	Ctx context.Context
}

func NewChiefService() *ChiefService {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithVersion(DEFAULT_API_VERSION))
	if err != nil {
		panic(err)
	}

	slog.Info("Successfully connected to docker.sock")

	return &ChiefService{
		Cli: cli,
		Ctx: ctx,
	}
}

func (cs *ChiefService) GetAllContainers() ([]types.Container, error) {
	containerList, err := cs.Cli.ContainerList(cs.Ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("label", CONTROL_LABEL)),
	})
	if err != nil {
		slog.Error("Error finding containers", "error", err)
		return nil, err
	}

	return containerList, nil
}

type ContainerLabel struct {
	Key   string
	Value string
}

func (cs *ChiefService) GetLabels(container types.Container) []ContainerLabel {
	labels := container.Labels

	labelList := []ContainerLabel{}

	for key, value := range labels {

		// Check if labelKey is in LABEL_LIST and labelValue is not empty
		exist := slices.Contains(LABEL_LIST, key)

		if exist {
			labelList = append(labelList, ContainerLabel{Key: key, Value: value})
		}

	}

	return labelList
}

func (cs *ChiefService) RestartContainer(id string) error {
	slog.Info("Restarting container", "id", id)
	err := cs.Cli.ContainerRestart(cs.Ctx, id, container.StopOptions{Signal: "", Timeout: nil})
	if err != nil {
		slog.Warn("Failed to restart container", "id", id, "error", err)
		return err
	}

	slog.Info("Successfully restarted container", "id", id)
	return nil
}
