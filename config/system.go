package config

type System struct {
	Level string `json:"level" yaml:"level"`
	Path  string `json:"path" yaml:"path"`
	Port  string `json:"port" yaml:"port"`
	Host  string `json:"host" yaml:"host"`
}
