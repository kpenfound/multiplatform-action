package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	err := build(ctx)
	if err != nil {
		panic(err)
	}
}

func build(ctx context.Context) error {
	buildImage := "golang:1.19"
	// define build matrix
	oses := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}


	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// get project
	src := client.Host().Directory(".",
		dagger.HostDirectoryOpts{Include: []string{"main.go", "go.mod", "go.sum"}},
  )

	outputDirectory := client.Directory()

	golang := client.Container().
		From(buildImage).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")

	for _, goos := range oses {
		for _, goarch := range arches {
			// create a directory for each os, arch and version
			bin := fmt.Sprintf("build/action_%s_%s", goos, goarch)

			// set GOARCH and GOOS in the build environment
			build := golang.WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-o", bin})

			outputDirectory = outputDirectory.WithFile(bin, build.File(bin))
		}
	}
	_, err = outputDirectory.Export(ctx, ".")
	return err
}
