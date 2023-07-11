/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	annotempconst "github.com/jacklv111/aifs/pkg/annotation-template-type"
	dzipjobconst "github.com/jacklv111/aifs/pkg/job/dataset-zip/constant"
	"github.com/jacklv111/common-sdk/annotation"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/utils"
)

type ClassificationHandler struct {
	HandlerBase
}

func NewClassificationHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) Handler {
	return &ClassificationHandler{
		HandlerBase: HandlerBase{
			client:     client,
			dataView:   dataView,
			datasetDir: datasetDir,
		},
	}
}

func (handler *ClassificationHandler) GetLabelName(filePath string) string {
	spList := strings.Split(filePath, "/")
	labelName := spList[len(spList)-2]
	return labelName
}

func (handler *ClassificationHandler) getLabelNameList(filePaths []string) []string {
	labelNameMap := make(map[string]string, 0)
	for _, filePath := range filePaths {
		labelName := handler.GetLabelName(filePath)
		labelNameMap[labelName] = ""
	}
	labelNameList := make([]string, 0)
	for name := range labelNameMap {
		labelNameList = append(labelNameList, name)
	}
	return labelNameList
}

func (handler *ClassificationHandler) createAnnotationTemplate(labelNameList []string) (annotation.AnnotationTemplateDetails, error) {
	if handler.dataView.AnnotationTemplateId == nil {
		var annoTempReq aifsclientgo.CreateAnnotationTemplateRequest
		labelList := make([]aifsclientgo.Label, 0)
		for _, name := range labelNameList {
			label := aifsclientgo.Label{
				Name: name,
			}
			labelList = append(labelList, label)
		}
		annoTempReq.SetLabels(labelList)
		annoTempReq.SetName(handler.dataView.GetName())
		annoTempReq.SetType(annotempconst.CATEGORY)
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

func (handler *ClassificationHandler) uploadAnnoData(dataViewId string, dataPathList []string, rawDataIdHashMap map[string]string, annoTemp annotation.AnnotationTemplateDetails) error {
	return collection.BatchRange(dataPathList, dzipjobconst.UPLOAD_CLASSIFICATION_ANNOTATION_BATCH_SIZE, func(batch []string) error {
		var writer bytes.Buffer
		for _, filePath := range batch {
			labelName := handler.GetLabelName(filePath)
			hash, err := utils.GetFileSha256FromFile(filePath)
			if err != nil {
				return err
			}
			rawDataId, ok := rawDataIdHashMap[hash]
			if ok {
				writer.WriteString(fmt.Sprintf("%s %s\n", rawDataId, annoTemp.GetIdByName(labelName)))
			}
		}
		resp, err := handler.client.DataViewUploadApi.UploadAnnotationToDataView(context.Background(), dataViewId).FileMeta(writer.Bytes()).Execute()
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("upload raw data status code %d", resp.StatusCode)
		}
		return nil
	})
}

func (handler *ClassificationHandler) Exec() error {
	trainDataDir := filepath.Join(handler.datasetDir, "train")
	valDataDir := filepath.Join(handler.datasetDir, "val")

	trainDataPathList, err := utils.ReadAllFiles(trainDataDir)
	if err != nil {
		return err
	}
	valDataPathList, err := utils.ReadAllFiles(valDataDir)
	if err != nil {
		return err
	}

	annoTemp, err := handler.createAnnotationTemplate(handler.getLabelNameList(trainDataPathList))
	if err != nil {
		return err
	}
	trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId, err := handler.CreateAllDataView(aifsclientgo.IMAGE, annoTemp)
	if err != nil {
		return err
	}

	err = handler.UploadRawData(trainRawDataViewId, trainDataPathList)
	if err != nil {
		return err
	}
	if err := handler.UpdateProgress(50.0, "upload train raw data completed"); err != nil {
		return err
	}
	trainRawDataIdHashMap, err := handler.GetRawDataHashIdMap(trainRawDataViewId)
	if err != nil {
		return err
	}
	err = handler.uploadAnnoData(trainAnnoViewId, trainDataPathList, trainRawDataIdHashMap, annoTemp)
	if err != nil {
		return err
	}
	if err := handler.UpdateProgress(70.0, "upload train annotation completed"); err != nil {
		return err
	}
	err = handler.UploadRawData(valRawDataViewId, valDataPathList)
	if err != nil {
		return err
	}
	if err := handler.UpdateProgress(80.0, "upload valid raw data completed"); err != nil {
		return err
	}
	valRawDataIdHashMap, err := handler.GetRawDataHashIdMap(valRawDataViewId)
	if err != nil {
		return err
	}
	err = handler.uploadAnnoData(valAnnoViewId, valDataPathList, valRawDataIdHashMap, annoTemp)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(100.0, "dataset zip view decompression completed"); err != nil {
		return err
	}
	return nil
}
