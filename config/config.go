package config

type PreviewConfig struct {
	Open  bool
	Width int
}

type Defaults struct {
	Preview PreviewConfig
}

type Config struct {
	Defaults Defaults
}

var DefaultConfig = &Config{
	Defaults: Defaults{
		Preview: PreviewConfig{
			Open:  true,
			Width: 50,
		},
	},
}

func ParseConfig() Config {
	return *DefaultConfig
}
