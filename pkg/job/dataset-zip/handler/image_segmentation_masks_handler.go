/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	annotempconst "github.com/jacklv111/aifs/pkg/annotation-template-type"
	"github.com/jacklv111/aifs/pkg/job/dataset-zip/constant"
	"github.com/jacklv111/common-sdk/annotation"
	"github.com/jacklv111/common-sdk/utils"
)

type annotationTemplate struct {
	Name   string  `json:"name"`
	Labels []label `json:"labels"`
}

type label struct {
	Name  string `json:"name"`
	Color int32  `json:"color"`
}

type ImageSegmentationMasksHandler struct {
	HandlerBase
}

func NewImageSegmentationMasksHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) Handler {
	return &ImageSegmentationMasksHandler{
		HandlerBase: HandlerBase{
			client:     client,
			dataView:   dataView,
			datasetDir: datasetDir,
		},
	}
}

func (handler *ImageSegmentationMasksHandler) Exec() error {
	trainDataDir := filepath.Join(handler.datasetDir, "train")
	trainRawDataDir := filepath.Join(trainDataDir, "raw-data")
	trainAnnoDir := filepath.Join(trainDataDir, "annotation")
	trainAnnoFile := filepath.Join(trainDataDir, "anno.txt")

	valDataDir := filepath.Join(handler.datasetDir, "val")
	valRawDataDir := filepath.Join(valDataDir, "raw-data")
	valAnnoDir := filepath.Join(valDataDir, "annotation")
	valAnnoFile := filepath.Join(valDataDir, "anno.txt")

	annotationTemplateFile := filepath.Join(handler.datasetDir, "annotation_template.json")

	rawDataAnnoFileMap, err := handler.parseAnnotationFile(trainAnnoFile)
	if err != nil {
		return err
	}

	annoTemp, err := handler.createAnnotationTemplate(annotationTemplateFile)
	if err != nil {
		return err
	}

	trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId, err := handler.CreateAllDataView(aifsclientgo.IMAGE, annoTemp)
	if err != nil {
		return err
	}

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
	trainRawDataIdHashMap, err := handler.GetRawDataHashIdMap(trainRawDataViewId)
	if err != nil {
		return err
	}
	formFiles, err := handler.getAnnoFormFiles(trainRawDataDir, trainAnnoDir, trainRawDataIdHashMap, rawDataAnnoFileMap, annoTemp)
	if err != nil {
		return err
	}
	err = handler.BatchUploadAnnotations(trainAnnoViewId, formFiles, constant.UPLOAD_IMAGE_SEGMENTATION_MASKS_ANNOTATION_BATCH_SIZE)
	if err != nil {
		return err
	}
	if err := handler.UpdateProgress(70.0, "upload train annotation completed"); err != nil {
		return err
	}

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
	valRawDataIdHashMap, err := handler.GetRawDataHashIdMap(valRawDataViewId)
	if err != nil {
		return err
	}

	rawDataAnnoFileMap, err = handler.parseAnnotationFile(valAnnoFile)
	if err != nil {
		return err
	}

	formFiles, err = handler.getAnnoFormFiles(valRawDataDir, valAnnoDir, valRawDataIdHashMap, rawDataAnnoFileMap, annoTemp)
	if err != nil {
		return err
	}
	err = handler.BatchUploadAnnotations(valAnnoViewId, formFiles, constant.UPLOAD_IMAGE_SEGMENTATION_MASKS_ANNOTATION_BATCH_SIZE)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(100.0, "dataset zip view decompression completed"); err != nil {
		return err
	}
	return nil
}

func (handler *ImageSegmentationMasksHandler) getAnnoFormFiles(rawDataDir string, annoDir string, rawDataIdHashMap map[string]string, rawDataAnooFileMap map[string]string, annoTemp annotation.AnnotationTemplateDetails) ([]aifsclientgo.FormFile, error) {
	// batch upload annotation files
	var formFiles []aifsclientgo.FormFile
	for rawDataFileName, annoFileName := range rawDataAnooFileMap {
		annoFile, err := os.Open(filepath.Join(annoDir, annoFileName))
		if err != nil {
			return nil, err
		}
		annoBytes, err := io.ReadAll(annoFile)
		if err != nil {
			return nil, err
		}
		err = annoFile.Close()
		if err != nil {
			return nil, err
		}
		rawDataHash, err := utils.GetFileSha256FromFile(filepath.Join(rawDataDir, rawDataFileName))
		if err != nil {
			return nil, err
		}
		rawDataId, ok := rawDataIdHashMap[rawDataHash]
		if !ok {
			return nil, fmt.Errorf("raw data id not found, raw data file name: %s", rawDataFileName)
		}
		formFiles = append(formFiles,
			aifsclientgo.FormFile{
				Key:      rawDataId,
				FileName: annoFileName,
				Value:    annoBytes,
			},
		)

	}

	return formFiles, nil
}

func (handler *ImageSegmentationMasksHandler) parseAnnotationFile(annoFilePath string) (map[string]string, error) {
	annoFile, err := os.Open(annoFilePath)
	if err != nil {
		return nil, err
	}
	defer annoFile.Close()
	scanner := bufio.NewScanner(annoFile)
	// key: raw data file name; value: annotation file name
	rawDataAnnoFileMap := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " \t\n\r")
		items := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '\r'
		})
		rawDataAnnoFileMap[items[0]] = items[1]
	}
	return rawDataAnnoFileMap, nil
}

func (handler *ImageSegmentationMasksHandler) createAnnotationTemplate(annoTempFile string) (annotation.AnnotationTemplateDetails, error) {
	var annoTemp annotationTemplate
	content, err := os.ReadFile(annoTempFile)
	if err != nil {
		return annotation.AnnotationTemplateDetails{}, err
	}
	err = json.Unmarshal(content, &annoTemp)
	if err != nil {
		return annotation.AnnotationTemplateDetails{}, err
	}

	if handler.dataView.AnnotationTemplateId == nil {
		var annoTempReq aifsclientgo.CreateAnnotationTemplateRequest
		labelList := make([]aifsclientgo.Label, 0)
		for _, label := range annoTemp.Labels {
			label := aifsclientgo.Label{
				Name:  label.Name,
				Color: label.Color,
			}
			labelList = append(labelList, label)
		}
		annoTempReq.SetLabels(labelList)
		annoTempReq.SetName(handler.dataView.GetName())
		annoTempReq.SetType(annotempconst.SEGMENTATION_MASKS)
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
