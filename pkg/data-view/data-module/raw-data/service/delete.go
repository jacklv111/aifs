/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	imagerepo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/repo"
	p3drepo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/repo"
	rgbdrepo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/repo"
	"github.com/jacklv111/aifs/pkg/store/manager"
	storeparams "github.com/jacklv111/aifs/pkg/store/value-object"
)

// DeleteRawData 删除 raw data 相关的 meta 和 store 中的数据，用于 gc，只能被系统调用
//
//	@param dataItemList require field: id and type
func DeleteRawData(dataItemList []basicdo.DataItemDo) error {
	dataItemMap := make(map[string][]basicdo.DataItemDo)
	for _, data := range dataItemList {
		dataItemMap[data.Type] = append(dataItemMap[data.Type], data)
	}

	storeParams := getStoreParamsForDelete(dataItemMap)
	err := manager.StoreMgr.Delete(storeParams)
	if err != nil {
		return err
	}
	// meta 必须最后删除，只要中间删除失败，都可以重新获取 meta 信息重新进入删除流程
	err = deleteMeta(dataItemMap)
	if err != nil {
		return err
	}

	return nil
}

func deleteMeta(dataItemMap map[string][]basicdo.DataItemDo) error {
	for key, rawDataList := range dataItemMap {
		switch key {
		case constant.IMAGE:
			err := imagerepo.ImageRawDataRepo.DeleteBatch(basicdo.GetDataItemIdList(rawDataList))
			if err != nil {
				return err
			}
		case constant.RGBD:
			err := rgbdrepo.RgbdRawDataRepo.DeleteBatch(basicdo.GetDataItemIdList(rawDataList))
			if err != nil {
				return err
			}
		case constant.POINTS_3D:
			err := p3drepo.P3dRawDataRepo.DeleteBatch(basicdo.GetDataItemIdList(rawDataList))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getStoreParamsForDelete(dataItemMap map[string][]basicdo.DataItemDo) storeparams.DeleteParams {
	var params storeparams.DeleteParams
	for _, rawDataList := range dataItemMap {
		for _, data := range rawDataList {
			params.AddItem(data.ID)
		}
	}
	return params
}
