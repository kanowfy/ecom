package config

type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"`
		Port int    `env:"SERVER_PORT" envDefault:"3005"`
	}
	Mail struct {
		Host     string `env:"MAIL_HOST" envDefaut:"smtp.mailtrap.io"`
		Port     int    `env:"MAIL_PORT" envDefault:"2525"`
		Username string `env:"MAIL_USERNAME" envDefault:"ad6c3dc5d1ad98"`
		Password string `env:"MAIL_PASSWORD" envDefault:"b45ea57aabed01"`
		Sender   string `env:"MAIL_SENDER" envDefault:"Ecom <no-reply@ecom.kanowfy.com>"`
	}
}
