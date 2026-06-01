package datatypes

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Db       string `env:"POSTGRES_DB"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
}

type Config struct {
	PgConf        PostgresConfig
	RecipePath    string `env:"RECIPIES_PATH" envDefault:"recipies"`
	StaticHtmlDir string `env:"STATIC_DIR"`
}
