package datatypes

type Config struct {
	DbPath     string `env:"DBPATH" envDefault:"data/clover-store"`
	RecipePath string `env:"RECIPEPATH" envDefault:"recipies"`
}
