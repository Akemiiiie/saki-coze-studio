/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rdb

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB 初始化 SQL Server 数据库连接
func NewDB() (*gorm.DB, error) {
	// 读取环境变量
	dsn := os.Getenv("SQL_CONN")
	if dsn == "" {
		return nil, fmt.Errorf("environment variable SQL_CONN is not set")
	}

	// 配置 GORM 日志
	gormLogger := logger.Default.LogMode(logger.Info)

	// 打开数据库
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlserver, dsn: %s, err: %w", dsn, err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 连接最大存活时间

	return db, nil
}
