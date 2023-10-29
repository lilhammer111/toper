package config

type ServerConfig struct {
	MysqlConfig MysqlConfig `mapstructure:"mysql-config" json:"mysql-config"`
	RedisConfig RedisConfig `mapstructure:"redis-config" json:"redis-config"`
	JwtConfig   JwtConfig   `mapstructure:"jwt-config" json:"jwt-config"`
}

type AddrConfig struct {
	Host string `mapstructure:"host" json:"host,omitempty"`
	Port string `mapstructure:"port" json:"port,omitempty"`
}

type MysqlConfig struct {
	Username   string     `mapstructure:"username" json:"username"`
	Password   string     `mapstructure:"password" json:"password"`
	DBName     string     `mapstructure:"db-name" json:"db-name"`
	AddrConfig AddrConfig `mapstructure:"addr-config" json:"addr"`
}

type RedisConfig struct {
	Expire     int        `mapstructure:"expire" json:"expire"`
	AddrConfig AddrConfig `mapstructure:"addr-config" json:"addr"`
}

type JwtConfig struct {
	JwtKey string `mapstructure:"jwt-key" json:"jwt-key"`
}
