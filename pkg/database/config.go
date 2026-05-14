package database

type Driver string

var Pgx Driver = "pgx"
var Mysql Driver = "mysql"

type Config struct {
	Driver        Driver `mapstructure:"DRIVER"`
	Host          string `mapstructure:"HOST"`
	Port          int    `mapstructure:"PORT"`
	User          string `mapstructure:"USER"`
	Password      string `mapstructure:"PASSWORD"`
	Name          string `mapstructure:"NAME"`
	QueryExecMode string `mapstructure:"QUERY_EXEC_MODE"`
	SslMode       string `mapstructure:"SSL_MODE"`
}
