package projects

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

func setupProjectsTestDB(t *testing.T) (*gorm.DB, func()) {
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
	if err := db.AutoMigrate(&dt.Project{}, &dt.Version{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	cleanup := func() {
		_ = pgContainer.Terminate(ctx)
	}
	return db, cleanup
}

func TestListProjects_IncludesVersionCountsGroupedByStatus(t *testing.T) {
	db, cleanup := setupProjectsTestDB(t)
	defer cleanup()

	pc := &ProjectControler{Db: db}

	project1 := dt.Project{Base: dt.Base{Name: "Project 1"}}
	project2 := dt.Project{Base: dt.Base{Name: "Project 2"}}
	project3 := dt.Project{Base: dt.Base{Name: "Project 3"}}
	assert.NoError(t, db.Create(&project1).Error)
	assert.NoError(t, db.Create(&project2).Error)
	assert.NoError(t, db.Create(&project3).Error)

	versions := []dt.Version{
		{Base: dt.Base{Name: "v1"}, ProjectId: project1.Id, Status: dt.RS_New},
		{Base: dt.Base{Name: "v2"}, ProjectId: project1.Id, Status: dt.RS_New},
		{Base: dt.Base{Name: "v3"}, ProjectId: project1.Id, Status: dt.RS_Done},
		{Base: dt.Base{Name: "v4"}, ProjectId: project2.Id, Status: dt.RS_InReview},
	}
	for _, v := range versions {
		assert.NoError(t, db.Create(&v).Error)
	}

	projects, err := pc.ListProjects(nil)
	assert.NoError(t, err)
	assert.Len(t, projects, 3)

	byID := map[string]dt.Project{}
	for _, p := range projects {
		byID[p.Id] = p
	}

	assert.Equal(t, 2, byID[project1.Id].VersionStatusCounts[dt.RS_New])
	assert.Equal(t, 1, byID[project1.Id].VersionStatusCounts[dt.RS_Done])
	assert.Equal(t, 1, byID[project2.Id].VersionStatusCounts[dt.RS_InReview])
	assert.Empty(t, byID[project3.Id].VersionStatusCounts)
}
