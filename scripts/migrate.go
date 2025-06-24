package main

import (
	"fmt"
	"log"
	"wz-wenzhan-backend/internal/config"
	"wz-wenzhan-backend/internal/model"
)

func main() {
	// 加载配置
	cfg := config.Load()
	
	// 初始化数据库连接
	db := config.InitDB(cfg)

	// 自动迁移数据库表
	err := db.AutoMigrate(
		&model.User{},
		&model.Folder{},
		&model.Document{},
		&model.Activity{},
		&model.RecycleItem{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	fmt.Println("数据库迁移成功！")
}
