package config

// postgres://root:psswrd@postgres:5432/catalog?sslmode=disable
type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"`
		Port int    `env:"SERVER_PORT" envDefault:"3001"`
	}
	DB struct {
		Url string `env:"DB_URL"`
	}
}
