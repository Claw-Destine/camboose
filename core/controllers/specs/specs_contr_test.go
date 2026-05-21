package specs

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	dt "claw-destine.com/camboose/core/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()
	_, filename, _, _ := runtime.Caller(0)
	initScriptPath := filepath.Join(filepath.Dir(filename), "../../../init.sql")
	pgContainer, err := tcpostgres.Run(ctx,
		"pgvector/pgvector:0.8.2-pg18-trixie",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		tcpostgres.WithInitScripts(initScriptPath),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}
	dsn := "host=" + host + " port=" + port.Port() + " user=testuser password=testpass dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	if err := db.AutoMigrate(&dt.Project{}, &dt.SpecItem{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	cleanup := func() {
		_ = pgContainer.Terminate(ctx)
	}
	return db, cleanup
}

func TestSpecsController_CRUD(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	sc := &SpecsController{Db: db}
	project := dt.Project{Base: dt.Base{Name: "Test Project"}, Recipe: "test-recipe"}
	assert.NoError(t, db.Create(&project).Error)

	si := dt.SpecItem{Base: dt.Base{Name: "Spec1"}, ProjectId: project.Id, Type: dt.Version}
	created, err := sc.CreateSpecItem(si)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Spec1", created.Name)

	fetched, err := sc.GetSpecItemById(created.Id)
	assert.NoError(t, err)
	assert.Equal(t, created.Id, fetched.Id)

	created.Name = "Spec1 Updated"
	updated, err := sc.UpdateSpecItem(*created)
	assert.NoError(t, err)
	assert.Equal(t, "Spec1 Updated", updated.Name)

	list, err := sc.ListSpecItems(project.Id)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	assert.NoError(t, sc.DeleteSpecItemById(created.Id))
	list, err = sc.ListSpecItems(project.Id)
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}
