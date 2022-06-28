package config

type Mysql struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Config   string `json:"config" yaml:"config"`
	Dbname   string `json:"dbname" yaml:"dbname"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}
