package projects

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dt "claw-destine.com/camboose/core/datatypes"
	"gopkg.in/yaml.v3"
)

type RecipeManager struct {
	Conf dt.Config
}

func (rm *RecipeManager) ListRecipes() ([]dt.Recipe, error) {
	entries, err := os.ReadDir(rm.Conf.RecipePath)
	if err != nil {
		return nil, fmt.Errorf("read recipes directory %q: %w", rm.Conf.RecipePath, err)
	}

	recipes := make([]dt.Recipe, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(rm.Conf.RecipePath, entry.Name())
		manifestPath, err := findManifestPath(dirPath)
		if err != nil {
			return nil, err
		}
		if manifestPath == "" {
			continue
		}

		recipe, err := loadRecipe(entry.Name(), dirPath, manifestPath)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

type recipeManifest struct {
	Id            string   `yaml:"id"`
	Version       string   `yaml:"version"`
	Description   string   `yaml:"description"`
	SpecHierarchy []string `yaml:"specHierarchy"`
}

func findManifestPath(dirPath string) (string, error) {
	manifestNames := []string{"manifest.yml", "manifest.yaml"}
	for _, fileName := range manifestNames {
		manifestPath := filepath.Join(dirPath, fileName)
		info, err := os.Stat(manifestPath)
		if err == nil {
			if info.IsDir() {
				continue
			}
			return manifestPath, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("check manifest in %q: %w", dirPath, err)
		}
	}

	return "", nil
}

func loadRecipe(defaultId, dirPath, manifestPath string) (dt.Recipe, error) {
	manifestContent, err := os.ReadFile(manifestPath)
	if err != nil {
		return dt.Recipe{}, fmt.Errorf("read manifest %q: %w", manifestPath, err)
	}

	var manifest recipeManifest
	if err := yaml.Unmarshal(manifestContent, &manifest); err != nil {
		return dt.Recipe{}, fmt.Errorf("parse manifest %q: %w", manifestPath, err)
	}

	recipe := dt.Recipe{
		Id:           strings.TrimSpace(manifest.Id),
		Description:  strings.TrimSpace(manifest.Description),
		Version:      strings.TrimSpace(manifest.Version),
		SpecHierachy: manifest.SpecHierarchy,
	}

	if recipe.Id == "" {
		recipe.Id = defaultId
	}

	promptPath := filepath.Join(dirPath, "design.prompt.md")
	prompt, err := os.ReadFile(promptPath)
	if err == nil {
		recipe.SpecPrompt = string(prompt)
	} else if !errors.Is(err, os.ErrNotExist) {
		return dt.Recipe{}, fmt.Errorf("read design prompt %q: %w", promptPath, err)
	}

	return recipe, nil
}
