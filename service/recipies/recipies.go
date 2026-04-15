package recipies

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"claw-destine.com/camboose/service/datatypes"
	"gopkg.in/yaml.v3"
)

type RecipeConfig struct {
	RecipePath string `env:"PATH"`
}

type RecipeController struct {
	cfg RecipeConfig
}

type recipeDocument struct {
	Description       string `yaml:"description"`
	ProjectManagement struct {
		Relations map[string]map[string]any `yaml:"relations"`
	} `yaml:"project_management"`
}

func NewRecipeController(cfg RecipeConfig) RecipeController {
	return RecipeController{cfg: cfg}
}

func (rc RecipeController) ListRecipies() []datatypes.Recipe {
	entries, err := os.ReadDir(rc.cfg.RecipePath)
	if err != nil {
		slog.Error("Could not read recipe directory", "path", rc.cfg.RecipePath, "error", err)
		return []datatypes.Recipe{}
	}

	recipies := make([]datatypes.Recipe, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)
		baseName := strings.TrimSuffix(name, ext)
		recipePath := filepath.Join(rc.cfg.RecipePath, name)

		recipe := datatypes.Recipe{Name: baseName}
		rc.enrichRecipeFromYAML(recipePath, &recipe)

		recipies = append(recipies, recipe)
	}

	sort.Slice(recipies, func(i, j int) bool {
		return recipies[i].Name < recipies[j].Name
	})

	return recipies
}

func (rc RecipeController) enrichRecipeFromYAML(recipePath string, recipe *datatypes.Recipe) {
	ext := strings.ToLower(filepath.Ext(recipePath))
	if ext != ".yaml" && ext != ".yml" {
		return
	}

	content, err := os.ReadFile(recipePath)
	if err != nil {
		slog.Warn("Could not read recipe file", "path", recipePath, "error", err)
		return
	}

	var doc recipeDocument
	if err := yaml.Unmarshal(content, &doc); err != nil {
		slog.Warn("Could not parse recipe yaml", "path", recipePath, "error", err)
		return
	}

	recipe.Description = doc.Description
}
