/*
 * Created on Thu Jul 06 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package apigin

import (
	"io"

	"github.com/gin-gonic/gin"
)

func getAnnotationDataReaderFromMultipartForm(ctx *gin.Context) (fileMeta io.Reader, dataFileMap map[string]io.ReadSeeker, dataFileNameMap map[string]string, closers []io.Closer, err error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return
	}

	closers = make([]io.Closer, 0, len(form.File))
	dataFileMap = make(map[string]io.ReadSeeker, len(form.File))
	dataFileNameMap = make(map[string]string, len(form.File))
	for key, fileList := range form.File {
		if key == "fileMeta" {
			file, err := fileList[0].Open()
			if err != nil {
				return nil, nil, nil, nil, err
			}
			fileMeta = file
			closers = append(closers, file)
		} else {
			file, err := fileList[0].Open()
			if err != nil {
				return nil, nil, nil, nil, err
			}
			dataFileMap[key] = file
			dataFileNameMap[key] = fileList[0].Filename
			closers = append(closers, file)
		}
	}
	return
}

func getRawDataReaderFromMultipartForm(ctx *gin.Context) (fileMeta io.Reader, dataFileMap map[string]io.ReadSeeker, closers []io.Closer, err error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return
	}

	closers = make([]io.Closer, 0, len(form.File))
	dataFileMap = make(map[string]io.ReadSeeker, len(form.File))

	if fileList, ok := form.File["fileMeta"]; ok && len(fileList) == 1 {
		file, err := fileList[0].Open()
		if err != nil {
			return nil, nil, nil, err
		}
		fileMeta = file
		closers = append(closers, file)
	}

	for _, formFile := range form.File["files"] {
		file, err := formFile.Open()
		if err != nil {
			return nil, nil, nil, err
		}
		dataFileMap[formFile.Filename] = file
		closers = append(closers, file)
	}
	return
}

func getModelDataFromMultipartForm(ctx *gin.Context) (pairs map[string]string, dataFileMap map[string]io.ReadSeeker, closers []io.Closer, err error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return
	}

	closers = make([]io.Closer, 0, len(form.File))
	dataFileMap = make(map[string]io.ReadSeeker, len(form.File))
	pairs = make(map[string]string, len(form.Value))
	for key, fileList := range form.File {
		// 对 logs field 的特殊处理，这个字段可以放一个 list 的 log 文件
		if key == "logs" {
			for _, formFile := range fileList {
				file, err := formFile.Open()
				if err != nil {
					return nil, nil, nil, err
				}
				dataFileMap[formFile.Filename] = file
				closers = append(closers, file)
			}
		} else {
			file, err := fileList[0].Open()
			if err != nil {
				return nil, nil, nil, err
			}
			dataFileMap[key] = file
			closers = append(closers, file)
		}
	}
	for key, value := range form.Value {
		pairs[key] = value[0]
	}
	return
}
