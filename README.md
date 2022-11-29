# Multiplatform Action

This is an example of a single dagger pipeline distributed as both a Github Action and a CircleCI orb. The example action takes an input of a python script to run, and the action runs it.

## How does it work?

The dagger pipeline to be executed by either GitHub Actions or CircleCI is in [main.go](./main.go). This pipeline is compiled and distributed as a GitHub release when a tag is pushed.

The github action entrypoint is defined by [action.yml](./action.yml). It reads the input and runs [src/scripts/action.sh](./src/scripts/action.sh).

The CircleCI orb entrypoint is defined by [src/@orb.yml](./src/%40orb.yml). It reads the input and runs [src/scripts/action.sh](./src/scripts/action.sh).

The [src/scripts/action.sh](./src/scripts/action.sh) script will download the latest released compiled pipeline from this repo's github releases for the appropriate architecture of the CI runner and execute the pipeline binary. The pipeilne binary will run the Dagger pipeline by creating a python container, mounting the input script, and executing it.
