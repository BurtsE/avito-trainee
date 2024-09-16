package config

type Config struct {
	BillingDB `json:"billing_db"`
	Service   `json:"service"`
}

type Service struct {
	Port string `json:"server_port"`
}

type BillingDB struct {
	Host     string `env:"POSTGRES_HOST,notEmpty"`
	Port     string `env:"POSTGRES_PORT,notEmpty"`
	DB       string `env:"POSTGRES_DATABASE,notEmpty"`
	User     string `env:"POSTGRES_USERNAME,notEmpty"`
	Password string `env:"POSTGRES_PASSWORD,notEmpty"`
	MaxConns int    `json:"max_conns"`
	Sslmode  string `json:"sslmode"`
}
