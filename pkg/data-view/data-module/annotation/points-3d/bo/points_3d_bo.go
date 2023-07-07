/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package points3d

import (
	"bufio"
	"compress/zlib"
	"fmt"

	"github.com/google/uuid"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/repo"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

type Points3DAnnotationBo struct {
	basicbo.AnnotationDataImpl
}

func (bo *Points3DAnnotationBo) LoadFromBuffer() (err error) {
	reader, err := zlib.NewReader(bo.ReadSeeker)
	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(reader)
	labelIdMap := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		labelIdMap[line] = ""
	}
	for labelId := range labelIdMap {
		rawDataLabel := annodo.RawDataLabelDo{
			RawDataId:    bo.DataItemId,
			LabelId:      uuid.MustParse(labelId),
			AnnotationId: bo.DataItemDo.ID,
		}
		bo.RawDataLabelList = append(bo.RawDataLabelList, rawDataLabel)
	}

	if err = bo.ResetReader(); err != nil {
		return err
	}
	return nil
}

func (bo *Points3DAnnotationBo) GetLabels() []uuid.UUID {
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
		rgbdAnno := data.(*Points3DAnnotationBo)
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
