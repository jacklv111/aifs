/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"encoding/json"
	"io"
	"os"

	"github.com/google/uuid"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/repo"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

type CocoBo struct {
	basicbo.AnnotationDataImpl
	AnnoData []CocoAnnoData
}

type CocoAnnoData struct {
	LabelId uuid.UUID
	// [polygon] 或者 RLE 格式
	Segmentation interface{}
	IsCrowd      int
	Area         float32
	NumKeyPoints int
	KeyPoints    []int
	// x, y, width, height
	Bbox           []float32
	PredictedIou   float32
	PointCoords    [][]float32
	CropBox        []float32
	StabilityScore float32
	Id             int
}

func (cocoBo *CocoBo) UnmarshalJSON(data []byte) error {
	type Alias CocoBo
	wrapper := &struct {
		RawDataId string
		*Alias
	}{
		Alias: (*Alias)(cocoBo),
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	cocoBo.DataItemId = uuid.MustParse(wrapper.RawDataId)
	return nil
}

func (bo *CocoBo) LoadFromLocal() error {
	content, err := os.ReadFile(bo.GetLocalPath())
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, bo)
	if err != nil {
		return err
	}

	for _, data := range bo.AnnoData {
		// 存在没有 label 的数据，比如 sam 数据集
		if data.LabelId == uuid.Nil {
			continue
		}
		rawDataLabel := annodo.RawDataLabelDo{
			RawDataId:    bo.DataItemId,
			LabelId:      data.LabelId,
			AnnotationId: bo.DataItemDo.ID,
		}
		bo.RawDataLabelList = append(bo.RawDataLabelList, rawDataLabel)
	}

	return nil
}

func (bo *CocoBo) LoadFromBuffer() (err error) {
	bytes, err := io.ReadAll(bo.ReadSeeker)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, bo)
	if err != nil {
		return err
	}
	for _, data := range bo.AnnoData {
		// 存在没有 label 的数据，比如 sam 数据集
		if data.LabelId == uuid.Nil {
			continue
		}
		rawDataLabel := annodo.RawDataLabelDo{
			RawDataId:    bo.DataItemId,
			LabelId:      data.LabelId,
			AnnotationId: bo.DataItemDo.ID,
		}
		bo.RawDataLabelList = append(bo.RawDataLabelList, rawDataLabel)
	}

	if err = bo.ResetReader(); err != nil {
		return err
	}
	return nil
}

func (bo *CocoBo) GetLabels() []uuid.UUID {
	var labels []uuid.UUID
	for _, data := range bo.RawDataLabelList {
		labels = append(labels, data.LabelId)
	}
	return labels
}

// CreateBatch 批量插入 metadata
//
//	@param dataList
//	@return []uuid.UUID
//	@return error
func CreateBatch(dataList []basicbo.AnnotationData) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	var dataItemDoList []basicdo.DataItemDo
	var annoDoList []annodo.AnnotationDo
	var rawDataLabelList []annodo.RawDataLabelDo
	for _, data := range dataList {
		cocoData := data.(*CocoBo)
		idList = append(idList, cocoData.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, cocoData.DataItemDo)
		annoDoList = append(annoDoList, cocoData.AnnotationDo)
		rawDataLabelList = append(rawDataLabelList, cocoData.RawDataLabelList...)
	}
	err := repo.AnnotationRepo.CreateBatch(dataItemDoList, annoDoList, rawDataLabelList)
	if err != nil {
		return nil, err
	}
	return idList, nil
}
