package config

type Config struct {
	Url   Url    `mapstructure:"url" json:"url"`
	Token string `mapstructure:"token" json:"token"`
}

type Url struct {
	Root string `mapstructure:"root" json:"root"`

	Register string `mapstructure:"register" json:"register"`
	Login    string `mapstructure:"login" json:"login"`
	Sms      string `mapstructure:"sms" json:"sms"`
	Ping     string `mapstructure:"ping" json:"ping"`
	User     string `mapstructure:"user" json:"user"`
	Toper    string `mapstructure:"toper" json:"toper"`
	Done     string `mapstructure:"done" json:"done"`
	History  string `mapstructure:"history" json:"history"`
}
