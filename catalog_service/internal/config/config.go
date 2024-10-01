package config

type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"`
		Port int    `env:"SERVER_PORT" envDefault:"3001"`
	}
	DB struct {
		Url string `env:"POSTGRES_URL"`
	}
	Otel struct {
		GrpcEndpoint string `env:"OTLP_GRPC_ENDPOINT"`
	}
}
