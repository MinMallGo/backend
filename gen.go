package main

//
//import (
//	"gorm.io/driver/mysql"
//	"gorm.io/gen"
//	"gorm.io/gorm"
//)
//
//func main() {
//	// 配置生成器
//	g := gen.NewGenerator(gen.Config{
//		OutPath: "./dao/model", // 生成代码的输出目录
//	})
//
//	// 连接数据库（需自行配置DSN）
//	gormdb, err := gorm.Open(mysql.Open("root:root@(127.0.0.1:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"))
//	if err != nil {
//		panic(err)
//	}
//	g.UseDB(gormdb) // reuse your gorm db
//
//	// 仅生成所有表的模型结构体
//	g.GenerateAllTable()
//	g.Execute()
//}
