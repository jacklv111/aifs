/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"io"

	"github.com/google/uuid"
	datamodule "github.com/jacklv111/aifs/pkg/data-view/data-module"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/model/bo"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/model/repo"
	"github.com/jacklv111/aifs/pkg/store/manager"
	storevb "github.com/jacklv111/aifs/pkg/store/value-object"
	"github.com/jacklv111/common-sdk/collection/mapset"
)

// getCurrentModelData 查询现有的模型数据，得出统计数据
//
//	@param idList 模型文件数据的 id 列表
//	@return nameMap key: name; value: data item id
//	@return hashSet 现在存在的文件 hash 集合
//	@return err
func getCurrentModelData(idList []uuid.UUID) (nameMap map[string]uuid.UUID, hashSet mapset.Set[string], err error) {
	dataItemList, modelExts, err := repo.ModelRepo.GetByIdList(idList)
	if err != nil {
		return nil, nil, err
	}
	nameMap = make(map[string]uuid.UUID, 0)
	hashSet = mapset.NewSet[string]()
	for _, data := range dataItemList {
		nameMap[data.Name] = data.ID
	}
	for _, data := range modelExts {
		hashSet.Add(data.Sha256)
	}
	return nameMap, hashSet, nil
}

// UploadModelData 上传模型数据。hash 相同的表示这个数据没有变动，就忽略。名字相同，新数据会覆盖旧数据，这里会删除旧的 item，添加新的 item。
//
//	@param input
//	@param dataItemList
//	@return addIdList
//	@return deleteIdList
//	@return err
func UploadModelData(DataFileMap map[string]io.ReadSeeker, dataItemList []uuid.UUID) (addIdList []uuid.UUID, deleteIdList []uuid.UUID, err error) {
	nameMap, hashSet, err := getCurrentModelData(dataItemList)
	if err != nil {
		return nil, nil, err
	}
	modelList := make([]bo.ModelDataBo, 0)
	for key, data := range DataFileMap {
		modelBo := bo.BuildFromBuffer(data, key)
		err := modelBo.LoadFromBuffer()
		if err != nil {
			return nil, nil, err
		}

		// 现有数据存在相同名字，进行数据覆盖，即删除旧数据，添加新数据
		if id, existed := nameMap[modelBo.GetName()]; existed {
			// 现有数据同名且存在相同 hash， 不需要再上传
			if hashSet.Contains(modelBo.GetHash()) {
				continue
			}
			deleteIdList = append(deleteIdList, id)
		}
		// 添加数据
		modelList = append(modelList, *modelBo)
	}

	// save meta
	addIdList, err = bo.CreateBatch(modelList)
	if err != nil {
		return
	}

	// save the annotation data
	storeParams, err := getStoreParamRemote(modelList)

	if err != nil {
		return nil, nil, err
	}

	err = manager.StoreMgr.Upload(storeParams)
	if err != nil {
		return nil, nil, err
	}
	return
}

func getStoreParamRemote(modelList []bo.ModelDataBo) (storevb.UploadParams, error) {
	var params storevb.UploadParams
	params.DataType = datamodule.MODEL
	for _, model := range modelList {
		params.AddItem(model.GetId(), model.GetReader(), model.GetName())
	}
	return params, nil
}
