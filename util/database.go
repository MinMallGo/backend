package util

import (
	"errors"
	"fmt"
	"log"
	"mall_backend/structure"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	database *structure.DatabaseClient
	dbOnce   sync.Once
)

// DBInstance 数据库连接单例。  .Client.Action()
func DBInstance() (*structure.DatabaseClient, error) {
	dbOnce.Do(func() {
		database = connectDatabase()
	})

	if database == nil {
		return nil, errors.New("failed to initialize Database client")
	}

	return database, nil
}

func connectDatabase() *structure.DatabaseClient {
	config, err := GetDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Name)
	open, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database")
		return nil
	}

	database = &structure.DatabaseClient{
		Client: open,
	}

	return database
}

// DBClient 我甚至不想把error带出去
func DBClient() *gorm.DB {
	return database.Client
}
