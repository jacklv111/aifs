/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import "github.com/google/uuid"

const (
	TABLE_DATA_VIEW_ITEM = "data_view_items"
)

// data view 跟踪的数据
type DataViewItemDo struct {
	DataViewId uuid.UUID `gorm:"primaryKey"`
	DataItemId uuid.UUID `gorm:"primaryKey"`
}

func (DataViewItemDo) TableName() string {
	return TABLE_DATA_VIEW_ITEM
}

func GetDataItemIdList(doList []DataViewItemDo) (result []uuid.UUID) {
	for _, data := range doList {
		result = append(result, data.DataItemId)
	}
	return
}
