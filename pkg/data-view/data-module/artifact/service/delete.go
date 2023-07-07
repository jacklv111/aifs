/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"github.com/jacklv111/aifs/pkg/data-view/data-module/artifact/repo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/store/manager"
	storeparams "github.com/jacklv111/aifs/pkg/store/value-object"
)

// DeleteArtifactData 删除 artifact 相关的 meta 和 store 中的数据，用于 gc，只能被系统调用
//
//	@param dataItemList require field: id and type
func DeleteArtifactData(dataItemList []basicdo.DataItemDo) error {
	storeParams := getStoreParamsForDelete(dataItemList)
	err := manager.StoreMgr.Delete(storeParams)
	if err != nil {
		return err
	}

	// meta 必须最后删除，只要中间删除失败，都可以重新获取 meta 信息重新进入删除流程
	err = repo.ArtifactRepo.DeleteBatch(basicdo.GetDataItemIdList(dataItemList))
	if err != nil {
		return err
	}

	return nil
}

func getStoreParamsForDelete(dataItemList []basicdo.DataItemDo) storeparams.DeleteParams {
	var params storeparams.DeleteParams
	for _, data := range dataItemList {
		params.AddItem(data.ID)
	}
	return params
}
