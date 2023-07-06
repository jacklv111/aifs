/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

import (
	"io"

	"github.com/google/uuid"
)

// 存储请求的输入参数
type UploadParams struct {
	// 数据类型
	DataType     string
	DataItemList []ReaderItem
}

func (params *UploadParams) AddItem(id uuid.UUID, reader io.Reader, name string) {
	params.DataItemList = append(params.DataItemList, ReaderItem{DataItemId: id, Reader: reader, Name: name})
}

func (params *UploadParams) GetIdList() (res []uuid.UUID) {
	for _, data := range params.DataItemList {
		res = append(res, data.DataItemId)
	}
	return
}

func (params *UploadParams) IsEmpty() bool {
	return len(params.DataItemList) == 0
}

type ReaderItem struct {
	// 数据的 id
	DataItemId uuid.UUID
	// 数据上传或者下载时的文件
	Reader io.Reader
	// 数据的名称。有些 data item 会有多个要存储的数据，这里需要用 name 来区分
	Name string
}

func (item ReaderItem) GetUniqueName() string {
	return item.DataItemId.String() + "-" + item.Name
}
