/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"github.com/google/uuid"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/repo"
)

type Points3DBo struct {
	basicbo.DataBaseImpl
	do.Points3DExtDo
}

// FixDataItemId 对于已经存在的 raw data，使用其已经存在的 id
//
//	@param imageList
//	@return error
func FixDataItemId(dataList []basicbo.DataInterface) error {
	var extDoList []do.Points3DExtDo
	for _, data := range dataList {
		p3dRawDataBo := data.(*Points3DBo)
		extDoList = append(extDoList, p3dRawDataBo.Points3DExtDo)
	}
	existedHashMap, err := repo.P3dRawDataRepo.FindExistedByHash(do.GetHashList(extDoList))
	if err != nil {
		return err
	}
	for _, data := range dataList {
		p3dRawDataBo := data.(*Points3DBo)
		id, ok := existedHashMap[p3dRawDataBo.Points3DExtDo.Sha256]
		if !ok {
			continue
		}
		p3dRawDataBo.DataItemDo.ID = id
		p3dRawDataBo.Points3DExtDo.ID = id
	}
	return nil
}

// CreateBatch 批量插入数据
//
//	@param imageList
//	@return []uuid.UUID
//	@return error
func CreateBatch(imageList []basicbo.DataInterface) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	var dataItemDoList []basicdo.DataItemDo
	var extDoList []do.Points3DExtDo

	for _, data := range imageList {
		p3dRawDataBo := data.(*Points3DBo)
		idList = append(idList, p3dRawDataBo.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, p3dRawDataBo.DataItemDo)
		extDoList = append(extDoList, p3dRawDataBo.Points3DExtDo)
	}

	err := repo.P3dRawDataRepo.CreateBatch(dataItemDoList, extDoList)
	if err != nil {
		return nil, err
	}

	return idList, nil
}
