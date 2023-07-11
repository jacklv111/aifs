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

type CocoHandler struct {
	HandlerBase
}

func NewCocoHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) Handler {
	return &CocoHandler{
		HandlerBase: HandlerBase{
			client:     client,
			dataView:   dataView,
			datasetDir: datasetDir,
		},
	}
}

func (handler *CocoHandler) createAnnotationTemplate(categories []annotation.CocoCategory) (annotation.AnnotationTemplateDetails, error) {
	if handler.dataView.AnnotationTemplateId == nil {
		var annoTempReq aifsclientgo.CreateAnnotationTemplateRequest
		labelList := make([]aifsclientgo.Label, 0)
		for _, category := range categories {
			label := aifsclientgo.Label{
				Name:              category.Name,
				SuperCategoryName: &category.SuperCategory,
				KeyPointDef:       category.KeyPoints,
				KeyPointSkeleton:  category.Skeleton,
			}
			labelList = append(labelList, label)
		}
		annoTempReq.SetLabels(labelList)
		annoTempReq.SetName(handler.dataView.GetName())
		annoTempReq.SetType(annotempconst.COCO_TYPE)
		annoTempReq.SetDescription(fmt.Sprintf("extract from zip %s", handler.dataView.GetName()))

		createResp, _, err := handler.client.AnnotationTemplateApi.CreateAnnotationTemplate(context.Background()).CreateAnnotationTemplateRequest(annoTempReq).Execute()
		if err != nil {
			return annotation.AnnotationTemplateDetails{}, err
		}
		handler.dataView.AnnotationTemplateId = createResp.AnnotationTemplateId

		var updateReq aifsclientgo.UpdateDatasetZipRequest
		var updateResp *http.Response
		updateReq.SetAnnotationTemplateId(*createResp.AnnotationTemplateId)
		updateResp, err = handler.client.DataViewApi.UpdateDatasetZipView(context.Background(), *handler.dataView.Id).UpdateDatasetZipRequest(updateReq).Execute()
		if err != nil {
			return annotation.AnnotationTemplateDetails{}, err
		}
		if updateResp.StatusCode != http.StatusOK {
			return annotation.AnnotationTemplateDetails{}, fmt.Errorf("update dataset zip view got incorrect status code %d", updateResp.StatusCode)
		}
	}

	annoTempId := *handler.dataView.AnnotationTemplateId
	annoTempDetails, _, err := handler.client.AnnotationTemplateApi.GetAnnoTemplateDetails(context.Background(), annoTempId).Execute()
	return annotation.NewAnnotationTemplateDetails(*annoTempDetails), err
}

func (handler *CocoHandler) getAnnoFormFiles(annos *annotation.CocoAnno, rawDataDir string, rawDataIdHashMap map[string]string, annoTemp annotation.AnnotationTemplateDetails) ([]aifsclientgo.FormFile, error) {
	log.Info("get annotation form files")
	annoMap := make(map[int]*annotation.CocoAnnoFormat)

	for _, data := range annos.Images {
		filePath := filepath.Join(rawDataDir, data.FileName)
		hash, err := utils.GetFileSha256FromFile(filePath)
		if err != nil {
			return nil, err
		}

		var cocoRawDataAnno annotation.CocoAnnoFormat
		cocoRawDataAnno.AnnotationTemplateId = annoTemp.GetAnnoTempId()
		if rawDataIdHashMap[hash] == "" {
			return nil, fmt.Errorf("raw data id not found for raw data %s", data.FileName)
		}
		cocoRawDataAnno.RawDataId = rawDataIdHashMap[hash]
		cocoRawDataAnno.AnnoData = make([]annotation.CocoAnnoData, 0)
		annoMap[data.Id] = &cocoRawDataAnno
	}

	for _, data := range annos.Annotations {
		var cocoAnnoData annotation.CocoAnnoData
		cocoAnnoData.Area = data.Area
		cocoAnnoData.Bbox = data.Bbox
		cocoAnnoData.Id = data.Id
		cocoAnnoData.IsCrowd = data.IsCrowd
		cocoAnnoData.KeyPoints = data.KeyPoints
		cocoAnnoData.NumKeyPoints = data.NumKeyPoints
		cocoAnnoData.Segmentation = data.Segmentation
		cocoAnnoData.LabelId = annoTemp.GetIdByName(annos.CategoryMap[data.CategoryId])
		if cocoAnnoData.LabelId == "" {
			return nil, fmt.Errorf("label id not found for label name %s", annos.CategoryMap[data.CategoryId])
		}
		annoMap[data.ImageId].AnnoData = append(annoMap[data.ImageId].AnnoData, cocoAnnoData)
	}

	formFiles := make([]aifsclientgo.FormFile, 0)
	for _, data := range annoMap {
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

func (handler *CocoHandler) Exec() error {
	trainDataDir := filepath.Join(handler.datasetDir, "train")
	trainRawDataDir := filepath.Join(trainDataDir, "raw-data")
	trainAnnoFile := filepath.Join(trainDataDir, "annotations.json")

	valDataDir := filepath.Join(handler.datasetDir, "val")
	valRawDataDir := filepath.Join(valDataDir, "raw-data")
	valAnnoFile := filepath.Join(valDataDir, "annotations.json")

	annos := &annotation.CocoAnno{}

	err := annos.ParseAnnotationFile(trainAnnoFile)
	if err != nil {
		return err
	}

	annoTemp, err := handler.createAnnotationTemplate(annos.Categories)
	if err != nil {
		return err
	}

	trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId, err := handler.CreateAllDataView(aifsclientgo.IMAGE, annoTemp)
	if err != nil {
		return err
	}
	// train
	dataPathList, err := utils.ReadAllFiles(trainRawDataDir)
	if err != nil {
		return err
	}
	err = handler.UploadRawData(trainRawDataViewId, dataPathList)
	if err != nil {
		return err
	}
	if err := handler.UpdateProgress(50.0, "upload train raw data completed"); err != nil {
		return err
	}
	rawDataIdHashMap, err := handler.GetRawDataHashIdMap(trainRawDataViewId)
	if err != nil {
		return err
	}
	formFiles, err := handler.getAnnoFormFiles(annos, trainRawDataDir, rawDataIdHashMap, annoTemp)
	if err != nil {
		return err
	}
	err = handler.BatchUploadAnnotations(trainAnnoViewId, formFiles, constant.UPLOAD_COCO_ANNOTATION_BATCH_SIZE)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(70.0, "upload train annotation completed"); err != nil {
		return err
	}
	// val
	dataPathList, err = utils.ReadAllFiles(valRawDataDir)
	if err != nil {
		return err
	}
	err = handler.UploadRawData(valRawDataViewId, dataPathList)
	if err != nil {
		return err
	}
	if err := handler.UpdateProgress(80.0, "upload valid raw data completed"); err != nil {
		return err
	}
	rawDataIdHashMap, err = handler.GetRawDataHashIdMap(valRawDataViewId)
	if err != nil {
		return err
	}

	err = annos.ParseAnnotationFile(valAnnoFile)
	if err != nil {
		return err
	}

	formFiles, err = handler.getAnnoFormFiles(annos, valRawDataDir, rawDataIdHashMap, annoTemp)
	if err != nil {
		return err
	}
	err = handler.BatchUploadAnnotations(valAnnoViewId, formFiles, constant.UPLOAD_COCO_ANNOTATION_BATCH_SIZE)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(100.0, "dataset zip view decompression completed"); err != nil {
		return err
	}
	return nil
}
