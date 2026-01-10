package config

type FilterType string

const (
	Priority FilterType = "priority"
	Status   FilterType = "status"
	Category FilterType = "category"
)

type PageInfo struct {
	HasNextPage bool
	StartCursor string
	EndCursor   string
}

type SectionConfig struct {
	Title   string
	Filters string
	Limit   *int
	Type    *FilterType
}

type PreviewConfig struct {
	Open  bool
	Width int
}

type Defaults struct {
	Preview PreviewConfig
	View    FilterType
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
		View: Status,
	},
}

func ParseConfig() Config {
	return *DefaultConfig
}
