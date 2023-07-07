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
	TABLE_ANNOTATION_TEMPLATE_EXT = "annotation_template_exts"
)

// AnnotationTemplateExtDo 存大字段，只有查询详情会查的数据
type AnnotationTemplateExtDo struct {
	AnnotationTemplateId uuid.UUID    `gorm:"uniqueIndex:anno_temp_ext_ukey,priority:1;<-:create"`
	WordList             WordListType `gorm:"column:word_list;type:text"`
	// soft delete
	DeleteAt soft_delete.DeletedAt `gorm:"softDelete:milli;uniqueIndex:anno_temp_ext_ukey,priority:2"`
}

func (AnnotationTemplateExtDo) TableName() string {
	return TABLE_ANNOTATION_TEMPLATE_EXT
}

func (do AnnotationTemplateExtDo) IsEmpty() bool {
	return do.WordList.IsEmpty()
}
