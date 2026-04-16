package graphdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"claw-destine.com/camboose/service/datatypes"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
)

var (
	ErrProjectExists   = errors.New("project already exists")
	ErrProjectInvalid  = errors.New("project payload is invalid")
	ErrProjectNotFound = errors.New("project not found")
	ErrVersionExists   = errors.New("version already exists")
	ErrVersionInvalid  = errors.New("version payload is invalid")
)

type GraphDBConfig struct {
	URI      string `env:"URI" envDefault:"bolt://localhost:7687"`
	Username string `env:"USERNAME" envDefault:"neo4j"`
	Password string `env:"PASSWORD" envDefault:"password"`
	Database string `env:"DATABASE" envDefault:"neo4j"`
}

type ProjectController interface {
	CreateProject(ctx context.Context, project datatypes.Project) (datatypes.Project, error)
	ListProjects(ctx context.Context) ([]datatypes.Project, error)
	GetProject(ctx context.Context, name string) (datatypes.Project, bool, error)
	DeleteProject(ctx context.Context, name string) (bool, error)
	CreateVersion(ctx context.Context, projectName string, version datatypes.Version) (datatypes.Version, error)
	ListVersions(ctx context.Context, projectName string) ([]datatypes.Version, error)
	GetVersion(ctx context.Context, projectName string, number int) (datatypes.Version, bool, error)
	Close(ctx context.Context) error
}

type Neo4jProjectController struct {
	driver   neo4j.DriverWithContext
	database string
}

func NewNeo4jProjectController(ctx context.Context, cfg GraphDBConfig) (*Neo4jProjectController, error) {
	driver, err := neo4j.NewDriverWithContext(cfg.URI, neo4j.BasicAuth(cfg.Username, cfg.Password, ""))
	if err != nil {
		return nil, fmt.Errorf("create neo4j driver: %w", err)
	}

	if err := driver.VerifyConnectivity(ctx); err != nil {
		_ = driver.Close(ctx)
		return nil, fmt.Errorf("verify neo4j connectivity: %w", err)
	}

	controller := &Neo4jProjectController{driver: driver, database: cfg.Database}
	if err := controller.ensureConstraints(ctx); err != nil {
		_ = driver.Close(ctx)
		return nil, err
	}

	return controller, nil
}

func (c *Neo4jProjectController) ensureConstraints(ctx context.Context) error {
	session := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.database})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(
			ctx,
			`CREATE CONSTRAINT project_name_unique IF NOT EXISTS
			 FOR (p:Project)
			 REQUIRE p.name IS UNIQUE`,
			nil,
		)
		if err != nil {
			return nil, err
		}

		_, err = tx.Run(
			ctx,
			`CREATE CONSTRAINT version_project_number_unique IF NOT EXISTS
			 FOR (v:Version)
			 REQUIRE (v.projectName, v.number) IS UNIQUE`,
			nil,
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("ensure project constraints: %w", err)
	}

	return nil
}

func (c *Neo4jProjectController) CreateProject(ctx context.Context, project datatypes.Project) (datatypes.Project, error) {
	project.Name = strings.TrimSpace(project.Name)
	project.Recipe = strings.TrimSpace(project.Recipe)

	if project.Name == "" || project.Recipe == "" {
		return datatypes.Project{}, ErrProjectInvalid
	}

	session := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.database})
	defer session.Close(ctx)

	_, exists, err := c.GetProject(ctx, project.Name)
	if err != nil {
		return datatypes.Project{}, err
	}
	if exists {
		return datatypes.Project{}, ErrProjectExists
	}

	created, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			`CREATE (p:Project {name: $name, recipe: $recipe})
			 RETURN p.name AS name, p.recipe AS recipe`,
			map[string]any{"name": project.Name, "recipe": project.Recipe},
		)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return datatypes.Project{
			Name:   getStringValue(record, "name"),
			Recipe: getStringValue(record, "recipe"),
		}, nil
	})
	if err != nil {
		return datatypes.Project{}, fmt.Errorf("create project: %w", err)
	}

	return created.(datatypes.Project), nil
}

func (c *Neo4jProjectController) ListProjects(ctx context.Context) ([]datatypes.Project, error) {
	session := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.database})
	defer session.Close(ctx)

	projectsAny, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			`MATCH (p:Project)
			 RETURN p.name AS name, p.recipe AS recipe
			 ORDER BY p.name`,
			nil,
		)
		if err != nil {
			return nil, err
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		projects := make([]datatypes.Project, 0, len(records))
		for _, record := range records {
			projects = append(projects, datatypes.Project{
				Name:   getStringValue(record, "name"),
				Recipe: getStringValue(record, "recipe"),
			})
		}

		return projects, nil
	})
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	return projectsAny.([]datatypes.Project), nil
}

func (c *Neo4jProjectController) GetProject(ctx context.Context, name string) (datatypes.Project, bool, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return datatypes.Project{}, false, ErrProjectInvalid
	}

	session := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.database})
	defer session.Close(ctx)

	projectAny, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			`MATCH (p:Project {name: $name})
			 RETURN p.name AS name, p.recipe AS recipe`,
			map[string]any{"name": name},
		)
		if err != nil {
			return nil, err
		}

		if !result.Next(ctx) {
			if err := result.Err(); err != nil {
				return nil, err
			}

			return nil, nil
		}

		record := result.Record()

		return datatypes.Project{
			Name:   getStringValue(record, "name"),
			Recipe: getStringValue(record, "recipe"),
		}, nil
	})
	if err != nil {
		return datatypes.Project{}, false, fmt.Errorf("get project: %w", err)
	}

	if projectAny == nil {
		return datatypes.Project{}, false, nil
	}

	return projectAny.(datatypes.Project), true, nil
}

func (c *Neo4jProjectController) DeleteProject(ctx context.Context, name string) (bool, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return false, ErrProjectInvalid
	}

	session := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.database})
	defer session.Close(ctx)

	deletedAny, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			`MATCH (p:Project {name: $name})
			 OPTIONAL MATCH (p)-[:HAS_VERSION]->(v:Version)
			 WITH p, COLLECT(v) AS versions
			 FOREACH (version IN versions | DETACH DELETE version)
			 WITH p
			 DELETE p
			 RETURN COUNT(p) AS deleted`,
			map[string]any{"name": name},
		)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		deleted, _ := record.Get("deleted")
		deletedCount, ok := deleted.(int64)
		if !ok {
			return false, nil
		}

		return deletedCount > 0, nil
	})
	if err != nil {
		return false, fmt.Errorf("delete project: %w", err)
	}

	return deletedAny.(bool), nil
}

func (c *Neo4jProjectController) CreateVersion(ctx context.Context, projectName string, version datatypes.Version) (datatypes.Version, error) {
	projectName = strings.TrimSpace(projectName)
	version.Name = strings.TrimSpace(version.Name)
	version.Status = strings.TrimSpace(version.Status)

	if projectName == "" || version.Number <= 0 || version.Status == "" {
		return datatypes.Version{}, ErrVersionInvalid
	}

	session := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.database})
	defer session.Close(ctx)

	project, found, err := c.GetProject(ctx, projectName)
	if err != nil {
		return datatypes.Version{}, err
	}
	if !found || project.Name == "" {
		return datatypes.Version{}, ErrProjectNotFound
	}

	_, exists, err := c.GetVersion(ctx, projectName, version.Number)
	if err != nil {
		return datatypes.Version{}, err
	}
	if exists {
		return datatypes.Version{}, ErrVersionExists
	}

	createdAny, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			`MATCH (p:Project {name: $projectName})
			 CREATE (p)-[:HAS_VERSION]->(v:Version {
				projectName: $projectName,
				number: $number,
				name: $name,
				status: $status
			 })
			 RETURN v.number AS number, v.name AS name, v.status AS status`,
			map[string]any{
				"projectName": projectName,
				"number":      int64(version.Number),
				"name":        version.Name,
				"status":      version.Status,
			},
		)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return datatypes.Version{
			Number: getIntValue(record, "number"),
			Name:   getStringValue(record, "name"),
			Status: getStringValue(record, "status"),
		}, nil
	})
	if err != nil {
		return datatypes.Version{}, fmt.Errorf("create version: %w", err)
	}

	return createdAny.(datatypes.Version), nil
}

func (c *Neo4jProjectController) ListVersions(ctx context.Context, projectName string) ([]datatypes.Version, error) {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return nil, ErrVersionInvalid
	}

	project, found, err := c.GetProject(ctx, projectName)
	if err != nil {
		return nil, err
	}
	if !found || project.Name == "" {
		return nil, ErrProjectNotFound
	}

	session := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.database})
	defer session.Close(ctx)

	versionsAny, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			`MATCH (:Project {name: $projectName})-[:HAS_VERSION]->(v:Version)
			 RETURN v.number AS number, v.name AS name, v.status AS status
			 ORDER BY v.number`,
			map[string]any{"projectName": projectName},
		)
		if err != nil {
			return nil, err
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		versions := make([]datatypes.Version, 0, len(records))
		for _, record := range records {
			versions = append(versions, datatypes.Version{
				Number: getIntValue(record, "number"),
				Name:   getStringValue(record, "name"),
				Status: getStringValue(record, "status"),
			})
		}

		return versions, nil
	})
	if err != nil {
		return nil, fmt.Errorf("list versions: %w", err)
	}

	return versionsAny.([]datatypes.Version), nil
}

func (c *Neo4jProjectController) GetVersion(ctx context.Context, projectName string, number int) (datatypes.Version, bool, error) {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" || number <= 0 {
		return datatypes.Version{}, false, ErrVersionInvalid
	}

	project, found, err := c.GetProject(ctx, projectName)
	if err != nil {
		return datatypes.Version{}, false, err
	}
	if !found || project.Name == "" {
		return datatypes.Version{}, false, ErrProjectNotFound
	}

	session := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.database})
	defer session.Close(ctx)

	versionAny, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			`MATCH (:Project {name: $projectName})-[:HAS_VERSION]->(v:Version {projectName: $projectName, number: $number})
			 RETURN v.number AS number, v.name AS name, v.status AS status`,
			map[string]any{"projectName": projectName, "number": int64(number)},
		)
		if err != nil {
			return nil, err
		}

		if !result.Next(ctx) {
			if err := result.Err(); err != nil {
				return nil, err
			}

			return nil, nil
		}

		record := result.Record()

		return datatypes.Version{
			Number: getIntValue(record, "number"),
			Name:   getStringValue(record, "name"),
			Status: getStringValue(record, "status"),
		}, nil
	})
	if err != nil {
		return datatypes.Version{}, false, fmt.Errorf("get version: %w", err)
	}

	if versionAny == nil {
		return datatypes.Version{}, false, nil
	}

	return versionAny.(datatypes.Version), true, nil
}

func (c *Neo4jProjectController) Close(ctx context.Context) error {
	return c.driver.Close(ctx)
}

func getStringValue(record *neo4j.Record, key string) string {
	value, ok := record.Get(key)
	if !ok {
		return ""
	}

	stringValue, ok := value.(string)
	if !ok {
		return ""
	}

	return stringValue
}

func getIntValue(record *neo4j.Record, key string) int {
	value, ok := record.Get(key)
	if !ok {
		return 0
	}

	intValue, ok := value.(int64)
	if !ok {
		return 0
	}

	return int(intValue)
}
