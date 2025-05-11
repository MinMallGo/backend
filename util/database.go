package util

import (
	"errors"
	"fmt"
	"log"
	"mall_backend/structure"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	database *structure.DatabaseClient
	dbOnce   sync.Once
)

// DBInstance 数据库连接单例。  .Client.Action()
func DBInstance() (*structure.DatabaseClient, error) {
	if database == nil {
		dbOnce.Do(func() {
			database = connectDatabase()
		})
	}

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
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database")
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("failed to connect database")
	}

	// 优化连接池配置
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)                                    // 根据服务器性能调整
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)                                    // 保持一定数量的空闲连接
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second) // 避免长时间持有连接
	sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)

	// 启动监控
	monitorDBStats(db)

	database = &structure.DatabaseClient{
		Client: db,
	}

	return database
}

func monitorDBStats(db *gorm.DB) {
	sqlDB, _ := db.DB() // 从 GORM 的 db 对象中获取底层的 *sql.DB 对象，以便获取连接池统计信息

	// 启动一个协程来定期监控连接池状态
	go func() {
		ticker := time.NewTicker(time.Minute) // 创建一个定时器，每隔一分钟触发一次
		defer ticker.Stop()                   // 函数结束时停止定时器

		for range ticker.C { // 持续监听定时器的通道
			stats := sqlDB.Stats() // 获取连接池的统计信息
			fmt.Printf("连接池状态：打开=%d, 空闲=%d, 等待=%d, 最大打开=%d\n",
				stats.OpenConnections,
				stats.Idle,
				stats.WaitCount,
				stats.MaxOpenConnections)
		}
	}()
}

// DBClient 我甚至不想把error带出去
func DBClient() *gorm.DB {
	// TODO 这里读取配置文件，如果是debug就加上debug，全局开启sql日志记录
	return database.Client //.Debug()
}
