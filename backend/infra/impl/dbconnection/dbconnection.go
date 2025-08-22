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

package dbconnection

import (
	"fmt"
	"os"

	"github.com/coze-dev/coze-studio/backend/types/consts"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func New() (*gorm.DB, error) {
	dbType := os.Getenv("DB_TYPE")
	if dbType == consts.SQLSERVER {
		return NewSqlServer(os.Getenv("SQL_CONN"))
	}
	return NewMySql(os.Getenv("SQL_CONN"))
}

func NewMySql(sqlConn string) (*gorm.DB, error) {
		db, err := gorm.Open(mysql.Open(sqlConn))
		if err != nil {
			return nil, fmt.Errorf("mysql open, dsn: %s, err: %w", sqlConn, err)
		}
		return db, nil
}

func NewSqlServer(sqlConn string) (*gorm.DB, error) {
		db, err := gorm.Open(sqlserver.Open(sqlConn))
		if err != nil {
			return nil, fmt.Errorf("sqlserver open, dsn: %s, err: %w", sqlConn, err)
		}
		return db, nil
}
