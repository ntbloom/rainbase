package database

import (
	"github.com/sirupsen/logrus"
	"os/exec"
)

/* Simple wrapper around podman for setting up test fixtures */
const container = "podman"

type containerConfig struct {
	name  string
	image string
	port  int
}

var podman = containerConfig{ //nolint:gochecknoglobals
	name: "postgresql",
	//image: "docker.io/library.postgres:latest",
	image: "postgres",
	port:  5432,
}

// PostgresUp fixture for spinning up postgresql podman container
func PostgresContainerUp() error {
	_, err := exec.Command(container, //nolint:gosec
		"run",
		"--rm",
		"--name",
		podman.name,
		"-e",
		"POSTGRES_PASSWORD=password",
		podman.image).Output()
	if err != nil {
		return err
	}
	return nil
}

func PostgresContainerDown() {
	logrus.Info("shutting down podman container gracefully")
	_, err := exec.Command(container, "kill", podman.name).Output() //nolint:gosec
	if err != nil {
		logrus.Error(err)
	}
}
