package config

type Config struct {
	Addr string
	Port int
}

func NewConfig(addr string, port int) *Config {
	return &Config{
		Addr: addr,
		Port: port,
	}
}
