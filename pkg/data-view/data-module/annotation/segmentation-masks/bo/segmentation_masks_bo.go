/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"fmt"
	"image"
	"image/color"

	"github.com/google/uuid"
	annotempbo "github.com/jacklv111/aifs/pkg/annotation-template/bo"
	annoDo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/repo"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/common-sdk/collection/mapset"
	"github.com/jacklv111/common-sdk/utils"
)

type SegmentationMasksBo struct {
	basicbo.AnnotationDataImpl
	annoData image.Image
	AnnoTemp annotempbo.AnnotationTemplateBoInterface
}

func (bo *SegmentationMasksBo) GetLabels() []uuid.UUID {
	var labels []uuid.UUID
	for _, data := range bo.RawDataLabelList {
		labels = append(labels, data.LabelId)
	}
	return labels
}

func (bo *SegmentationMasksBo) LoadFromBuffer() (err error) {
	bo.annoData, err = utils.ReadImage(bo.ReadSeeker)
	if err != nil {
		return err
	}

	lableIdSet := mapset.NewSet[uuid.UUID]()
	for i := 0; i < bo.annoData.Bounds().Dx(); i++ {
		for j := 0; j < bo.annoData.Bounds().Dy(); j++ {
			color, ok := bo.annoData.At(i, j).(color.Gray)
			if !ok {
				return fmt.Errorf("segmentation masks annotation should be gray")
			}
			labelId := bo.AnnoTemp.GetLabelIdByColor(int32(color.Y))
			if labelId == uuid.Nil {
				return fmt.Errorf("pixel value %d is not defined", color.Y)
			}
			if lableIdSet.Contains(labelId) {
				continue
			}
			lableIdSet.Add(labelId)

			rawDataLabel := annoDo.RawDataLabelDo{
				RawDataId:    bo.DataItemId,
				LabelId:      labelId,
				AnnotationId: bo.DataItemDo.ID,
			}
			bo.RawDataLabelList = append(bo.RawDataLabelList, rawDataLabel)
		}
	}
	if err = bo.ResetReader(); err != nil {
		return err
	}
	return nil
}

// CreateBatch 批量插入 metadata
//
//	@param dataList
//	@return []uuid.UUID
//	@return error
func CreateBatch(dataList []basicbo.AnnotationData) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	var dataItemDoList []basicdo.DataItemDo
	var annoDoList []annoDo.AnnotationDo
	var rawDataLabelList []annoDo.RawDataLabelDo
	for _, data := range dataList {
		segMasksData := data.(*SegmentationMasksBo)
		idList = append(idList, segMasksData.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, segMasksData.DataItemDo)
		annoDoList = append(annoDoList, segMasksData.AnnotationDo)
		rawDataLabelList = append(rawDataLabelList, segMasksData.RawDataLabelList...)
	}
	err := repo.AnnotationRepo.CreateBatch(dataItemDoList, annoDoList, rawDataLabelList)
	if err != nil {
		return nil, err
	}
	return idList, nil
}
