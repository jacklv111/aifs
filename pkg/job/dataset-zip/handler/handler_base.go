/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	"github.com/jacklv111/aifs/pkg/job/dataset-zip/constant"
	"github.com/jacklv111/common-sdk/annotation"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/log"
)

type HandlerBase struct {
	client     *aifsclientgo.APIClient
	dataView   *aifsclientgo.DataViewDetails
	datasetDir string
}

func (handler *HandlerBase) UpdateProgress(progress float32, status string) error {
	dataViewId := *handler.dataView.Id
	var updateReq aifsclientgo.UpdateDatasetZipRequest
	updateReq.SetProgress(progress)
	updateReq.SetStatus(status)

	resp, err := handler.client.DataViewApi.UpdateDatasetZipView(context.Background(), dataViewId).UpdateDatasetZipRequest(updateReq).Execute()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update dataset zip view got incorrect status code %d", resp.StatusCode)
	}
	return nil
}

func (handler *HandlerBase) GetRawDataHashIdMap(rawDataViewId string) (map[string]string, error) {
	log.Info("get raw data hash id map")
	ans := make(map[string]string, 0)
	offset := 0
	limit := 1000
	for {
		log.Infof("get raw data hash id map offset %d", offset)
		res, _, err := handler.client.DataViewApi.
			GetRawDataHashListInDataView(context.Background(), rawDataViewId).
			Offset(int32(offset)).
			Limit(int32(limit)).
			Execute()
		if err != nil {
			return nil, err
		}
		for _, data := range res {
			ans[*data.Sha256] = *data.RawDataId
		}

		if len(res) == 0 || len(res) < limit {
			break
		}

		offset += limit
	}
	return ans, nil
}

func (handler *HandlerBase) UploadRawData(dataViewId string, dataPathList []string) error {
	count := 0
	return collection.BatchRange(dataPathList, constant.UPLOAD_IMAGE_BATCH, func(batch []string) error {
		count += len(batch)
		log.Infof("upload raw data %d", count)

		files := make([]aifsclientgo.FormFile, 0)
		for _, filePath := range batch {
			file, err := os.Open(filePath)
			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()
			if err != nil {
				return err
			}
			byteArray, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			files = append(files, aifsclientgo.FormFile{
				Key:      "files",
				Value:    byteArray,
				FileName: filepath.Base(filePath),
			})
		}

		resp, err := handler.client.DataViewUploadApi.UploadRawDataToDataView(context.Background(), dataViewId).Files(files).Execute()
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("upload raw data status code %d", resp.StatusCode)
		}
		return nil
	})
}

func (handler *HandlerBase) UploadFormFiles(dataViewId string, formFiles []aifsclientgo.FormFile) error {
	resp, err := handler.client.DataViewUploadApi.UploadAnnotationToDataView(context.Background(), dataViewId).Files(formFiles).Execute()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload raw data failed, status code: %d", resp.StatusCode)
	}
	return nil
}

func (handler *HandlerBase) CreateAllDataView(rawDataType aifsclientgo.RawDataType, annoTemp annotation.AnnotationTemplateDetails) (trainRawDataViewId, trainAnnotationViewId, valRawDataViewId, valAnnotationViewId string, err error) {
	if handler.dataView.TrainRawDataViewId == nil || handler.dataView.TrainAnnotationViewId == nil || handler.dataView.ValRawDataViewId == nil || handler.dataView.ValAnnotationViewId == nil {
		trainRawDataViewId, trainAnnotationViewId, err = handler.createDataView(rawDataType, annoTemp, "-train")
		if err != nil {
			return
		}
		valRawDataViewId, valAnnotationViewId, err = handler.createDataView(rawDataType, annoTemp, "-val")
		if err != nil {
			return
		}
		var updateReq aifsclientgo.UpdateDatasetZipRequest
		var resp *http.Response
		updateReq.SetTrainRawDataViewId(trainRawDataViewId)
		updateReq.SetTrainAnnotationViewId(trainAnnotationViewId)
		updateReq.SetValRawDataViewId(valRawDataViewId)
		updateReq.SetValAnnotationViewId(valAnnotationViewId)
		resp, err = handler.client.DataViewApi.UpdateDatasetZipView(context.Background(), *handler.dataView.Id).UpdateDatasetZipRequest(updateReq).Execute()
		if err != nil {
			return
		}
		if resp.StatusCode != http.StatusOK {
			return "", "", "", "", fmt.Errorf("update dataset zip view got incorrect status code %d", resp.StatusCode)
		}
		return
	}
	return *handler.dataView.TrainRawDataViewId, *handler.dataView.TrainAnnotationViewId, *handler.dataView.ValRawDataViewId, *handler.dataView.ValAnnotationViewId, nil
}

func (handler *HandlerBase) CreateRawDataView(rawDataType aifsclientgo.RawDataType) (rawDataViewId string, err error) {
	if handler.dataView.RawDataViewId == nil {
		var rawDataViewReq aifsclientgo.CreateDataViewRequest
		rawDataViewReq.SetRawDataType(rawDataType)
		rawDataViewReq.SetDataViewName(*handler.dataView.Name)
		rawDataViewReq.SetViewType(aifsclientgo.RAW_DATA)
		rawDataViewResp, _, err := handler.client.DataViewApi.CreateDataView(context.Background()).CreateDataViewRequest(rawDataViewReq).Execute()
		if err != nil {
			return "", err
		}
		rawDataViewId = rawDataViewResp.GetDataViewId()
		var updateReq aifsclientgo.UpdateDatasetZipRequest
		var resp *http.Response
		updateReq.SetRawDataViewId(rawDataViewId)

		resp, err = handler.client.DataViewApi.UpdateDatasetZipView(context.Background(), *handler.dataView.Id).UpdateDatasetZipRequest(updateReq).Execute()
		if err != nil {
			return "", err
		}
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("update dataset zip view got incorrect status code %d", resp.StatusCode)
		}
		return rawDataViewId, nil
	}
	return *handler.dataView.RawDataViewId, nil
}

func (handler *HandlerBase) CreateRawDataAndAnnotationView(rawDataType aifsclientgo.RawDataType, annoTemp annotation.AnnotationTemplateDetails) (rawDataViewId, annotationViewId string, err error) {
	if handler.dataView.RawDataViewId == nil || handler.dataView.AnnotationViewId == nil {
		rawDataViewId, annotationViewId, err = handler.createDataView(rawDataType, annoTemp, "")
		if err != nil {
			return
		}

		var updateReq aifsclientgo.UpdateDatasetZipRequest
		var resp *http.Response
		updateReq.SetRawDataViewId(rawDataViewId)
		updateReq.SetAnnotationViewId(annotationViewId)

		resp, err = handler.client.DataViewApi.UpdateDatasetZipView(context.Background(), *handler.dataView.Id).UpdateDatasetZipRequest(updateReq).Execute()
		if err != nil {
			return
		}
		if resp.StatusCode != http.StatusOK {
			return "", "", fmt.Errorf("update dataset zip view got incorrect status code %d", resp.StatusCode)
		}
		return
	}
	return *handler.dataView.RawDataViewId, *handler.dataView.AnnotationViewId, nil
}

func (handler *HandlerBase) createDataView(rawDataType aifsclientgo.RawDataType, annoTemp annotation.AnnotationTemplateDetails, viewNameSuffix string) (rawDataViewId, annotationViewId string, err error) {
	var rawDataViewReq aifsclientgo.CreateDataViewRequest
	rawDataViewReq.SetRawDataType(rawDataType)
	rawDataViewReq.SetDataViewName(*handler.dataView.Name + viewNameSuffix)
	rawDataViewReq.SetViewType(aifsclientgo.RAW_DATA)
	rawDataViewResp, _, err := handler.client.DataViewApi.CreateDataView(context.Background()).CreateDataViewRequest(rawDataViewReq).Execute()
	if err != nil {
		return "", "", err
	}

	var annoViewReq aifsclientgo.CreateDataViewRequest
	annoViewReq.SetDataViewName(*handler.dataView.Name + viewNameSuffix)
	annoViewReq.SetViewType(aifsclientgo.ANNOTATION)
	annoViewReq.SetAnnotationTemplateId(*annoTemp.Id)
	annoViewReq.SetRelatedDataViewId(rawDataViewResp.GetDataViewId())
	annoViewResp, _, err := handler.client.DataViewApi.CreateDataView(context.Background()).CreateDataViewRequest(annoViewReq).Execute()
	if err != nil {
		return "", "", err
	}
	return *rawDataViewResp.DataViewId, *annoViewResp.DataViewId, nil
}

func (handler *HandlerBase) BatchUploadAnnotations(dataViewId string, formFiles []aifsclientgo.FormFile, batchSize int) (err error) {
	log.Infof("upload annotation data to data view %s", dataViewId)
	// batch upload annotation files
	count := 0.0
	total := float64(len(formFiles))
	return collection.BatchRange(formFiles, batchSize, func(batch []aifsclientgo.FormFile) error {
		if err = handler.UploadFormFiles(dataViewId, batch); err != nil {
			return err
		}
		count += float64(len(batch))
		log.Infof("upload annotation data %.3f", count/total)
		return nil
	})
}
