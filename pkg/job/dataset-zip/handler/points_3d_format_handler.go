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
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flatbuffers "github.com/google/flatbuffers/go"
	aifsclientgo "github.com/jacklv111/aifs-client-go"
	annotempconst "github.com/jacklv111/aifs/pkg/annotation-template-type"
	p3dvb "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/value-object"
	"github.com/jacklv111/aifs/pkg/job/dataset-zip/constant"
	"github.com/jacklv111/common-sdk/annotation"
	"github.com/jacklv111/common-sdk/flatbuffer/raw-data/go/RawData/Points3D"
	"github.com/jacklv111/common-sdk/utils"
)

type Points3DHandler struct {
	HandlerBase
}

func NewPoints3DHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) Handler {
	return &Points3DHandler{
		HandlerBase: HandlerBase{
			client:     client,
			dataView:   dataView,
			datasetDir: datasetDir,
		},
	}
}

func (handler *Points3DHandler) createAnnotationTemplate(annoTempFile string) (annotation.AnnotationTemplateDetails, error) {
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
		annoTempReq.SetType(annotempconst.POINTS_3D)
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

func (handler *Points3DHandler) UploadRawData(dataViewId string, inputDir string) (map[string]string, error) {
	dataBinHashMap := make(map[string]string)
	posDir := filepath.Join(inputDir, "pos")
	rgbDir := filepath.Join(inputDir, "rgb")

	posMap, err := utils.ReadFilesReturnMap(posDir)
	if err != nil {
		return nil, err
	}

	// rgb map may be empty
	rgbMap, err := utils.ReadFilesReturnMap(rgbDir)
	if err != nil {
		return nil, err
	}

	count := 0
	var files []aifsclientgo.FormFile
	var metaData *bytes.Buffer
	for fileBaseName, posFilePath := range posMap {
		rgbFilePath := rgbMap[fileBaseName]

		if count == 0 {
			files = make([]aifsclientgo.FormFile, 0)
			metaData = bytes.NewBuffer([]byte{})
		}

		fileName := fileBaseName + ".bin"

		bin, meta, err := handler.GetBinAndMeta(fileName, posFilePath, rgbFilePath)
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
		if count == constant.UPLOAD_POINTS_3D_BIN_BATCH {
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

func (handler *Points3DHandler) UploadAnnoFiles(dataViewId, annoDir string, rawDataIdHashMap map[string]string, dataBinHashMap map[string]string, annoTemp annotation.AnnotationTemplateDetails) error {
	annoFileMap, _ := utils.ReadFilesReturnMap(annoDir)
	var files []aifsclientgo.FormFile
	count := 0
	for fileName, annoFile := range annoFileMap {
		if count == 0 {
			files = make([]aifsclientgo.FormFile, 0)
		}
		hash, ok := dataBinHashMap[fileName]
		if !ok {
			return fmt.Errorf("file %s not found in data bin map", fileName)
		}
		rawDataId, ok := rawDataIdHashMap[hash]
		if !ok {
			return fmt.Errorf("file %s not found in raw data id map", fileName)
		}
		file, err := os.Open(annoFile)
		if err != nil {
			return err
		}

		annoBuf := bytes.NewBuffer([]byte{})
		annoWriter := zlib.NewWriter(annoBuf)

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			annoWriter.Write([]byte(annoTemp.GetIdByName(line) + "\n"))
		}
		err = annoWriter.Close()
		if err != nil {
			return err
		}

		if err = file.Close(); err != nil {
			return err
		}

		files = append(files, aifsclientgo.FormFile{
			FileName: fileName,
			Value:    annoBuf.Bytes(),
			Key:      rawDataId,
		})
		count++
		if count == constant.UPLOAD_POINTS_3D_ANNOTATION_BATCH_SIZE {
			err = handler.BatchUploadAnnotations(dataViewId, files, constant.UPLOAD_POINTS_3D_ANNOTATION_BATCH_SIZE)
			if err != nil {
				return err
			}
			count = 0
		}
	}
	if count != 0 {
		err := handler.BatchUploadAnnotations(dataViewId, files, constant.UPLOAD_POINTS_3D_ANNOTATION_BATCH_SIZE)
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler *Points3DHandler) uploadRawDataBatch(dataViewId string, files []aifsclientgo.FormFile, metaData *bytes.Buffer) error {
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

func getPoints(filePath string) (points [][]float64, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stList := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == '	'
		})
		if len(stList) != 3 {
			return nil, fmt.Errorf("point position must have 3 parameters required per line")
		}
		x, err := strconv.ParseFloat(stList[0], 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(stList[1], 64)
		if err != nil {
			return nil, err
		}
		z, err := strconv.ParseFloat(stList[2], 64)
		if err != nil {
			return nil, err
		}
		points = append(points, []float64{x, y, z})
	}
	return points, nil
}

func getRgbList(filePath string) (rgbList [][]float64, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stList := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == '	'
		})
		if len(stList) != 3 {
			return nil, fmt.Errorf("point rgb must have 3 parameters required per line")
		}
		r, err := strconv.ParseFloat(stList[0], 64)
		if err != nil {
			return nil, err
		}
		g, err := strconv.ParseFloat(stList[1], 64)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseFloat(stList[2], 64)
		if err != nil {
			return nil, err
		}
		rgbList = append(rgbList, []float64{r, g, b})
	}
	return rgbList, nil
}

func getPointsMeta(points [][]float64) (xmin, ymin, zmin, xmax, ymax, zmax float64, err error) {
	xmin = math.MaxFloat64
	ymin = math.MaxFloat64
	zmin = math.MaxFloat64
	xmax = -math.MaxFloat64
	ymax = -math.MaxFloat64
	zmax = -math.MaxFloat64

	for _, point := range points {
		x, y, z := point[0], point[1], point[2]
		if x < xmin {
			xmin = x
		}
		if x > xmax {
			xmax = x
		}
		if y < ymin {
			ymin = y
		}
		if y > ymax {
			ymax = y
		}
		if z < zmin {
			zmin = z
		}
		if z > zmax {
			zmax = z
		}
	}
	return xmin, ymin, zmin, xmax, ymax, zmax, nil
}

func getRgbMeta(rgbList [][]float64) (rmean, gmean, bmean, rstd, gstd, bstd float64, err error) {
	if len(rgbList) == 0 {
		return -1, -1, -1, -1, -1, -1, nil
	}

	rchannel := 0.0
	gchannel := 0.0
	bchannel := 0.0
	rchannelSquare := 0.0
	gchannelSquare := 0.0
	bchannelSquare := 0.0
	count := 0.0

	for _, rgb := range rgbList {
		r, g, b := float64(rgb[0]), float64(rgb[1]), float64(rgb[2])

		rchannel += r
		gchannel += g
		bchannel += b
		rchannelSquare += r * r
		gchannelSquare += g * g
		bchannelSquare += b * b
		count++
	}
	if count == 0 {
		return -1, -1, -1, -1, -1, -1, nil
	}
	rmean = rchannel / count
	gmean = gchannel / count
	bmean = bchannel / count
	rstd = math.Sqrt(rchannelSquare/count - rmean*rmean)
	gstd = math.Sqrt(gchannelSquare/count - gmean*gmean)
	bstd = math.Sqrt(bchannelSquare/count - bmean*bmean)
	return rmean, gmean, bmean, rstd, gstd, bstd, nil
}

func (bo *Points3DHandler) GetBinAndMeta(fileName string, posPath string, rgbPath string) (bin []byte, meta p3dvb.Points3DMeta, err error) {
	// initial size of the buffer (here 1024 bytes), which will grow automatically if needed
	builder := flatbuffers.NewBuilder(1024)

	meta.FileName = fileName
	points, err := getPoints(posPath)
	if err != nil {
		return nil, meta, err
	}
	meta.Size = int64(len(points))
	meta.Xmin, meta.Ymin, meta.Zmin, meta.Xmax, meta.Ymax, meta.Zmax, err = getPointsMeta(points)
	if err != nil {
		return nil, meta, err
	}

	rgbList, err := getRgbList(rgbPath)
	if err != nil {
		return nil, meta, err
	}

	meta.Rmean, meta.Gmean, meta.Bmean, meta.Rstd, meta.Gstd, meta.Bstd, err = getRgbMeta(rgbList)
	if err != nil {
		return nil, meta, err
	}

	// put data into flatbuffer
	Points3D.Points3DStartPosVector(builder, len(points))
	for i := len(points) - 1; i >= 0; i-- {
		Points3D.CreatePoint3(builder, float32(points[i][0]), float32(points[i][1]), float32(points[i][2]))
	}
	posData := builder.EndVector(len(points))

	// rgb data can be empty
	Points3D.Points3DStartRgbVector(builder, len(rgbList))
	for i := len(rgbList) - 1; i >= 0; i-- {
		Points3D.CreateRgb(builder, uint8(rgbList[i][0]), uint8(rgbList[i][1]), uint8(rgbList[i][2]))
	}
	rgbData := builder.EndVector(len(rgbList))

	Points3D.Points3DStart(builder)
	Points3D.Points3DAddPos(builder, posData)
	Points3D.Points3DAddRgb(builder, rgbData)
	Points3D.Points3DAddXmax(builder, float32(meta.Xmax))
	Points3D.Points3DAddXmin(builder, float32(meta.Xmin))
	Points3D.Points3DAddYmax(builder, float32(meta.Ymax))
	Points3D.Points3DAddYmin(builder, float32(meta.Ymin))
	Points3D.Points3DAddZmax(builder, float32(meta.Zmax))
	Points3D.Points3DAddZmin(builder, float32(meta.Zmin))
	Points3D.Points3DAddRmean(builder, float32(meta.Rmean))
	Points3D.Points3DAddGmean(builder, float32(meta.Gmean))
	Points3D.Points3DAddBmean(builder, float32(meta.Bmean))
	Points3D.Points3DAddRstd(builder, float32(meta.Rstd))
	Points3D.Points3DAddGstd(builder, float32(meta.Gstd))
	Points3D.Points3DAddBstd(builder, float32(meta.Bstd))
	p3ds := Points3D.Points3DEnd(builder)

	builder.Finish(p3ds)
	buf := builder.FinishedBytes()

	hashStr, err := utils.GetFileSha256Bytes(buf)
	if err != nil {
		return nil, meta, err
	}
	meta.Sha256 = hashStr

	return buf, meta, nil
}

func (handler *Points3DHandler) Exec() error {
	trainDataDir := filepath.Join(handler.datasetDir, "train")
	trainRawDataDir := filepath.Join(trainDataDir, "raw-data")
	trainAnnoDir := filepath.Join(trainDataDir, "annotation")

	valDataDir := filepath.Join(handler.datasetDir, "val")
	valRawDataDir := filepath.Join(valDataDir, "raw-data")
	valAnnoDir := filepath.Join(valDataDir, "annotation")

	annoTempFile := filepath.Join(handler.datasetDir, "annotation_template.json")

	annoTemp, err := handler.createAnnotationTemplate(annoTempFile)
	if err != nil {
		return err
	}

	trainRawDataViewId, trainAnnoViewId, valRawDataViewId, valAnnoViewId, err := handler.CreateAllDataView(aifsclientgo.POINTS_3D, annoTemp)
	if err != nil {
		return err
	}

	var dataBinHashMap map[string]string
	dataBinHashMap, err = handler.UploadRawData(trainRawDataViewId, trainRawDataDir)
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
	err = handler.UploadAnnoFiles(trainAnnoViewId, trainAnnoDir, trainRawDataIdHashMap, dataBinHashMap, annoTemp)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(70.0, "upload train annotation completed"); err != nil {
		return err
	}
	dataBinHashMap, err = handler.UploadRawData(valRawDataViewId, valRawDataDir)
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

	err = handler.UploadAnnoFiles(valAnnoViewId, valAnnoDir, valRawDataIdHashMap, dataBinHashMap, annoTemp)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(100.0, "dataset zip view decompression completed"); err != nil {
		return err
	}
	return nil
}
