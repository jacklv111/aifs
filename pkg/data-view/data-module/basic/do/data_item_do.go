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
)

const (
	TABLE_DATA_ITEM = "data_items"
)

// 该数据是只读的
type DataItemDo struct {
	ID   uuid.UUID `gorm:"primaryKey;<-:create"`
	Name string
	// example: image, video, text, ...
	Type string
	// allow read and create
	CreateAt int64 `gorm:"autoCreateTime:milli;<-:create"`
}

func (DataItemDo) TableName() string {
	return TABLE_DATA_ITEM
}

func GetDataItemIdList(dataItemList []DataItemDo) []uuid.UUID {
	idList := make([]uuid.UUID, 0, len(dataItemList))
	for _, data := range dataItemList {
		idList = append(idList, data.ID)
	}
	return idList
}
