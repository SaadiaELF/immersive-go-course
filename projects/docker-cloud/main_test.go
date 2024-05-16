package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
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
	resource, err := pool.BuildAndRun("docker-cloud-test", "./Dockerfile", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	fmt.Println(resource.GetPort("80/tcp"))
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
