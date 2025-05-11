package structure

import (
	"gorm.io/gorm"
)

type DataBaseConfig struct {
	Type            string `ini:"type"`
	Host            string `ini:"host"`
	Port            int    `ini:"port"`
	Name            string `ini:"name"`
	User            string `ini:"user"`
	Password        string `ini:"password"`
	MaxIdleConns    int    `ini:"max_idle_conns"`
	MaxOpenConns    int    `ini:"max_open_conns"`
	ConnMaxLifetime int    `ini:"conn_max_lifetime"`
	ConnMaxIdleTime int    `ini:"conn_max_idle_time"`
}

type DatabaseClient struct {
	Client *gorm.DB
}
