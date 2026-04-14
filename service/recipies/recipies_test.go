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
		gotNames = append(gotNames, recipe)
	}

	want := []string{"alpha_recipe", "notes", "zeta_recipe"}
	if !reflect.DeepEqual(gotNames, want) {
		t.Fatalf("unexpected recipe names: got %v, want %v", gotNames, want)
	}
}
