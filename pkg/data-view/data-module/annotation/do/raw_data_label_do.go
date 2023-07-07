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
	TABLE_RAW_DATA_LABEL = "raw_data_labels"
)

// 该数据是只读的
type RawDataLabelDo struct {
	AnnotationId uuid.UUID `gorm:"index:raw_data_labels_anno_id_idx;<-:create"`
	RawDataId    uuid.UUID
	LabelId      uuid.UUID `gorm:"index:raw_data_labels_label_id_idx;<-:create"`
}

func (RawDataLabelDo) TableName() string {
	return TABLE_RAW_DATA_LABEL
}

func GetLabelIdStrList(dataList []RawDataLabelDo) []string {
	res := make([]string, 0)
	for _, data := range dataList {
		res = append(res, data.LabelId.String())
	}
	return res
}
