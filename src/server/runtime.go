package server

type Runtime struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Dir      string `yaml:"dir"`
	LogLevel string `yaml:"log_level"`
}
