package control

import (
	"context"
	"errors"
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

var (
	LABEL_ENABLE      = "chief.enable"
	LABEL_ID          = "chief.id"
	LABEL_NAME        = "chief.name"
	LABEL_DESCRIPTION = "chief.description"
)

var LABEL_LIST = []string{LABEL_ENABLE, LABEL_DESCRIPTION, LABEL_NAME, LABEL_ID}

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

type ChiefContainer struct {
	ID          string
	DockerID    string
	Name        string
	Description string
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

func (cs *ChiefService) GetInfo(container types.Container) ChiefContainer {
	labels := container.Labels

	return ChiefContainer{
		ID:          labels[LABEL_ID],
		DockerID:    container.ID,
		Name:        labels[LABEL_NAME],
		Description: labels[LABEL_DESCRIPTION],
	}
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

var ErrContainerNotFound = errors.New("container not found")

func (cs *ChiefService) GetContainerId(id string) (string, error) {
	containers, err := cs.GetAllContainers()
	if err != nil {
		return "", err
	}

	for _, container := range containers {
		chiefContainer := cs.GetInfo(container)
		if chiefContainer.ID == id {
			return chiefContainer.ID, nil
		}
	}

	return "", ErrContainerNotFound
}
