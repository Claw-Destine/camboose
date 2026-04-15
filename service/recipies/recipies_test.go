package recipies

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListRecipiesReturnsFileNamesWithoutSuffix(t *testing.T) {
	recipeDir := t.TempDir()

	files := []string{
		"zeta_recipe.yaml",
		"alpha_recipe.yml",
		"notes.md",
	}

	for _, fileName := range files {
		filePath := filepath.Join(recipeDir, fileName)
		if err := os.WriteFile(filePath, []byte("content"), 0o600); err != nil {
			t.Fatalf("failed to create test file %q: %v", fileName, err)
		}
	}

	if err := os.Mkdir(filepath.Join(recipeDir, "nested"), 0o700); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	controller := NewRecipeController(RecipeConfig{RecipePath: recipeDir})

	got := controller.ListRecipies()

	gotNames := make([]string, 0, len(got))
	for _, recipe := range got {
		gotNames = append(gotNames, recipe.Name)
	}

	want := []string{"alpha_recipe", "notes", "zeta_recipe"}
	if !reflect.DeepEqual(gotNames, want) {
		t.Fatalf("unexpected recipe names: got %v, want %v", gotNames, want)
	}
}

func TestListRecipiesParsesYAMLFields(t *testing.T) {
	recipeDir := t.TempDir()

	content := `description: Static HTML with blog
project_management:
  relations:
    epic:
      view:
`

	filePath := filepath.Join(recipeDir, "static_html_with_blog.yaml")
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to create yaml recipe file: %v", err)
	}

	controller := NewRecipeController(RecipeConfig{RecipePath: recipeDir})
	got := controller.ListRecipies()

	if len(got) != 1 {
		t.Fatalf("unexpected recipe count: got %d, want 1", len(got))
	}

	recipe := got[0]
	if recipe.Name != "static_html_with_blog" {
		t.Fatalf("unexpected recipe name: got %q, want %q", recipe.Name, "static_html_with_blog")
	}

	if recipe.Description != "Static HTML with blog" {
		t.Fatalf("unexpected description: got %q, want %q", recipe.Description, "Static HTML with blog")
	}
}
