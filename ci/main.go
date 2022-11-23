package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

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
	if err != nil {
		return err
	}

	if os.Getenv("RELEASE") == "true" {
		tag := os.Getenv("GITHUB_REF_NAME")
		buildDir := client.Host().Directory("build")
		ghcli := fmt.Sprintf("https://github.com/cli/cli/releases/download/v2.20.0/gh_2.20.0_linux_%s.tar.gz", runtime.GOARCH)
		ghcliPath := fmt.Sprintf("gh_2.20.0_linux_%s/bin/gh", runtime.GOARCH)

		alpine := client.Container().
			From("alpine").
			WithExec([]string{"apk", "add", "curl"}). // Download github cli
			WithExec([]string{"curl", "-L", "-o", "ghcli.tar.gz", ghcli}).
			WithExec([]string{"tar", "-xvf", "ghcli.tar.gz"}).
			WithMountedDirectory("/build", buildDir).
			WithEnvVariable("GH_TOKEN", os.Getenv("GH_ELEVATED_TOKEN")). // create github release
			WithExec([]string{ghcliPath, "release", "create", tag, "/build/*"})
		_, err := alpine.ExitCode(ctx)
		return err
	}
	
	return nil
}
