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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	annotempconst "github.com/jacklv111/aifs/pkg/annotation-template-type"
	dzipjobconst "github.com/jacklv111/aifs/pkg/job/dataset-zip/constant"
	"github.com/jacklv111/common-sdk/annotation"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/utils"
)

type OcrHandler struct {
	HandlerBase
}

func NewOcrHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) Handler {
	return &OcrHandler{
		HandlerBase: HandlerBase{
			client:     client,
			dataView:   dataView,
			datasetDir: datasetDir,
		},
	}
}

func (handler *OcrHandler) Exec() error {
	trainDataDir := filepath.Join(handler.datasetDir, "train")
	trainRawDataDir := filepath.Join(trainDataDir, "images")
	trainAnnoDataFile := filepath.Join(trainDataDir, "annotations.txt")

	valDataDir := filepath.Join(handler.datasetDir, "val")
	valRawDataDir := filepath.Join(valDataDir, "images")
	valAnnoDataFile := filepath.Join(valDataDir, "annotations.txt")

	wordsFile := filepath.Join(handler.datasetDir, "words.txt")

	annoTemp, err := handler.createAnnotationTemplate(wordsFile)
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
	rawDataIdHashMap, err := handler.GetRawDataHashIdMap(trainRawDataViewId)
	if err != nil {
		return err
	}
	rawDataFileAnnoMap, err := handler.parseAnnotationFile(trainAnnoDataFile)
	if err != nil {
		return err
	}
	err = handler.uploadAnnoData(trainAnnoViewId, dataPathList, rawDataFileAnnoMap, rawDataIdHashMap, annoTemp)
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
	rawDataIdHashMap, err = handler.GetRawDataHashIdMap(valRawDataViewId)
	if err != nil {
		return err
	}
	rawDataFileAnnoMap, err = handler.parseAnnotationFile(valAnnoDataFile)
	if err != nil {
		return err
	}
	err = handler.uploadAnnoData(valAnnoViewId, dataPathList, rawDataFileAnnoMap, rawDataIdHashMap, annoTemp)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(100.0, "dataset zip view decompression completed"); err != nil {
		return err
	}
	return nil
}

func (handler *OcrHandler) uploadAnnoData(dataViewId string, dataPathList []string, rawDataFileAnnoMap map[string]string, rawDataIdHashMap map[string]string, annoTemp annotation.AnnotationTemplateDetails) error {
	log.Info("upload annotation data")
	count := 0
	return collection.BatchRange(dataPathList, dzipjobconst.UPLOAD_OCR_ANNOTATION_BATCH_SIZE, func(batch []string) error {
		count += len(batch)
		log.Infof("upload annotation data %d/%d", count, len(dataPathList))

		var writer bytes.Buffer
		for _, filePath := range batch {
			if _, ok := rawDataFileAnnoMap[filepath.Base(filePath)]; !ok {
				log.Warnf("file %s not found in annotation file", filePath)
				continue
			}
			hash, err := utils.GetFileSha256FromFile(filePath)
			if err != nil {
				return err
			}
			rawDataId, ok := rawDataIdHashMap[hash]
			if ok {
				writer.WriteString(fmt.Sprintf("%s %s\n", rawDataId, rawDataFileAnnoMap[filepath.Base(filePath)]))
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

func (handler *OcrHandler) createAnnotationTemplate(wordsFile string) (annotation.AnnotationTemplateDetails, error) {
	if handler.dataView.AnnotationTemplateId == nil {
		log.Info("create annotation template")

		wordList := make([]string, 0)
		file, err := os.Open(wordsFile)
		if err != nil {
			return annotation.AnnotationTemplateDetails{}, err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			wordList = append(wordList, scanner.Text())
		}

		var annoTempReq aifsclientgo.CreateAnnotationTemplateRequest
		annoTempReq.SetName(handler.dataView.GetName())
		annoTempReq.SetType(annotempconst.OCR)
		annoTempReq.SetWordList(wordList)
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
	} else {
		log.Infof("annotation template %s already exist, skip", *handler.dataView.AnnotationTemplateId)
	}

	annoTempId := *handler.dataView.AnnotationTemplateId
	annoTempDetails, _, err := handler.client.AnnotationTemplateApi.GetAnnoTemplateDetails(context.Background(), annoTempId).Execute()
	return annotation.NewAnnotationTemplateDetails(*annoTempDetails), err
}

// key is raw data file name, value is words in the raw data
func (handler *OcrHandler) parseAnnotationFile(annoFilePath string) (map[string]string, error) {
	log.Infof("parse annotation file %s", annoFilePath)
	annoMap := make(map[string]string)
	file, err := os.Open(annoFilePath)
	if err != nil {
		return annoMap, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == '\t' || r == '	'
		})
		if len(split) != 2 {
			return annoMap, fmt.Errorf("incorrect annotation file format")
		}
		annoMap[split[0]] = split[1]
	}
	return annoMap, nil
}
