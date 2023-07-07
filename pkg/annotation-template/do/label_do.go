/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import (
	"github.com/google/uuid"
	"gorm.io/plugin/soft_delete"
)

const (
	TABLE_LABELS = "labels"
)

type LabelDo struct {
	ID                   uuid.UUID `gorm:"primaryKey;<-:create"`
	AnnotationTemplateId uuid.UUID `gorm:"uniqueIndex:label_ukey,priority:1;<-:create"`
	Name                 string    `gorm:"uniqueIndex:label_ukey,priority:2"`
	SuperCategoryName    string
	Color                int32
	// allow read and create
	CreateAt int64 `gorm:"autoCreateTime:milli;<-:create"`
	// allow read and update
	UpdateAt int64 `gorm:"autoUpdateTime:milli;<-:update,create"`
	// soft delete
	DeleteAt soft_delete.DeletedAt `gorm:"softDelete:milli;uniqueIndex:label_ukey,priority:3"`

	KeyPointDef      KeyPointDefType      `gorm:"column:key_point_def;type:varchar(512)"`
	KeyPointSkeleton KeyPointSkeletonType `gorm:"column:key_point_skeleton;type:varchar(512)"`

	CoverImageUrl string `gorm:"column:cover_image_url;type:varchar(512)"`
}

func (LabelDo) TableName() string {
	return TABLE_LABELS
}
