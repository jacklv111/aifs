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
	TABLE_ANNOTATION_TEMPLATES = "annotation_templates"
)

type AnnotationTemplateDo struct {
	ID   uuid.UUID `gorm:"primaryKey;<-:create"`
	Name string
	// annotation template type
	Type        string
	Description string
	// allow read and create
	CreateAt int64 `gorm:"autoCreateTime:milli;<-:create"`
	// allow read and update
	UpdateAt int64 `gorm:"autoUpdateTime:milli;<-:update,create"`
	// soft delete
	DeleteAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
}

func (AnnotationTemplateDo) TableName() string {
	return TABLE_ANNOTATION_TEMPLATES
}
