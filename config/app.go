package config

type App struct {
	System  System  `json:"system" yaml:"system"`
	Mysql   Mysql   `json:"mysql" yaml:"mysql"`
	Service Service `json:"service" yaml:"service"`
	Jwt     JWT     `json:"jwt" yaml:"jwt" mapstructure:"jwt"`
	Redis   Redis   `json:"redis" yaml:"redis" mapstructure:"redis"`
}
