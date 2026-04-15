package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"claw-destine.com/camboose/service/datatypes"
	"claw-destine.com/camboose/service/graphdb"
	"claw-destine.com/camboose/service/recipies"
)

type fakeProjectController struct {
	projects map[string]datatypes.Project
}

func newFakeProjectController() *fakeProjectController {
	return &fakeProjectController{projects: make(map[string]datatypes.Project)}
}

func (f *fakeProjectController) CreateProject(_ context.Context, project datatypes.Project) (datatypes.Project, error) {
	project.Name = strings.TrimSpace(project.Name)
	project.Recipe = strings.TrimSpace(project.Recipe)

	if project.Name == "" || project.Recipe == "" {
		return datatypes.Project{}, graphdb.ErrProjectInvalid
	}

	if _, exists := f.projects[project.Name]; exists {
		return datatypes.Project{}, graphdb.ErrProjectExists
	}

	f.projects[project.Name] = project
	return project, nil
}

func (f *fakeProjectController) ListProjects(_ context.Context) ([]datatypes.Project, error) {
	projects := make([]datatypes.Project, 0, len(f.projects))
	for _, project := range f.projects {
		projects = append(projects, project)
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	return projects, nil
}

func (f *fakeProjectController) GetProject(_ context.Context, name string) (datatypes.Project, bool, error) {
	project, ok := f.projects[name]
	if !ok {
		return datatypes.Project{}, false, nil
	}

	return project, true, nil
}

func (f *fakeProjectController) DeleteProject(_ context.Context, name string) (bool, error) {
	if _, ok := f.projects[name]; !ok {
		return false, nil
	}

	delete(f.projects, name)
	return true, nil
}

func (f *fakeProjectController) Close(_ context.Context) error {
	return nil
}

func newTestConfig(t *testing.T) RestConfig {
	t.Helper()

	recipeDir := t.TempDir()
	recipePath := filepath.Join(recipeDir, "example.yaml")
	recipeContent := "description: Example recipe\n"
	if err := os.WriteFile(recipePath, []byte(recipeContent), 0o600); err != nil {
		t.Fatalf("failed to write recipe file: %v", err)
	}

	return RestConfig{
		RecipiesCtr: recipies.NewRecipeController(recipies.RecipeConfig{RecipePath: recipeDir}),
		ProjectsCtr: newFakeProjectController(),
	}
}

func TestRecipiesHandlerReturnsMethodNotAllowedForNonGet(t *testing.T) {
	cfg := newTestConfig(t)
	h := newRecipiesHandler(cfg)

	req := httptest.NewRequest(http.MethodPost, "/api/recipies", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}

	if got := rr.Header().Get("Allow"); got != http.MethodGet {
		t.Fatalf("unexpected Allow header: got %q, want %q", got, http.MethodGet)
	}
}

func TestRecipiesHandlerRequiresBasicAuthWhenConfigured(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.BasicUser = "user"
	cfg.BasicPassword = "pass"
	h := newRecipiesHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/recipies", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusUnauthorized)
	}

	if got := rr.Header().Get("WWW-Authenticate"); got == "" {
		t.Fatal("expected WWW-Authenticate header to be set")
	}
}

func TestRecipiesHandlerReturnsUnauthorizedForWrongCredentials(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.BasicUser = "user"
	cfg.BasicPassword = "pass"
	h := newRecipiesHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/recipies", nil)
	req.SetBasicAuth("user", "wrong")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestRecipiesHandlerAllowsGetWithoutAuthWhenNotFullyConfigured(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.BasicUser = "user"
	h := newRecipiesHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/recipies", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestRecipiesHandlerReturnsRecipesForValidGetWithCredentials(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.BasicUser = "user"
	cfg.BasicPassword = "pass"
	h := newRecipiesHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/recipies", nil)
	req.SetBasicAuth("user", "pass")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusOK)
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected Content-Type: got %q, want %q", got, "application/json")
	}

	var gotRecipes []datatypes.Recipe
	if err := json.Unmarshal(rr.Body.Bytes(), &gotRecipes); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(gotRecipes) != 1 {
		t.Fatalf("unexpected recipes count: got %d, want 1", len(gotRecipes))
	}

	if gotRecipes[0].Name != "example" {
		t.Fatalf("unexpected recipe name: got %q, want %q", gotRecipes[0].Name, "example")
	}
}

func TestProjectsCollectionHandlerRejectsUnsupportedMethod(t *testing.T) {
	cfg := newTestConfig(t)
	h := newProjectsCollectionHandler(cfg)

	req := httptest.NewRequest(http.MethodDelete, "/api/project", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestProjectsCollectionHandlerCreatesProject(t *testing.T) {
	cfg := newTestConfig(t)
	h := newProjectsCollectionHandler(cfg)

	body := bytes.NewBufferString(`{"name":"demo","recipe":"static_html_with_blog"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/project", body)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusCreated)
	}

	var project datatypes.Project
	if err := json.Unmarshal(rr.Body.Bytes(), &project); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if project.Name != "demo" || project.Recipe != "static_html_with_blog" {
		t.Fatalf("unexpected project: got %+v", project)
	}
}

func TestProjectsCollectionHandlerReturnsConflictWhenProjectExists(t *testing.T) {
	cfg := newTestConfig(t)
	h := newProjectsCollectionHandler(cfg)

	body := bytes.NewBufferString(`{"name":"demo","recipe":"static_html_with_blog"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/project", body)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	req = httptest.NewRequest(http.MethodPost, "/api/project", bytes.NewBufferString(`{"name":"demo","recipe":"other"}`))
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusConflict)
	}
}

func TestProjectsCollectionHandlerListsProjects(t *testing.T) {
	cfg := newTestConfig(t)
	h := newProjectsCollectionHandler(cfg)

	seedProjects(t, h)

	req := httptest.NewRequest(http.MethodGet, "/api/project", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusOK)
	}

	var projects []datatypes.Project
	if err := json.Unmarshal(rr.Body.Bytes(), &projects); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	want := []datatypes.Project{
		{Name: "demo", Recipe: "static_html_with_blog"},
		{Name: "other", Recipe: "static_html_with_blog"},
	}
	if !reflect.DeepEqual(projects, want) {
		t.Fatalf("unexpected projects: got %+v, want %+v", projects, want)
	}
}

func TestProjectsItemHandlerReturnsNotFoundForUnknownProject(t *testing.T) {
	cfg := newTestConfig(t)
	h := newProjectsItemHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/project/missing", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestProjectsItemHandlerReturnsProject(t *testing.T) {
	cfg := newTestConfig(t)
	collection := newProjectsCollectionHandler(cfg)
	seedProject(t, collection, `{"name":"demo","recipe":"static_html_with_blog"}`)

	h := newProjectsItemHandler(cfg)
	req := httptest.NewRequest(http.MethodGet, "/api/project/demo", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusOK)
	}

	var project datatypes.Project
	if err := json.Unmarshal(rr.Body.Bytes(), &project); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if project.Name != "demo" {
		t.Fatalf("unexpected project name: got %q, want %q", project.Name, "demo")
	}
}

func TestProjectsItemHandlerDeletesProject(t *testing.T) {
	cfg := newTestConfig(t)
	collection := newProjectsCollectionHandler(cfg)
	seedProject(t, collection, `{"name":"demo","recipe":"static_html_with_blog"}`)

	h := newProjectsItemHandler(cfg)
	req := httptest.NewRequest(http.MethodDelete, "/api/project/demo", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusNoContent)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/project/demo", nil)
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("unexpected status after delete: got %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestProjectHandlersRequireBasicAuthWhenConfigured(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.BasicUser = "user"
	cfg.BasicPassword = "pass"
	h := newProjectsCollectionHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/project", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestNewMuxRoutesProjectEndpoints(t *testing.T) {
	cfg := newTestConfig(t)
	h := newMux(cfg)

	seedReq := httptest.NewRequest(http.MethodPost, "/api/project", bytes.NewBufferString(`{"name":"demo","recipe":"static_html_with_blog"}`))
	seedRes := httptest.NewRecorder()
	h.ServeHTTP(seedRes, seedReq)
	if seedRes.Code != http.StatusCreated {
		t.Fatalf("unexpected seed status: got %d, want %d", seedRes.Code, http.StatusCreated)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/project/demo", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusOK)
	}
}

func seedProjects(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	seedProject(t, handler, `{"name":"demo","recipe":"static_html_with_blog"}`)
	seedProject(t, handler, `{"name":"other","recipe":"static_html_with_blog"}`)
}

func seedProject(t *testing.T, handler http.HandlerFunc, payload string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/project", bytes.NewBufferString(payload))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("unexpected seed status: got %d, want %d", rr.Code, http.StatusCreated)
	}
}

func TestProjectsCollectionHandlerReturnsBadRequestForInvalidJSON(t *testing.T) {
	cfg := newTestConfig(t)
	h := newProjectsCollectionHandler(cfg)

	req := httptest.NewRequest(http.MethodPost, "/api/project", bytes.NewBufferString(`{"name":`))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestProjectsCollectionHandlerReturnsBadRequestForInvalidPayload(t *testing.T) {
	cfg := newTestConfig(t)
	h := newProjectsCollectionHandler(cfg)

	req := httptest.NewRequest(http.MethodPost, "/api/project", bytes.NewBufferString(`{"name":"demo"}`))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestProjectsCollectionHandlerReturnsInternalServerErrorWhenControllerFails(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.ProjectsCtr = errProjectController{}
	h := newProjectsCollectionHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/project", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got %d, want %d", rr.Code, http.StatusInternalServerError)
	}
}

type errProjectController struct{}

func (errProjectController) CreateProject(context.Context, datatypes.Project) (datatypes.Project, error) {
	return datatypes.Project{}, errors.New("boom")
}

func (errProjectController) ListProjects(context.Context) ([]datatypes.Project, error) {
	return nil, errors.New("boom")
}

func (errProjectController) GetProject(context.Context, string) (datatypes.Project, bool, error) {
	return datatypes.Project{}, false, errors.New("boom")
}

func (errProjectController) DeleteProject(context.Context, string) (bool, error) {
	return false, errors.New("boom")
}

func (errProjectController) Close(context.Context) error {
	return nil
}
