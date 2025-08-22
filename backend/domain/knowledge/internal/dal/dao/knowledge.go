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

package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coze-dev/coze-studio/backend/domain/knowledge/entity"
	"github.com/coze-dev/coze-studio/backend/domain/knowledge/internal/dal/model"
	"github.com/coze-dev/coze-studio/backend/domain/knowledge/internal/dal/query"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/pagination"
)

type KnowledgeDAO struct {
	DB    *gorm.DB
	Query *query.Query
}

func (dao *KnowledgeDAO) Create(ctx context.Context, knowledge *model.Knowledge) error {
	return dao.Query.Knowledge.WithContext(ctx).Create(knowledge)
}
func (dao *KnowledgeDAO) Upsert(ctx context.Context, knowledge *model.Knowledge) error {
	return dao.Query.Knowledge.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(knowledge)
}
func (dao *KnowledgeDAO) Update(ctx context.Context, knowledge *model.Knowledge) error {
	k := dao.Query.Knowledge
	knowledge.UpdatedAt = time.Now().UnixMilli()
	err := k.WithContext(ctx).Where(k.ID.Eq(knowledge.ID)).Save(knowledge)
	return err
}

func (dao *KnowledgeDAO) Delete(ctx context.Context, id int64) error {
	k := dao.Query.Knowledge
	_, err := k.WithContext(ctx).Where(k.ID.Eq(id)).Delete()
	return err
}

func (dao *KnowledgeDAO) MGetByID(ctx context.Context, ids []int64) ([]*model.Knowledge, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	k := dao.Query.Knowledge
	pos, err := k.WithContext(ctx).Where(k.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}

	return pos, nil
}

//过滤出在启用状态的知识（Knowledge）记录
func (dao *KnowledgeDAO) FilterEnableKnowledge(ctx context.Context, knowledgeIDs []int64) ([]*model.Knowledge, error) {
	if len(knowledgeIDs) == 0 {
		return nil, nil
	}
	k := dao.Query.Knowledge
	knowledgeModels, err := k.WithContext(ctx).
		Select(k.ID, k.FormatType).
		Where(k.ID.In(knowledgeIDs...)).
		Where(k.Status.Eq(int32(entity.DocumentStatusEnable))).
		Find()

	return knowledgeModels, err
}

func (dao *KnowledgeDAO) InitTx() (tx *gorm.DB, err error) {
	tx = dao.DB.Begin()
	if tx.Error != nil {
		return nil, err
	}
	return
}

func (dao *KnowledgeDAO) UpdateWithTx(ctx context.Context, tx *gorm.DB, knowledgeID int64, updateMap map[string]interface{}) error {
	return tx.WithContext(ctx).Model(&model.Knowledge{}).Where("id = ?", knowledgeID).Updates(updateMap).Error
}

func (dao *KnowledgeDAO) FindKnowledgeByCondition(ctx context.Context, opts *entity.WhereKnowledgeOption) (knowledge []*model.Knowledge, total int64, err error) {
	if opts == nil {
		return nil, 0, nil
	}

	// 显式指定泛型类型 *model.Knowledge
	items, totalCount, err := pagination.FindWithPagination[*model.Knowledge](
		dao.Query.Knowledge.WithContext(ctx).UnderlyingDB(), // 移除 .Debug()，除非调试需要
		func(db *gorm.DB) *gorm.DB {
			q := db.Model(&model.Knowledge{})

			// 过滤条件
			if opts.Query != nil && len(*opts.Query) > 0 {
				q = q.Where("name LIKE ?", "%"+*opts.Query+"%")
			}
			if opts.Name != nil && len(*opts.Name) > 0 {
				q = q.Where("name = ?", *opts.Name)
			}
			if len(opts.KnowledgeIDs) > 0 {
				q = q.Where("id IN ?", opts.KnowledgeIDs)
			}
			if ptr.From(opts.AppID) != 0 {
				q = q.Where("app_id = ?", ptr.From(opts.AppID))
			} else if len(opts.KnowledgeIDs) == 0 {
				q = q.Where("app_id = 0")
			}
			if ptr.From(opts.SpaceID) != 0 {
				q = q.Where("space_id = ?", *opts.SpaceID)
			}
			if len(opts.Status) > 0 {
				q = q.Where("status IN ?", opts.Status)
			}
			if opts.UserID != nil && ptr.From(opts.UserID) != 0 {
				q = q.Where("creator_id = ?", *opts.UserID)
			}
			if opts.FormatType != nil {
				q = q.Where("format_type = ?", int32(*opts.FormatType))
			}

			return q
		},
		&pagination.PageOption{
			Page:     opts.Page,
			PageSize: opts.PageSize,
			OrderBy: func() string {
				if opts.Order != nil {
					if *opts.Order == entity.OrderCreatedAt {
						if opts.OrderType != nil && *opts.OrderType == entity.OrderTypeAsc {
							return "created_at ASC"
						}
						return "created_at DESC"
					} else if *opts.Order == entity.OrderUpdatedAt {
						if opts.OrderType != nil && *opts.OrderType == entity.OrderTypeAsc {
							return "updated_at ASC"
						}
						return "updated_at DESC"
					}
				}
				return "created_at DESC"
			}(),
		},
	)

	if err != nil {
		return nil, 0, err
	}

	knowledge = items
	total = totalCount
	return
}


func (dao *KnowledgeDAO) GetByID(ctx context.Context, id int64) (*model.Knowledge, error) {
	k := dao.Query.Knowledge
	knowledge, err := k.WithContext(ctx).Where(k.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return knowledge, nil
}
