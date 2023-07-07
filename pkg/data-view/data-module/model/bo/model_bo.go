/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"io"

	"github.com/google/uuid"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	modeldo "github.com/jacklv111/aifs/pkg/data-view/data-module/model/do"
	modelrepo "github.com/jacklv111/aifs/pkg/data-view/data-module/model/repo"
	"github.com/jacklv111/common-sdk/utils"
)

type ModelDataBo struct {
	bo.DataBaseImpl
	modeldo.ModelExtDo
}

func (bo *ModelDataBo) LoadFromBuffer() error {
	bytes, err := io.ReadAll(bo.GetReader())
	if err != nil {
		return err
	}
	err = bo.ResetReader()
	if err != nil {
		return err
	}

	bo.Sha256, err = utils.GetFileSha256Bytes(bytes)
	if err != nil {
		return err
	}
	bo.Size = len(bytes)

	return nil
}

func (bo *ModelDataBo) GetHash() string {
	return bo.Sha256
}

func (bo *ModelDataBo) GetName() string {
	return bo.Name
}

func (bo *ModelDataBo) GetId() uuid.UUID {
	return bo.ID
}

// CreateBatch 批量插入 metadata
//
//	@param dataList
//	@return []uuid.UUID
//	@return error
func CreateBatch(dataList []ModelDataBo) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	var dataItemDoList []basicdo.DataItemDo
	var modelExts []modeldo.ModelExtDo
	for _, data := range dataList {
		idList = append(idList, data.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, data.DataItemDo)
		modelExts = append(modelExts, data.ModelExtDo)
	}
	err := modelrepo.ModelRepo.CreateBatch(dataItemDoList, modelExts)
	if err != nil {
		return nil, err
	}
	return idList, nil
}
