# LocalDev-CLI

A command-line tool developed in Go to automate the management of local development environments using Docker. It parses YAML configuration files to define services (e.g., databases like PostgreSQL or caches like Redis) and interacts with the Docker Engine API to spin up, tear down, or force-recreate complex setups with simple commands. This streamlines developer workflows by enabling quick environment management without manual Docker commands.

## Features
- **YAML Configuration**: Define services with images, ports, and environment variables in a simple YAML file.
- **Docker API Integration**: Pulls images, creates/starts/stops/removes containers programmatically.
- **Force Recreation**: Optional `--force` flag for `up` to automatically handle existing container conflicts.
- **Progress Feedback**: Displays real-time image pull progress.
- **Error Handling**: Graceful warnings for non-existent containers, validation for empty configs.

## Tech Stack
- Go (core language for CLI and API calls)
- Docker Engine API (via github.com/docker/docker/client)
- YAML parsing (gopkg.in/yaml.v2)

## Installation
1. Ensure Go (1.20+) and Docker Desktop are installed and running.
2. Clone the repo: `git clone https://github.com/martinchernyavskiy/localdev-cli.git`
3. Navigate to the directory: `cd localdev-cli`
4. Build the binary: `go build -o localdev`

## Usage
Create a `config.yaml` file (example included):
```yaml
services:
  db:
    image: postgres:13
    ports:
      5432/tcp: 5432
    env:
      - POSTGRES_PASSWORD=mysecret
  cache:
    image: redis:6
    ports:
      6379/tcp: 6379