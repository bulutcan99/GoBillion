package baseconfig

type (
	App struct {
		Name string `env-required:"true" yaml:"name" env:"APP_NAME"`
	}

	Log struct {
		Level int `env-required:"true" yaml:"level"   env:"LOG_LEVEL"`
	}
)
