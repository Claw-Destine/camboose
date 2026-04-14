package recipies

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type RecipeConfig struct {
	RecipePath string `env:"PATH"`
}

type ReqStatus string

const (
	ReqStatusNew       ReqStatus = "new"
	ReqStatusDefined   ReqStatus = "defined"
	ReqStatusDelivered ReqStatus = "delivered"
)

type ReqEntity string

const (
	ReqEntityEpic                ReqEntity = "epic"
	ReqEntityStory               ReqEntity = "story"
	ReqEntityAcceptanceCriterion ReqEntity = "acceptance_cryterion"
)

type DesignEntity string

const (
	DesignEntityView DesignEntity = "view"
)

type ReqDescription struct {
	Me    ReqEntity
	Child *ReqEntity
}

type ProjectManagement struct {
	Relations []ReqDescription
}

type Recipe struct {
	Name               string
	Description        string
	Project_Management ProjectManagement
}

type RecipeController struct {
	cfg RecipeConfig
}

func NewRecipeController(cfg RecipeConfig) RecipeController {
	return RecipeController{cfg: cfg}
}

func (rc RecipeController) ListRecipies() []string {
	entries, err := os.ReadDir(rc.cfg.RecipePath)
	if err != nil {
		slog.Error("Could not read recipe directory", "path", rc.cfg.RecipePath, "error", err)
		return []string{}
	}

	recipies := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)
		recipies = append(recipies, strings.TrimSuffix(name, ext))

	}

	// sort.Slice(recipies, func(i, j int) bool {
	// 	return recipies[i] < recipies[j]
	// })

	return recipies
}
