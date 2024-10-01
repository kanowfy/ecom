package config

type UpstreamAddr struct {
	Catalog  string `env:"CATALOG_SVC_ADDR"`
	Cart     string `env:"CART_SVC_ADDR"`
	Shipping string `env:"SHIPPING_SVC_ADDR"`
	Email    string `env:"EMAIL_SVC_ADDR"`
	Payment  string `env:"PAYMENT_SVC_ADDR"`
	Order    string `env:"ORDER_SVC_ADDR"`
}

type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"`
		Port int    `env:"SERVER_PORT" envDefault:"4000"`
	}
	Cookie struct {
		SID    string `env:"COOKIE_SID" envDefault:"ecom_sid"`
		MaxAge int    `env:"COOKIE_MAXAGE" envDefault:"86400"`
	}
	UpstreamAddr UpstreamAddr
	Otel         struct {
		GrpcEndpoint string `env:"OTLP_GRPC_ENDPOINT"`
	}
}
