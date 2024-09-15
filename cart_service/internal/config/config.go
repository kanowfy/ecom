package config

type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"`
		Port int    `env:"SERVER_PORT" envDefault:"3002"`
	}
	DB struct {
		Addr     string `env:"DB_URL"`
		Password string `env:"DB_PASSWORD"`
		Index    int    `env:"DB_INDEX" envDefault:"0"`
	}
}
