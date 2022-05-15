package conf

type Config struct {
	URL string
	PKA string
}

func (c *Config) Check() {
}
