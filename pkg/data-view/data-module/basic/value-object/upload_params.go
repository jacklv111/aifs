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

type UploadRawDataParam struct {
	FileMeta io.Reader
	// key: file name; value: file
	DataFileMap map[string]io.ReadSeeker
}

// 通过 multipart form 的形式上传数据
type UploadAnnotationParam struct {
	// 检查 annotation 的 raw data id 是否合法
	RawDataIdChecker func(dataItemIdList []uuid.UUID) error

	FileMeta io.Reader
	// key: raw data id; value: file
	DataFileMap map[string]io.ReadSeeker
	// key: raw data id; value: file name
	DataFileNameMap map[string]string
}

type UploadModelParams struct {
	DataFileMap map[string]io.ReadSeeker
	Pairs       map[string]string
}

type UploadDatasetZipParams struct {
	File     io.Reader
	FileName string
}

type UploadArtifactFileParams struct {
	File     io.Reader
	FileName string
}
