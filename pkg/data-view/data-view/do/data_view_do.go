/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import (
	"database/sql"

	"github.com/google/uuid"
	"gorm.io/plugin/soft_delete"
)

const (
	TABLE_DATA_VIEW = "data_views"
)

// data view 的 meta data
// 支持软删除（相当于放回收站）。如果需要硬删除，就会彻底清除和其相关的所有数据，包括 data view item
type DataViewDo struct {
	ID                    uuid.UUID `gorm:"primaryKey;<-:create"`
	RelatedDataViewId     uuid.UUID `gorm:"<-:create"`
	AnnotationTemplateId  uuid.UUID
	Name                  string
	ViewType              string
	RawDataType           string
	Description           string
	Progress              float64
	ZipFormat             sql.NullString
	CommitId              sql.NullString
	Status                sql.NullString
	RawDataViewId         sql.NullString
	AnnotationViewId      sql.NullString
	TrainRawDataViewId    sql.NullString
	TrainAnnotationViewId sql.NullString
	ValRawDataViewId      sql.NullString
	ValAnnotationViewId   sql.NullString
	// allow read and create
	CreateAt int64 `gorm:"autoCreateTime:milli;<-:create"`
	// allow read and update
	UpdateAt int64 `gorm:"autoUpdateTime:milli;<-:update,create"`
	// soft delete
	DeleteAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
}

func (DataViewDo) TableName() string {
	return TABLE_DATA_VIEW
}

// 不同类型的 data view 共用这个 do，所以这里对不同 view 用到的属性进行说明
// 这里没有将不同 view 特有的属性拆到其他表，用不同的 do，主要考虑是降低代码的复杂性，不然每个 view 都要实现自己特有的方法
// 现在的设计，不同的 view 会将自己 view 的属性存入，没有该属性则该属性的位置为空值。在查询操作，转成 vo 的时候每个 view 写自己 assemble 拿自己的属性到 vo。
/* 公共属性
ID
Name
ViewType
Description
CreateAt
UpdateAt
DeleteAt
*/

/* raw data view 特有属性
RawDataType
*/

/* annotation 特有属性
RelatedDataViewId 在该 data view 的 raw data 上进行标注
AnnotationTemplateId 这个 data view 的 annotation 使用的 annotation template
*/

/* model 特有属性
Progress 当前模型训练的进度
CommitId 模型代码的 commit id
*/

/* dataset zip 特有属性
Progress 解压的进度
Status 解压状态
ZipFormat zip 文件解压后的目录结构
RawDataViewId 上传 raw data 的时候，raw data 上传的 data view
AnnotationViewId 上传 annotation 的时候，annotation 上传的 data view
TrainRawDataViewId 解析后训练集 raw data 所在的 data view
TrainAnnotationViewId 解析后训练集 annotation 所在的 data view
ValRawDataViewId 解析后验证集 raw data 所在的 data view
ValAnnotationViewId 解析后验证集 annotation 所在的 data view
AnnotationTemplateId 解析后 annotation template 的 id
*/
