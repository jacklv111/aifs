/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package handler

import (
	aifsclientgo "github.com/jacklv111/aifs-client-go"
	"github.com/jacklv111/common-sdk/utils"
)

type RawDataImagesHandler struct {
	HandlerBase
}

func NewRawDataImagesHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) Handler {
	return &RawDataImagesHandler{
		HandlerBase: HandlerBase{
			client:     client,
			dataView:   dataView,
			datasetDir: datasetDir,
		},
	}
}

func (handler *RawDataImagesHandler) Exec() error {
	allDataPathList, err := utils.ReadAllFiles(handler.datasetDir)
	if err != nil {
		return err
	}
	imageDataPathList := make([]string, 0)
	for _, filePath := range allDataPathList {
		if utils.IsImageFromFile(filePath) {
			imageDataPathList = append(imageDataPathList, filePath)
		}
	}

	rawDataViewId, err := handler.CreateRawDataView(aifsclientgo.IMAGE)
	if err != nil {
		return err
	}

	err = handler.UploadRawData(rawDataViewId, imageDataPathList)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(100.0, "dataset zip view decompression completed"); err != nil {
		return err
	}
	return nil
}
