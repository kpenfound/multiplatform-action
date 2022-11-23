package main

import (
	"dagger.io/dagger"
	"context"
	"os"
	"fmt"
)

const pythonimg = "python:3.9-slim"

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		fmt.Printf("Error connecting to dagger engine: %+v", err)
		os.Exit(1)
	}
	defer client.Close()

	script := os.Getenv("INPUT_SCRIPT")
	if script == "" {
		fmt.Println("Please set input 'script'")
		os.Exit(1)
	}
	scriptDir := client.Host().Directory(".",
		dagger.HostDirectoryOpts{ Include: []string{script}},
  )

	py := client.Container().
	From(pythonimg).
	WithMountedDirectory("/src", scriptDir).
	WithWorkdir("/src").
	WithExec([]string{"python", script})

	out, err := py.Stdout(ctx)
	if err != nil {
		fmt.Printf("Error running pipeline: %+v\n", err)
		os.Exit(1)
	}
	fmt.Println(out)
}