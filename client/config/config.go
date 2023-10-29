package config

type Config struct {
	Url   Url    `mapstructure:"url" json:"url"`
	Token string `mapstructure:"token" json:"token"`
}

type Url struct {
	Base     string `mapstructure:"base" json:"base"`
	Register string `mapstructure:"register" json:"register"`
	Sms      string `mapstructure:"sms" json:"sms"`
	Ping     string `mapstructure:"ping" json:"ping"`
}
