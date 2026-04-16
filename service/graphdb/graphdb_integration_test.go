//go:build integration

package graphdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"claw-destine.com/camboose/service/datatypes"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNeo4jProjectControllerIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "neo4j:5.26.1",
			ExposedPorts: []string{"7687/tcp"},
			Env: map[string]string{
				"NEO4J_AUTH": "neo4j/password",
			},
			WaitingFor: wait.ForListeningPort("7687/tcp").WithStartupTimeout(90 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start neo4j container: %v", err)
	}
	defer func() {
		_ = container.Terminate(context.Background())
	}()

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "7687/tcp")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	cfg := GraphDBConfig{
		URI:      fmt.Sprintf("bolt://%s:%s", host, port.Port()),
		Username: "neo4j",
		Password: "password",
		Database: "neo4j",
	}

	controller, err := newControllerWithRetry(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to create project controller: %v", err)
	}
	defer func() {
		_ = controller.Close(context.Background())
	}()

	created, err := controller.CreateProject(ctx, datatypes.Project{Name: "alpha", Recipe: "static_html_with_blog"})
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	if created.Name != "alpha" {
		t.Fatalf("unexpected created project name: got %q, want %q", created.Name, "alpha")
	}

	if _, err := controller.CreateProject(ctx, datatypes.Project{Name: "alpha", Recipe: "other"}); err != ErrProjectExists {
		t.Fatalf("expected ErrProjectExists, got %v", err)
	}

	listed, err := controller.ListProjects(ctx)
	if err != nil {
		t.Fatalf("failed to list projects: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("unexpected projects count: got %d, want 1", len(listed))
	}

	project, found, err := controller.GetProject(ctx, "alpha")
	if err != nil {
		t.Fatalf("failed to get project: %v", err)
	}
	if !found || project.Name != "alpha" {
		t.Fatalf("unexpected get project result: found=%t, project=%+v", found, project)
	}

	version1, err := controller.CreateVersion(ctx, "alpha", datatypes.Version{Number: 1, Name: "MVP", Status: "active"})
	if err != nil {
		t.Fatalf("failed to create version: %v", err)
	}
	if version1.Number != 1 || version1.Name != "MVP" || version1.Status != "active" {
		t.Fatalf("unexpected created version: %+v", version1)
	}

	if _, err := controller.CreateVersion(ctx, "alpha", datatypes.Version{Number: 1, Status: "released"}); err != ErrVersionExists {
		t.Fatalf("expected ErrVersionExists, got %v", err)
	}

	if _, err := controller.CreateVersion(ctx, "missing", datatypes.Version{Number: 1, Status: "active"}); err != ErrProjectNotFound {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}

	versions, err := controller.ListVersions(ctx, "alpha")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("unexpected versions count: got %d, want 1", len(versions))
	}

	version, found, err := controller.GetVersion(ctx, "alpha", 1)
	if err != nil {
		t.Fatalf("failed to get version: %v", err)
	}
	if !found || version.Number != 1 {
		t.Fatalf("unexpected get version result: found=%t, version=%+v", found, version)
	}

	deleted, err := controller.DeleteProject(ctx, "alpha")
	if err != nil {
		t.Fatalf("failed to delete project: %v", err)
	}
	if !deleted {
		t.Fatal("expected project to be deleted")
	}

	_, found, err = controller.GetProject(ctx, "alpha")
	if err != nil {
		t.Fatalf("failed to get project after delete: %v", err)
	}
	if found {
		t.Fatal("expected deleted project to be missing")
	}

	if _, _, err := controller.GetVersion(ctx, "alpha", 1); err != ErrProjectNotFound {
		t.Fatalf("expected ErrProjectNotFound after project delete, got %v", err)
	}
}

func newControllerWithRetry(ctx context.Context, cfg GraphDBConfig) (*Neo4jProjectController, error) {
	deadline := time.Now().Add(45 * time.Second)
	var lastErr error

	for time.Now().Before(deadline) {
		controller, err := NewNeo4jProjectController(ctx, cfg)
		if err == nil {
			return controller, nil
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("neo4j did not become ready before deadline")
	}

	return nil, lastErr
}
