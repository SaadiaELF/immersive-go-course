package e2e_test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()

	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Build and run the given Dockerfile
	resource, err := pool.BuildAndRunWithBuildOptions(&dockertest.BuildOptions{
		ContextDir: "../",
		Dockerfile: "Dockerfile",
	}, &dockertest.RunOptions{
		Name:         "docker-cloud-e2e",
		ExposedPorts: []string{"80/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"80/tcp": {{HostPort: "80"}},
		}})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		fmt.Println("Checking API connection...")
		_, err := http.Get("http://localhost:80/")
		if err != nil {
			log.Printf("Could not connect to server: %s", err)
		}
		return err
	}); err != nil {
		fmt.Printf("could not start resource: %v", err)
	}

	code := m.Run()

	// When we are done, kill and remove the container
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
