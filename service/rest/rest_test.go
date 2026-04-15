package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"claw-destine.com/camboose/service/recipies"
)

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

	var gotRecipes []recipies.Recipe
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
