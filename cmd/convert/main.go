package main

import (
	"context"
	"fmt"
	"os"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
	composek8s "github.com/slaskis/compose-k8s"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	workingDir, _ := os.Getwd()

	project, err := NewProject(Options{
		Files:      []string{"compose.yaml"},
		Profiles:   []string{"*"},
		WorkingDir: workingDir,
	})
	if err != nil {
		return fmt.Errorf("unable to load files: %w", err)
	}

	resources, err := composek8s.Convert(project)
	if err != nil {
		return fmt.Errorf("unable to convert: %w", err)
	}

	err = resources.Write(os.Stdout)
	if err != nil {
		return fmt.Errorf("unable to write resources to file: %w", err)
	}

	return nil
}

type Options struct {
	Files      []string
	Profiles   []string
	WorkingDir string
}

func NewProject(opts Options) (*types.Project, error) {
	projectOptions, err := cli.NewProjectOptions(
		opts.Files,
		cli.WithOsEnv,
		cli.WithWorkingDirectory(opts.WorkingDir),
		cli.WithInterpolation(true),
		cli.WithProfiles(opts.Profiles),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create compose options: %w", err)
	}

	project, err := cli.ProjectFromOptions(context.Background(), projectOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to load files: %w", err)
	}

	return project, nil
}
