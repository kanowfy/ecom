package config

type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"`
		Port int    `env:"SERVER_PORT" envDefault:"3002"`
	}
	DB struct {
		Addr     string `env:"REDIS_URL"`
		Password string `env:"REDIS_PASSWORD"`
		Index    int    `env:"REDIS_DB_INDEX" envDefault:"0"`
	}
	Otel struct {
		GrpcEndpoint string `env:"OTLP_GRPC_ENDPOINT"`
	}
}
