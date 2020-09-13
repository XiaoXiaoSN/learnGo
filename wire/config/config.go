package config

type Config struct {
	DBCfg DBConfig
}

type DBConfig struct {
	Prefix string
}

func NewConfig() *Config {
	return &Config{
		DBCfg: DBConfig{Prefix: "Hi~"},
	}
}
