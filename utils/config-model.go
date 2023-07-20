package utils

type Config struct {
	DB      DBConfig      `yaml:"db"`
	Handler HandlerConfig `yaml:"handler"`
	Models  []string      `yaml:"models"`
	Repos   []string      `yaml:"repositories"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type HandlerConfig struct {
	Path     string                   `yaml:"path"`
	Prefix   string                   `yaml:"prefix"`
	Versions map[string]VersionConfig `yaml:"versions"`
}

type VersionConfig struct {
	Prefix string                      `yaml:"prefix"`
	Routes map[string]RouteGroupConfig `yaml:"routes"`
}

type RouteGroupConfig struct {
	Prefix  string                 `yaml:"prefix"`
	Actions map[string]RouteConfig `yaml:"actions"`
}

type RouteConfig struct {
	Method  string `yaml:"method"`
	Path    string `yaml:"path"`
	Handler string `yaml:"handler"`
}
