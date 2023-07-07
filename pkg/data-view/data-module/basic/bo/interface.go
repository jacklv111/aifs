/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"io"

	"github.com/google/uuid"
)

type DataInterface interface {

	// GetId 获取 data item 的 id
	//  @return uuid.UUID
	GetId() uuid.UUID

	// GetLocalPath 获取与 bo 对应的数据本地路径
	//  @return string
	GetLocalPath() string

	GetReader() io.Reader
	GetWriterAt() io.WriterAt

	// LoadFromLocal 读取本地文件，从中获取需要的数据到 bo 中
	//  @return error
	LoadFromLocal() error

	// LoadFromBuffer 数据从其他节点传过来，存放在 multiform file 中。从 multiform file 中读取数据
	//  @return error
	LoadFromBuffer() error
}

type AnnotationData interface {
	DataInterface
	GetLabels() []uuid.UUID
	GetAnnotationTemplateId() uuid.UUID
	GetRawDataId() uuid.UUID
}

func GetRawDataIdList(annoDataList []AnnotationData) []uuid.UUID {
	result := make([]uuid.UUID, 0)
	for _, data := range annoDataList {
		result = append(result, data.GetRawDataId())
	}
	return result
}
