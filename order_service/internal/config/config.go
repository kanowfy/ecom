package config

type UpstreamAddr struct {
	Catalog  string `env:"CATALOG_SVC_ADDR"`
	Cart     string `env:"CART_SVC_ADDR"`
	Shipping string `env:"SHIPPING_SVC_ADDR"`
	Email    string `env:"EMAIL_SVC_ADDR"`
	Payment  string `env:"PAYMENT_SVC_ADDR"`
}

type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"`
		Port int    `env:"SERVER_PORT" envDefault:"3006"`
	}
	UpstreamAddr UpstreamAddr
}