package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ozsari/velour/internal/models"
)

type DockerManager struct {
	client  *client.Client
	dataDir string
}

func NewDockerManager(dataDir string) (*DockerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker: %w", err)
	}

	return &DockerManager{
		client:  cli,
		dataDir: dataDir,
	}, nil
}

func (dm *DockerManager) Install(ctx context.Context, def *models.ServiceDefinition) error {
	// Pull image
	reader, err := dm.client.ImagePull(ctx, def.Image, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", def.Image, err)
	}
	defer reader.Close()
	// Drain the reader to complete the pull
	buf := make([]byte, 1024)
	for {
		_, err := reader.Read(buf)
		if err != nil {
			break
		}
	}

	// Prepare port bindings
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, p := range def.Ports {
		containerPort := nat.Port(fmt.Sprintf("%d/%s", p.Container, p.Protocol))
		portBindings[containerPort] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", p.Host)},
		}
		exposedPorts[containerPort] = struct{}{}
	}

	// Prepare mounts
	var mounts []mount.Mount
	for _, v := range def.Volumes {
		hostPath := strings.ReplaceAll(v.Host, "${DATA_DIR}", dm.dataDir)
		os.MkdirAll(hostPath, 0755)
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: hostPath,
			Target: v.Container,
		})
	}

	// Prepare env
	var envList []string
	for k, v := range def.Env {
		envList = append(envList, fmt.Sprintf("%s=%s", k, v))
	}

	containerName := fmt.Sprintf("velour-%s", def.ID)

	// Create container
	resp, err := dm.client.ContainerCreate(ctx,
		&container.Config{
			Image:        def.Image,
			Env:          envList,
			ExposedPorts: exposedPorts,
			Labels: map[string]string{
				"velour.managed": "true",
				"velour.service": def.ID,
			},
		},
		&container.HostConfig{
			PortBindings:  portBindings,
			Mounts:        mounts,
			RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
		},
		&network.NetworkingConfig{},
		nil,
		containerName,
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := dm.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func (dm *DockerManager) Start(ctx context.Context, serviceID string) error {
	containerName := fmt.Sprintf("velour-%s", serviceID)
	return dm.client.ContainerStart(ctx, containerName, container.StartOptions{})
}

func (dm *DockerManager) Stop(ctx context.Context, serviceID string) error {
	containerName := fmt.Sprintf("velour-%s", serviceID)
	timeout := 30
	return dm.client.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout})
}

func (dm *DockerManager) Restart(ctx context.Context, serviceID string) error {
	containerName := fmt.Sprintf("velour-%s", serviceID)
	timeout := 30
	return dm.client.ContainerRestart(ctx, containerName, container.StopOptions{Timeout: &timeout})
}

func (dm *DockerManager) Remove(ctx context.Context, serviceID string) error {
	containerName := fmt.Sprintf("velour-%s", serviceID)
	return dm.client.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
}

func (dm *DockerManager) Status(ctx context.Context, serviceID string) (models.ServiceStatus, error) {
	containerName := fmt.Sprintf("velour-%s", serviceID)
	info, err := dm.client.ContainerInspect(ctx, containerName)
	if err != nil {
		return models.StatusStopped, nil
	}

	if info.State.Running {
		return models.StatusRunning, nil
	}
	return models.StatusStopped, nil
}

func (dm *DockerManager) ListManaged(ctx context.Context) ([]models.Service, error) {
	containers, err := dm.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var services []models.Service
	for _, c := range containers {
		if c.Labels["velour.managed"] != "true" {
			continue
		}

		serviceID := c.Labels["velour.service"]
		def := FindByID(serviceID)
		if def == nil {
			continue
		}

		status := models.StatusStopped
		if c.State == "running" {
			status = models.StatusRunning
		}

		webURL := ""
		if len(def.Ports) > 0 {
			webURL = fmt.Sprintf("http://localhost:%d", def.Ports[0].Host)
		}

		services = append(services, models.Service{
			ID:          def.ID,
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Category:    def.Category,
			Port:        def.Ports[0].Host,
			WebURL:      webURL,
			Status:      status,
			Type:        "docker",
			Image:       def.Image,
			Installed:   true,
		})
	}

	return services, nil
}

// Exec runs a command inside a running container and returns the output.
// This is the key method that enables cross-container post-processing:
// e.g. running FileBot AMC inside the filebot container when qBittorrent finishes.
func (dm *DockerManager) Exec(ctx context.Context, containerName string, cmd []string) (string, error) {
	execConfig := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := dm.client.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec: %w", err)
	}

	resp, err := dm.client.ContainerExecAttach(ctx, execID.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to attach exec: %w", err)
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %w", err)
	}

	// Check exit code
	inspect, err := dm.client.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return string(output), fmt.Errorf("failed to inspect exec: %w", err)
	}

	if inspect.ExitCode != 0 {
		return string(output), fmt.Errorf("command exited with code %d", inspect.ExitCode)
	}

	return string(output), nil
}

func (dm *DockerManager) GetDataDir() string {
	return filepath.Clean(dm.dataDir)
}
