package builder

type Builder struct {
	Config Config
}

type Config struct {
	Jobs []Job `mapstructure:"jobs"`
}

type Job struct {
	Name string `mapstructure:"name"`
	Stages []Stage `mapstructure:"stages"`
}

type Stage struct {
	Commands []string `mapstructure:"commands"`
	Env map[string]string `mapstructure:"env"`
	Image string `mapstructure:"image"`
	Name string `mapstructure:"name"`
}
