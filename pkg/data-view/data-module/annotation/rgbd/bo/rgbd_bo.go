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

type RgbdAnnotationBo struct {
	basicbo.AnnotationDataImpl
	BoundingBoxList []boundingBox `json:"BoundingBoxList"`
}

type boundingBox struct {
	LabelId       uuid.UUID `json:"LabelId"`
	BoundingBox2D bBox2D    `json:"BoundingBox2D"`
	BoundingBox3D bBox3D    `json:"BoundingBox3D"`
}

type bBox2D struct {
	X      float32 `json:"X"`
	Y      float32 `json:"Y"`
	Width  float32 `json:"Width"`
	Height float32 `json:"Height"`
}

type bBox3D struct {
	X     float32 `json:"X"`
	Y     float32 `json:"Y"`
	Z     float32 `json:"Z"`
	XSize float32 `json:"XSize"`
	YSize float32 `json:"YSize"`
	ZSize float32 `json:"ZSize"`
	YawX  float32 `json:"YawX"`
	YawY  float32 `json:"YawY"`
	YawZ  float32 `json:"YawZ"`
}

func (bo *RgbdAnnotationBo) UnmarshalJSON(data []byte) error {
	type Alias RgbdAnnotationBo
	wrapper := &struct {
		RawDataId string
		*Alias
	}{
		Alias: (*Alias)(bo),
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	bo.DataItemId = uuid.MustParse(wrapper.RawDataId)
	return nil
}

func (bo *RgbdAnnotationBo) LoadFromLocal() error {
	content, err := os.ReadFile(bo.GetLocalPath())
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, bo)
	if err != nil {
		return err
	}

	for _, data := range bo.BoundingBoxList {
		rawDataLabel := annodo.RawDataLabelDo{
			RawDataId:    bo.DataItemId,
			LabelId:      data.LabelId,
			AnnotationId: bo.DataItemDo.ID,
		}
		bo.RawDataLabelList = append(bo.RawDataLabelList, rawDataLabel)
	}

	return nil
}

func (bo *RgbdAnnotationBo) LoadFromBuffer() (err error) {
	bytes, err := io.ReadAll(bo.ReadSeeker)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, bo)
	if err != nil {
		return err
	}
	for _, data := range bo.BoundingBoxList {
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

func (bo *RgbdAnnotationBo) GetLabels() []uuid.UUID {
	var labels []uuid.UUID
	for _, data := range bo.RawDataLabelList {
		labels = append(labels, data.LabelId)
	}
	return labels
}

// CreateBatch 批量插入 metadata
//
//	@param annoList
//	@return []uuid.UUID
//	@return error
func CreateBatch(annoList []basicbo.AnnotationData) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	var dataItemDoList []basicdo.DataItemDo
	var annoDoList []annodo.AnnotationDo
	var rawDataLabelList []annodo.RawDataLabelDo
	for _, data := range annoList {
		rgbdAnno := data.(*RgbdAnnotationBo)
		idList = append(idList, rgbdAnno.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, rgbdAnno.DataItemDo)
		annoDoList = append(annoDoList, rgbdAnno.AnnotationDo)
		rawDataLabelList = append(rawDataLabelList, rgbdAnno.RawDataLabelList...)
	}
	err := repo.AnnotationRepo.CreateBatch(dataItemDoList, annoDoList, rawDataLabelList)
	if err != nil {
		return nil, err
	}
	return idList, nil
}
