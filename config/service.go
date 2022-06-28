package config

type Service struct {
	UserCenter       string `mapstructure:"user-center" json:"userCenter" yaml:"user-center"`
	UserCenterSecret string `mapstructure:"user-center-secret" json:"userCenterSecret" yaml:"user-center-secret"`
	AlertCenter      string `mapstructure:"alert-center" json:"alertCenter" yaml:"alert-center"`
}
