package config

type Config struct {
	Host string `env:"SERVER_HOST"`
	Port int    `env:"SERVER_PORT" envDefault:"3004"`
}
