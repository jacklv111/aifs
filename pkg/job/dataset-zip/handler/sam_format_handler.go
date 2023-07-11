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
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	annotempconst "github.com/jacklv111/aifs/pkg/annotation-template-type"
	"github.com/jacklv111/aifs/pkg/job/dataset-zip/constant"
	"github.com/jacklv111/common-sdk/annotation"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/utils"
)

type SamHandler struct {
	HandlerBase
}

func NewSamHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) Handler {
	return &SamHandler{
		HandlerBase: HandlerBase{
			client:     client,
			dataView:   dataView,
			datasetDir: datasetDir,
		},
	}
}

func splitData(dataPathList []string) (rawData map[string]string, anno map[string]string) {
	rawData = make(map[string]string)
	anno = make(map[string]string)

	for _, dataPath := range dataPathList {
		fileName := utils.GetFileNameWithoutSuffix(filepath.Base(dataPath))
		fileExt := filepath.Ext(dataPath)
		if fileExt == ".json" {
			anno[fileName] = dataPath
		} else {
			rawData[fileName] = dataPath
		}
	}
	return
}

func getValueList(dataMap map[string]string) []string {
	res := make([]string, 0)
	for _, v := range dataMap {
		res = append(res, v)
	}
	return res
}

func (handler *SamHandler) getAnnotationTemplate() (annotation.AnnotationTemplateDetails, error) {
	if handler.dataView.AnnotationTemplateId == nil {
		var annoTempReq aifsclientgo.CreateAnnotationTemplateRequest
		annoTempReq.SetName(handler.dataView.GetName())
		annoTempReq.SetType(annotempconst.COCO_TYPE)
		annoTempReq.SetDescription(fmt.Sprintf("extract from zip %s", handler.dataView.GetName()))

		createResp, _, err := handler.client.AnnotationTemplateApi.CreateAnnotationTemplate(context.Background()).CreateAnnotationTemplateRequest(annoTempReq).Execute()
		if err != nil {
			return annotation.AnnotationTemplateDetails{}, err
		}

		dataViewId := *handler.dataView.Id
		var updateReq aifsclientgo.UpdateDatasetZipRequest
		updateReq.AnnotationTemplateId = createResp.AnnotationTemplateId
		resp, err := handler.client.DataViewApi.UpdateDatasetZipView(context.Background(), dataViewId).UpdateDatasetZipRequest(updateReq).Execute()
		if err != nil {
			return annotation.AnnotationTemplateDetails{}, err
		}
		if resp.StatusCode != http.StatusOK {
			return annotation.AnnotationTemplateDetails{}, fmt.Errorf("update dataset zip view got incorrect status code %d", resp.StatusCode)
		}

		handler.dataView.AnnotationTemplateId = createResp.AnnotationTemplateId
	}

	annoTempId := *handler.dataView.AnnotationTemplateId
	annoTempDetails, _, err := handler.client.AnnotationTemplateApi.GetAnnoTemplateDetails(context.Background(), annoTempId).Execute()
	return annotation.NewAnnotationTemplateDetails(*annoTempDetails), err
}

func (handler *SamHandler) getAnnoFormFiles(rawDataMap, annoDataMap map[string]string, rawDataIdHashMap map[string]string, annoTemp annotation.AnnotationTemplateDetails) ([]aifsclientgo.FormFile, error) {
	log.Info("get annotation form files")
	annoList := make([]*annotation.CocoAnnoFormat, 0)

	for fileName, filePath := range annoDataMap {
		annoSam := &annotation.CocoAnno{}
		annoSam.ParseAnnotationFile(filePath)

		imageHash, err := utils.GetFileSha256FromFile(rawDataMap[fileName])
		if err != nil {
			return nil, err
		}
		rawDataId, ok := rawDataIdHashMap[imageHash]
		if !ok {
			log.Errorf("annotation %s doesn't have raw data", filePath)
			continue
		}

		anno := &annotation.CocoAnnoFormat{}
		anno.RawDataId = rawDataId
		anno.AnnotationTemplateId = *annoTemp.Id
		for _, data := range annoSam.Annotations {
			var cocoAnnoData annotation.CocoAnnoData
			cocoAnnoData.Area = data.Area
			cocoAnnoData.Bbox = data.Bbox
			cocoAnnoData.Id = data.Id
			cocoAnnoData.IsCrowd = data.IsCrowd
			cocoAnnoData.KeyPoints = data.KeyPoints
			cocoAnnoData.NumKeyPoints = data.NumKeyPoints
			cocoAnnoData.Segmentation = data.Segmentation
			cocoAnnoData.PredictedIou = data.PredictedIou
			cocoAnnoData.PointCoords = data.PointCoords
			cocoAnnoData.CropBox = data.CropBox
			cocoAnnoData.StabilityScore = data.StabilityScore
			anno.AnnoData = append(anno.AnnoData, cocoAnnoData)
		}
		annoList = append(annoList, anno)
	}

	formFiles := make([]aifsclientgo.FormFile, 0)
	log.Infof("total %d annotations", len(annoList))
	for _, data := range annoList {
		content, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		formFiles = append(formFiles,
			aifsclientgo.FormFile{
				Key:      data.RawDataId,
				FileName: data.RawDataId + ".json",
				Value:    content,
			},
		)
	}
	return formFiles, nil
}

func (handler *SamHandler) Exec() error {
	dataPathList, err := utils.ReadAllFiles(handler.datasetDir)
	if err != nil {
		return err
	}

	rawDataMap, annoMap := splitData(dataPathList)

	annoTemp, err := handler.getAnnotationTemplate()
	if err != nil {
		return err
	}

	rawDataViewId, annotationViewId, err := handler.CreateRawDataAndAnnotationView(aifsclientgo.IMAGE, annoTemp)
	if err != nil {
		return err
	}

	err = handler.UploadRawData(rawDataViewId, getValueList(rawDataMap))
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(50.0, "upload train raw data completed"); err != nil {
		return err
	}

	rawDataIdHashMap, err := handler.GetRawDataHashIdMap(rawDataViewId)
	if err != nil {
		return err
	}

	formFiles, err := handler.getAnnoFormFiles(rawDataMap, annoMap, rawDataIdHashMap, annoTemp)
	if err != nil {
		return err
	}
	err = handler.BatchUploadAnnotations(annotationViewId, formFiles, constant.UPLOAD_COCO_ANNOTATION_BATCH_SIZE)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(100.0, "dataset zip view decompression completed"); err != nil {
		return err
	}
	return nil
}
