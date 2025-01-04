package main

import (
	"archive/tar"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
)

//go:embed dockerfile/gptr.Dockerfile
var gptrDockerfile []byte

//go:embed dockerfile/gptr.env
var gptrEnvFile string

const gptrImageName = "openaide/gptr"
const gptrContainerName = "gptr-v3.1.7"

func parseEnvFile(envData string) []string {
	lines := strings.Split(envData, "\n")
	var envVars []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		envVars = append(envVars, line)
	}
	return envVars
}

type DockerBuildResponse struct {
	Stream   string `json:"stream,omitempty"`
	Error    string `json:"error,omitempty"`
	ErrorMsg string `json:"message,omitempty"`
}

// buildDockerImage constructs and builds a Docker image using provided contents and parameters.
func buildDockerImage(ctx context.Context, dockerfileName, tag string, dockerfileContent []byte) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	tarBuffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(tarBuffer)

	// Add Dockerfile to the tar buffer
	if err := addFileToTar(tarWriter, dockerfileName, dockerfileContent); err != nil {
		return err
	}

	if err := tarWriter.Close(); err != nil {
		return err
	}

	// Prepare to build the image
	tarReader := bytes.NewReader(tarBuffer.Bytes())
	buildOptions := types.ImageBuildOptions{
		Context:        tarReader,
		Dockerfile:     dockerfileName,
		Tags:           []string{tag},
		Remove:         true,
		SuppressOutput: false,
		PullParent:     true,
	}

	resp, err := cli.ImageBuild(ctx, tarReader, buildOptions)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Use jsonmessage to handle the response stream
	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stdout, os.Stdout.Fd(), true, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Build image with tag %s succeeded\n", tag)

	return nil
}

// Utility function to add a file to the tarball
func addFileToTar(tw *tar.Writer, name string, fileContent []byte) error {
	header := &tar.Header{
		Name: name,
		Mode: 0600,
		Size: int64(len(fileContent)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tw.Write(fileContent); err != nil {
		return err
	}
	return nil
}

// BuildGPTRImage builds the GPT Researcher Docker image
func BuildGPTRImage(ctx context.Context) error {
	return buildDockerImage(ctx, "gptr.Dockerfile", gptrImageName, gptrDockerfile)
}

func RunGPTRContainer(ctx context.Context, query string, outDir string) error {
	envVars := parseEnvFile(gptrEnvFile)
	output, err := filepath.Abs(outDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	args := []string{query, "--report_type", "research_report"}
	config := &container.Config{
		Image: gptrImageName,
		Env:   envVars,
		Cmd:   args,
	}

	fmt.Printf("config: %+v\n", config)

	hostConfig := &container.HostConfig{
		Binds: []string{output + ":/app/outputs/"},
	}

	fmt.Printf("hostConfig: %+v\n", hostConfig)

	_, err = runContainer(ctx, gptrContainerName, config, hostConfig)
	if err != nil {
		// Log the error
		fmt.Printf("Error running container: %v\n", err)

		// Attempt to remove the container
		if rmErr := removeContainer(ctx, gptrContainerName); rmErr != nil {
			fmt.Printf("Error removing container: %v\n", rmErr)
		}
		return err
	}

	return nil
}

func runContainer(ctx context.Context, containerName string, config *container.Config, hostConfig *container.HostConfig) (*container.CreateResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	// try remove and create again if container already exists
	if errdefs.IsConflict(err) {
		removeContainer(ctx, containerName)

		resp, err = cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	attachOptions := container.AttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	}

	hjResp, err := cli.ContainerAttach(ctx, resp.ID, attachOptions)
	if err != nil {
		return &resp, err
	}
	defer hjResp.Close()

	// progress output
	go func() {
		if _, err := stdcopy.StdCopy(os.Stderr, os.Stderr, hjResp.Reader); err != nil {
			fmt.Fprintf(os.Stderr, "error copying output: %v\n", err)
		}
	}()

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return &resp, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return &resp, err
		}
	case <-statusCh:
	}

	return &resp, nil
}

func removeContainer(ctx context.Context, containerName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// Get the container ID from the container name
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err
	}

	var containerID string
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+containerName {
				containerID = container.ID
				break
			}
		}
		if containerID != "" {
			break
		}
	}

	if containerID == "" {
		fmt.Printf("Container %s not found\n", containerName)
		return nil
	}

	// Remove the container
	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		return err
	}

	fmt.Printf("Container %s removed successfully\n", containerName)
	return nil
}
