package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gopkg.in/yaml.v2"
)

// Config holds the YAML structure
type Config struct {
	Services map[string]Service `yaml:"services"`
}

// Service defines a single container setup
type Service struct {
	Image string            `yaml:"image"`
	Ports map[string]string `yaml:"ports"` // e.g., "5432/tcp": "5432"
	Env   []string          `yaml:"env"`   // e.g., ["POSTGRES_PASSWORD=secret"]
}

func parseConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal YAML: %w", err)
	}
	if len(cfg.Services) == 0 {
		return nil, fmt.Errorf("no services defined in config")
	}
	return &cfg, nil
}

func startServices(configFile string, force bool) error {
	cfg, err := parseConfig(configFile)
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	ctx := context.Background()
	for name, svc := range cfg.Services {
		if force {
			// Stop and remove if exists
			timeoutSeconds := 30
			if err := cli.ContainerStop(ctx, name, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
				// Ignore if not found
			}
			if err := cli.ContainerRemove(ctx, name, container.RemoveOptions{Force: true}); err != nil {
				// Ignore if not found
			}
		}

		// Pull image and show progress
		reader, err := cli.ImagePull(ctx, svc.Image, image.PullOptions{})
		if err != nil {
			return fmt.Errorf("pull image %s for %s: %w", svc.Image, name, err)
		}
		defer reader.Close()
		// Print pull progress to stdout
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return fmt.Errorf("pull image %s for %s failed during streaming: %w", svc.Image, name, err)
		}

		// Create container
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: svc.Image,
			Env:   svc.Env,
		}, &container.HostConfig{
			PortBindings: convertPorts(svc.Ports),
			AutoRemove:   false,
		}, nil, nil, name)
		if err != nil {
			return fmt.Errorf("create container %s: %w", name, err)
		}
		// Start container
		if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			return fmt.Errorf("start container %s: %w", name, err)
		}
		fmt.Printf("Started %s\n", name)
	}
	return nil
}

func stopServices(configFile string) error {
	cfg, err := parseConfig(configFile)
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	ctx := context.Background()
	for name := range cfg.Services {
		// Stop container
		timeoutSeconds := 30 // Timeout in seconds
		if err := cli.ContainerStop(ctx, name, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
			fmt.Printf("Warning: stop %s: %v (may not exist)\n", name, err)
		}
		// Remove container
		if err := cli.ContainerRemove(ctx, name, container.RemoveOptions{Force: true}); err != nil {
			fmt.Printf("Warning: remove %s: %v\n", name, err)
		} else {
			fmt.Printf("Stopped and removed %s\n", name)
		}
	}
	return nil
}

func convertPorts(ports map[string]string) nat.PortMap {
	portMap := make(nat.PortMap)
	for containerPort, hostPort := range ports {
		portMap[nat.Port(containerPort)] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: hostPort},
		}
	}
	return portMap
}
