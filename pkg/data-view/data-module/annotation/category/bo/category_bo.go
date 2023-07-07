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
	"os"

	"github.com/google/uuid"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/repo"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

type CategoryBo struct {
	basicbo.AnnotationDataImpl
}

type Category struct {
	AnnotationTemplateId string
	LabelId              string
	RawDataId            string
}

func (bo *CategoryBo) UnmarshalJSON(data []byte) error {
	var category Category
	if err := json.Unmarshal(data, &category); err != nil {
		return err
	}
	bo.DataItemId = uuid.MustParse(category.RawDataId)
	bo.AnnotationTemplateId = uuid.MustParse(category.AnnotationTemplateId)
	rawDataLabel := annodo.RawDataLabelDo{
		RawDataId:    bo.DataItemId,
		LabelId:      uuid.MustParse(category.LabelId),
		AnnotationId: bo.DataItemDo.ID,
	}
	bo.RawDataLabelList = append(bo.RawDataLabelList, rawDataLabel)
	return nil
}

func (bo *CategoryBo) LoadFromLocal() error {
	content, err := os.ReadFile(bo.GetLocalPath())
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, bo)
	if err != nil {
		return err
	}

	return nil
}

func (bo *CategoryBo) GetLabels() []uuid.UUID {
	var labels []uuid.UUID
	for _, data := range bo.RawDataLabelList {
		labels = append(labels, data.LabelId)
	}
	return labels
}

// CreateBatch 批量插入 metadata
//
//	@param categoryList
//	@return []uuid.UUID
//	@return error
func CreateBatch(categoryList []basicbo.AnnotationData) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	var dataItemDoList []basicdo.DataItemDo
	var annoDoList []annodo.AnnotationDo
	var rawDataLabelList []annodo.RawDataLabelDo
	for _, data := range categoryList {
		category := data.(*CategoryBo)
		idList = append(idList, category.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, category.DataItemDo)
		annoDoList = append(annoDoList, category.AnnotationDo)
		rawDataLabelList = append(rawDataLabelList, category.RawDataLabelList...)
	}
	err := repo.AnnotationRepo.CreateBatch(dataItemDoList, annoDoList, rawDataLabelList)
	if err != nil {
		return nil, err
	}
	return idList, nil
}
