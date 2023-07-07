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

	"github.com/google/uuid"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/repo"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

type OcrBo struct {
	basicbo.AnnotationDataImpl
}

func (bo *OcrBo) LoadFromLocal() error {
	return fmt.Errorf("ocr type doesn't implement loadFromLocal")
}

func (bo *OcrBo) GetLabels() []uuid.UUID {
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
	var annoDoList []annodo.AnnotationDo
	for _, data := range dataList {
		cocoData := data.(*OcrBo)
		idList = append(idList, cocoData.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, cocoData.DataItemDo)
		annoDoList = append(annoDoList, cocoData.AnnotationDo)
	}
	err := repo.AnnotationRepo.CreateBatch(dataItemDoList, annoDoList, nil)
	if err != nil {
		return nil, err
	}
	return idList, nil
}
