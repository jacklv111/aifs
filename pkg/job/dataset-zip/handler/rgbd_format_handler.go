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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flatbuffers "github.com/google/flatbuffers/go"
	aifsclientgo "github.com/jacklv111/aifs-client-go"
	annotempconst "github.com/jacklv111/aifs/pkg/annotation-template-type"
	rgbdvb "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/value-object"
	"github.com/jacklv111/aifs/pkg/job/dataset-zip/constant"
	"github.com/jacklv111/common-sdk/annotation"
	"github.com/jacklv111/common-sdk/collection/mapset"
	"github.com/jacklv111/common-sdk/flatbuffer/raw-data/go/RawData/Rgbd"
	"github.com/jacklv111/common-sdk/utils"
)

type RgbdHandler struct {
	HandlerBase
}

func NewRgbdHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) Handler {
	return &RgbdHandler{
		HandlerBase: HandlerBase{
			client:     client,
			dataView:   dataView,
			datasetDir: datasetDir,
		},
	}
}

func (handler *RgbdHandler) createAnnotationTemplate(labelNameList []string) (annotation.AnnotationTemplateDetails, error) {
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
		annoTempReq.SetType(annotempconst.RGBD)
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

func (handler *RgbdHandler) getLabelNameList(annoDir string) ([]string, error) {
	files, err := utils.ReadAllFiles(annoDir)
	if err != nil {
		return nil, err
	}

	labelNameSet := mapset.NewSet[string]()
	labelList := make([]string, 0)

	for _, file := range files {
		annoFile, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(annoFile)
		for scanner.Scan() {
			line := scanner.Text()
			stList := strings.FieldsFunc(line, func(r rune) bool {
				return r == ' ' || r == '	'
			})
			labelName := stList[0]
			if labelNameSet.Contains(labelName) {
				continue
			}
			labelNameSet.Add(labelName)
			labelList = append(labelList, labelName)
		}
		err = annoFile.Close()
		if err != nil {
			return nil, err
		}
	}
	return labelList, nil
}

func (handler *RgbdHandler) UploadRawData(dataViewId string, inputDir string) (map[string]string, error) {
	dataBinHashMap := make(map[string]string)
	imageDir := filepath.Join(inputDir, "image")
	depthDir := filepath.Join(inputDir, "depth")
	calibDir := filepath.Join(inputDir, "calib")
	imageMap, err := utils.ReadFilesReturnMap(imageDir)
	if err != nil {
		return nil, err
	}
	depthMap, err := utils.ReadFilesReturnMap(depthDir)
	if err != nil {
		return nil, err
	}
	calibMap, err := utils.ReadFilesReturnMap(calibDir)
	if err != nil {
		return nil, err
	}

	if len(imageMap) != len(depthMap) || len(imageMap) != len(calibMap) {
		return nil, fmt.Errorf("invalid raw data, the size of image %d, depth %d, calib %d are not equal", len(imageMap), len(depthMap), len(calibMap))
	}

	count := 0
	var files []aifsclientgo.FormFile
	var metaData *bytes.Buffer
	for fileBaseName, imageFilePath := range imageMap {
		depthFilePath, ok1 := depthMap[fileBaseName]
		calibFilePath, ok2 := calibMap[fileBaseName]
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("image %s exists, but depth ok is %t, calib ok is %t", imageFilePath, ok1, ok2)
		}
		if count == 0 {
			files = make([]aifsclientgo.FormFile, 0)
			metaData = bytes.NewBuffer([]byte{})
		}

		fileName := fileBaseName + ".bin"

		bin, meta, err := handler.GetBinAndMeta(fileName, imageFilePath, depthFilePath, calibFilePath)
		if err != nil {
			return nil, err
		}
		files = append(files, aifsclientgo.FormFile{
			FileName: fileName,
			Value:    bin,
			Key:      "files",
		})
		metaData.WriteString(meta.String() + "\n")
		dataBinHashMap[fileBaseName] = meta.Sha256

		count++
		if count == constant.UPLOAD_RGBD_BIN_BATCH {
			err = handler.uploadRawDataBatch(dataViewId, files, metaData)
			if err != nil {
				return nil, err
			}
			count = 0
		}
	}

	if count != 0 {
		err = handler.uploadRawDataBatch(dataViewId, files, metaData)
		if err != nil {
			return nil, err
		}
	}
	return dataBinHashMap, nil
}

func (handler *RgbdHandler) getAnnoFormFiles(annoDir string, rawDataIdHashMap map[string]string, dataBinHashMap map[string]string, annoTemp annotation.AnnotationTemplateDetails) ([]aifsclientgo.FormFile, error) {
	annoFileMap, _ := utils.ReadFilesReturnMap(annoDir)
	var files []aifsclientgo.FormFile

	for fileName, annoFile := range annoFileMap {
		hash, ok := dataBinHashMap[fileName]
		if !ok {
			return nil, fmt.Errorf("file %s not found in data bin map", fileName)
		}
		rawDataId, ok := rawDataIdHashMap[hash]
		if !ok {
			return nil, fmt.Errorf("file %s not found in raw data id map", fileName)
		}
		file, err := os.Open(annoFile)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)

		annoData := annotation.RgbdAnnotation{
			RawDataId:            rawDataId,
			AnnotationTemplateId: annoTemp.GetAnnoTempId(),
			BoundingBoxList:      make([]annotation.RgbdBoundingBox, 0),
		}
		for scanner.Scan() {
			line := scanner.Text()
			stList := strings.FieldsFunc(line, func(r rune) bool {
				return r == ' ' || r == '	'
			})
			if len(stList) != 13 {
				return nil, fmt.Errorf("13 parameters required per line")
			}
			floatList := make([]float32, 0)
			for i := 1; i < len(stList); i++ {
				data, err := strconv.ParseFloat(stList[i], 32)
				if err != nil {
					return nil, err
				}
				floatList = append(floatList, float32(data))
			}
			labelName := stList[0]
			labelId := annoTemp.GetIdByName(labelName)
			bbox := annotation.RgbdBoundingBox{
				LabelId: labelId,
				BoundingBox2D: annotation.BBox2D{
					X:      floatList[0],
					Y:      floatList[1],
					Width:  floatList[2],
					Height: floatList[3],
				},
				BoundingBox3D: annotation.BBox3D{
					X:     floatList[4],
					Y:     floatList[5],
					Z:     floatList[6],
					XSize: floatList[7],
					YSize: floatList[8],
					ZSize: floatList[9],
					YawX:  floatList[10],
					YawY:  floatList[11],
					YawZ:  0.0,
				},
			}
			annoData.BoundingBoxList = append(annoData.BoundingBoxList, bbox)
		}
		if err = file.Close(); err != nil {
			return nil, err
		}
		annoBytes, err := json.Marshal(annoData)
		if err != nil {
			return nil, err
		}
		files = append(files, aifsclientgo.FormFile{
			FileName: fileName + ".json",
			Value:    annoBytes,
			Key:      rawDataId,
		})
	}

	return files, nil
}

func (handler *RgbdHandler) uploadRawDataBatch(dataViewId string, files []aifsclientgo.FormFile, metaData *bytes.Buffer) error {
	metaForm := aifsclientgo.FormFile{
		FileName: "fileMeta",
		Value:    metaData.Bytes(),
		Key:      "fileMeta",
	}
	resp, err := handler.client.DataViewUploadApi.UploadRawDataToDataView(context.Background(), dataViewId).Files(files).FileMeta(metaForm).Execute()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload raw data status code %d", resp.StatusCode)
	}
	return nil
}

func (bo *RgbdHandler) GetBinAndMeta(fileName string, imgPath string, depthPath string, calibPath string) (bin []byte, meta rgbdvb.RgbdRawDataMeta, err error) {
	meta.FileName = fileName
	// image ext
	imgMeta, err := utils.GetImageMeta(imgPath)
	if err != nil {
		return nil, meta, err
	}
	meta.ImageHeight = int64(imgMeta.Height)
	meta.ImageWidth = int64(imgMeta.Width)
	meta.ImageSize = int64(imgMeta.Size)

	// depth ext
	imgMeta, err = utils.GetImageMeta(depthPath)
	if err != nil {
		return nil, meta, err
	}
	meta.DepthHeight = int64(imgMeta.Height)
	meta.DepthWidth = int64(imgMeta.Width)
	meta.DepthSize = int64(imgMeta.Size)

	// initial size of the buffer (here 1024 bytes), which will grow automatically if needed
	builder := flatbuffers.NewBuilder(1024)
	imageBytes, err := os.ReadFile(imgPath)
	if err != nil {
		return nil, meta, err
	}
	depthBytes, err := os.ReadFile(depthPath)
	if err != nil {
		return nil, meta, err
	}
	imageOffset := builder.CreateByteVector(imageBytes)
	depthOffset := builder.CreateByteVector(depthBytes)

	file, err := os.Open(calibPath)
	if err != nil {
		return nil, meta, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if err != nil {
		return nil, meta, err
	}

	scanner.Split(bufio.ScanWords)
	extrinsics := make([]float64, 9)
	intrinsics := make([]float64, 9)
	for i := 0; i < 9; i++ {
		if !scanner.Scan() {
			return nil, meta, fmt.Errorf("data %s has wrong calib data", calibPath)
		}
		extrinsics[i], err = strconv.ParseFloat(scanner.Text(), 32)

		if err != nil {
			return nil, meta, err
		}
	}
	for i := 0; i < 9; i++ {
		if !scanner.Scan() {
			return nil, meta, fmt.Errorf("data %s has wrong calib data", calibPath)
		}
		intrinsics[i], err = strconv.ParseFloat(scanner.Text(), 32)

		if err != nil {
			return nil, meta, err
		}
	}
	Rgbd.CalibStartExtrinsicsVector(builder, 9)
	for i := 8; i >= 0; i-- {
		builder.PrependFloat32(float32(extrinsics[i]))
	}
	extrinsicsOffset := builder.EndVector(9)
	Rgbd.CalibStartIntrinsicsVector(builder, 9)
	for i := 8; i >= 0; i-- {
		builder.PrependFloat32(float32(intrinsics[i]))
	}
	intrinsicsOffset := builder.EndVector(9)

	Rgbd.CalibStart(builder)
	Rgbd.CalibAddExtrinsics(builder, extrinsicsOffset)
	Rgbd.CalibAddIntrinsics(builder, intrinsicsOffset)
	calibOffset := Rgbd.CalibEnd(builder)

	Rgbd.RgbdDataStart(builder)
	Rgbd.RgbdDataAddImage(builder, imageOffset)
	Rgbd.RgbdDataAddDepth(builder, depthOffset)
	Rgbd.RgbdDataAddCalib(builder, calibOffset)
	rgbdData := Rgbd.RgbdDataEnd(builder)

	builder.Finish(rgbdData)
	buf := builder.FinishedBytes()

	hashStr, err := utils.GetFileSha256Bytes(buf)
	if err != nil {
		return nil, meta, err
	}
	meta.Sha256 = hashStr

	return buf, meta, nil
}

func (handler *RgbdHandler) Exec() error {
	trainDataDir := filepath.Join(handler.datasetDir, "train")
	valDataDir := filepath.Join(handler.datasetDir, "val")
	trainAnnoDir := filepath.Join(trainDataDir, "label")
	valAnnoDir := filepath.Join(valDataDir, "label")

	labelNameList, err := handler.getLabelNameList(trainAnnoDir)
	if err != nil {
		return err
	}
	annoTemp, err := handler.createAnnotationTemplate(labelNameList)
	if err != nil {
		return err
	}
	trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId, err := handler.CreateAllDataView(aifsclientgo.RGBD, annoTemp)
	if err != nil {
		return err
	}

	var dataBinHashMap map[string]string
	dataBinHashMap, err = handler.UploadRawData(trainRawDataViewId, trainDataDir)
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
	formFiles, err := handler.getAnnoFormFiles(trainAnnoDir, trainRawDataIdHashMap, dataBinHashMap, annoTemp)
	if err != nil {
		return err
	}
	err = handler.BatchUploadAnnotations(trainAnnoViewId, formFiles, constant.UPLOAD_RGBD_ANNOTATION_BATCH_SIZE)
	if err != nil {
		return err
	}
	if err := handler.UpdateProgress(70.0, "upload train annotation completed"); err != nil {
		return err
	}
	dataBinHashMap, err = handler.UploadRawData(valRawDataViewId, valDataDir)
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

	formFiles, err = handler.getAnnoFormFiles(valAnnoDir, valRawDataIdHashMap, dataBinHashMap, annoTemp)
	if err != nil {
		return err
	}
	err = handler.BatchUploadAnnotations(valAnnoViewId, formFiles, constant.UPLOAD_RGBD_ANNOTATION_BATCH_SIZE)
	if err != nil {
		return err
	}
	if err := handler.UpdateProgress(100.0, "dataset zip view decompression completed"); err != nil {
		return err
	}
	return nil
}
