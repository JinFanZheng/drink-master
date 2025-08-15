package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 数据库连接配置
	dsn := "root:mysql_5WPYdQn@tcp(47.113.121.132:3306)/vm_app_test?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// 获取所有表名
	var tables []string
	db.Raw("SHOW TABLES").Scan(&tables)

	fmt.Println("=== 数据库表列表 ===")
	for _, table := range tables {
		fmt.Printf("- %s\n", table)
	}

	// 分析每个表的结构
	for _, table := range tables {
		fmt.Printf("\n=== 表: %s ===\n", table)

		var columns []map[string]interface{}
		db.Raw("DESCRIBE " + table).Scan(&columns)

		for _, col := range columns {
			fmt.Printf("  %s: %s %s %s %s %s\n",
				col["Field"],
				col["Type"],
				col["Null"],
				col["Key"],
				col["Default"],
				col["Extra"])
		}
	}
}
