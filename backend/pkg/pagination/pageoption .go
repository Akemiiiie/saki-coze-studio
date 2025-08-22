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

package pagination

import (
	"github.com/coze-dev/coze-studio/backend/domain/knowledge/entity"
	"gorm.io/gorm"
)

// PageOption 用于分页和排序
type PageOption struct {
	Page     *int
	PageSize *int
	OrderBy  string // 例如 "created_at desc"
}

// 分页参数获取接口
type Paginable interface {
	GetLimit() *int
	GetOffset() *int
	GetOrder() *entity.Order
	GetOrderType() *entity.OrderType
}

// FilterFunc 用于构造 where 条件
// 接受一个 GORM 查询对象，返回带过滤条件的查询
type FilterFunc[T any] func(query *gorm.DB) *gorm.DB

// FindWithPagination 通用分页 + 总数查询
func FindWithPagination[T any](baseQuery *gorm.DB, filter FilterFunc[T], pageOpt *PageOption) (items []T, total int64, err error) {
	// ------------------------
	// 1. 构造带过滤条件的分页查询
	// ------------------------
	query := filter(baseQuery)

	// 排序和分页
	if pageOpt != nil {
		if pageOpt.OrderBy != "" {
			query = query.Order(pageOpt.OrderBy)
		}
		if pageOpt.Page != nil && pageOpt.PageSize != nil {
			offset := (*pageOpt.Page - 1) * (*pageOpt.PageSize)
			query = query.Limit(*pageOpt.PageSize).Offset(offset)
		}
	}

	// 查询分页数据
	if err = query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	// ------------------------
	// 2. 总数统计（干净查询，不带排序和分页）
	// ------------------------
	countQuery := filter(baseQuery.Session(&gorm.Session{NewDB: true}))
	if err = countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// 通用分页
// PageOptionFromPaginable 将实现了 Paginable 接口的结构体转换为通用分页参数
func PageOptionFromPaginable(p Paginable) *PageOption {
	if p == nil {
		return nil
	}

	// 获取 Limit，如果不存在或小于等于0，表示不分页
	limitPtr := p.GetLimit()
	if limitPtr == nil || *limitPtr <= 0 {
		return nil
	}
	limit := *limitPtr

	// 计算页码
	offsetPtr := p.GetOffset()
	page := 1
	if offsetPtr != nil {
		page = (*offsetPtr)/limit + 1
	}

	// 构造排序字段
	orderBy := "created_at DESC" // 默认排序
	if orderPtr := p.GetOrder(); orderPtr != nil {
		field := "created_at"
		if *orderPtr == entity.OrderUpdatedAt {
			field = "updated_at"
		}

		if orderTypePtr := p.GetOrderType(); orderTypePtr != nil && *orderTypePtr == entity.OrderTypeAsc {
			orderBy = field + " ASC"
		} else {
			orderBy = field + " DESC"
		}
	}

	// 返回通用分页结构
	return &PageOption{
		Page:     &page,
		PageSize: &limit,
		OrderBy:  orderBy,
	}
}
