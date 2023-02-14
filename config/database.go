package config

import (
	"github.com/spf13/viper"
)

type Database struct {
	DBName   string `yaml:"dbname" json:"dbname" mapstructure:"dbname"`
	Host     string `yaml:"host" json:"host" mapstructure:"host"`             // the Host to connect to
	Port     int    `yaml:"port" json:"port" mapstructure:"port"`             // the Port to connect to
	User     string `yaml:"user" json:"user" mapstructure:"user"`             // the User for the database connection
	Password string `yaml:"password" json:"password" mapstructure:"password"` // the Password for the database connection
	SSLMode  string `yaml:"ssl" json:"ssl" mapstructure:"ssl"`
}

func (cfg Database) loadDefaultValues(v *viper.Viper) {
	v.SetDefault("Database.host", "127.0.0.1")
	v.SetDefault("Database.port", 3306)
	v.SetDefault("Database.user", "admin")
}
