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
)

const (
	TABLE_ANNOTATION = "annotations"
)

// 该数据是只读的
type AnnotationDo struct {
	ID                   uuid.UUID      `gorm:"primaryKey;<-:create"`
	DataItemId           uuid.UUID      `gorm:"index:anno_data_item_idx;<-:create"`
	AnnotationTemplateId uuid.UUID      `gorm:"index:anno_anno_temp_idx;<-:create"`
	TextData             sql.NullString `gorm:"column:text_data;type:varchar(1024)"`
}

func (AnnotationDo) TableName() string {
	return TABLE_ANNOTATION
}
