/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import "gorm.io/plugin/soft_delete"

// ListItem 查询列表时需要的数据
type ListItem struct {
	Id string

	// name of the annotation template
	Name string

	// Unix timestamp in ms
	CreateAt int64

	// the type of the annotation template
	Type string

	// the number of labels annotation template has
	LabelCount int32 `gorm:"column:label_count"`

	// soft delete
	DeleteAt soft_delete.DeletedAt `gorm:"softDelete:milli;uniqueIndex:label_ukey"`
}

func (ListItem) TableName() string {
	return TABLE_ANNOTATION_TEMPLATES
}
