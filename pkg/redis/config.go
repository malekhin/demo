package redis

type Config struct {
	Host     string `mapstructure:"redis_host"`
	Port     int    `mapstructure:"redis_port"`
	Username string `mapstructure:"redis_username"`
	Password string `mapstructure:"redis_password"`
	Prefix   string `mapstructure:"redis_prefix"`
}
