/*
 * Created on Sat Jul 08 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"fmt"

	"github.com/google/uuid"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	rawdataconst "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	rawdatasvc "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/service"
	dvdo "github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	dvrepo "github.com/jacklv111/aifs/pkg/data-view/data-view/repo"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
	"github.com/jacklv111/aifs/pkg/store/manager"
	"github.com/jacklv111/common-sdk/collection"
)

type rawDataViewBo struct {
	dataViewBo
}

func (bo *rawDataViewBo) Create() (uuid.UUID, error) {
	if !rawdataconst.HasRawDataType(bo.dataViewDo.RawDataType) {
		return uuid.Nil, fmt.Errorf("invalid raw data type %s, it should be in %v", bo.dataViewDo.RawDataType, rawdataconst.GetRawDataTypeList())
	}

	return bo.dataViewBo.Create()
}

func (bo *rawDataViewBo) UploadRawData(input basicvb.UploadRawDataParam) (err error) {
	if err = bo.loadDataViewDo(); err != nil {
		return
	}

	// upload raw data
	var dataItemIdList []uuid.UUID
	dataItemIdList, err = rawdatasvc.UploadRawData(input, bo.dataViewDo.RawDataType)
	if err != nil {
		return
	}

	// save data view items
	return dvrepo.DataViewRepo.CreateDataViewItemsIgnoreConflict(bo.dataViewDo.ID, dataItemIdList)
}

func (bo *rawDataViewBo) GetAllRawDataLocations() (result dvvb.RawDataLocationResult, err error) {
	result.DataViewId = bo.dataViewDo.ID.String()
	result.ViewType = bo.dataViewDo.ViewType
	result.RawDataType = bo.dataViewDo.RawDataType
	result.DataItemDoList, result.LocationMap, err = bo.getDataItemAndLocations()
	return
}

func (bo *rawDataViewBo) GetHashList(offset int, limit int) (result []basicdo.IdHash, err error) {
	dataItemIdList, err := dvrepo.DataViewRepo.GetDataViewItems(bo.dataViewDo.ID, offset, limit)
	if err != nil {
		return nil, err
	}
	return rawdatasvc.GetHashList(dataItemIdList, bo.dataViewDo.RawDataType)
}

func (bo *rawDataViewBo) GetRawDataType() string {
	return bo.dataViewDo.RawDataType
}

func (bo *rawDataViewBo) GetRawDataList(offset int, limit int, rawDataIdList []string, excludedAnnoViewId, includedAnnotationViewId string) (rawDataList []dvvb.RawDataItem, err error) {
	dataItemIdList, err := dvrepo.DataViewRepo.GetRawDataViewItems(bo.dataViewDo.ID, offset, limit, rawDataIdList, excludedAnnoViewId, includedAnnotationViewId)
	if err != nil {
		return nil, err
	}
	idUrlMap, err := manager.StoreMgr.GetUrlListByIdList(dataItemIdList)
	if err != nil {
		return nil, err
	}

	rawDataList = make([]dvvb.RawDataItem, len(dataItemIdList))
	for idx, dataItemId := range dataItemIdList {
		rawDataList[idx].RawDataId = dataItemId.String()
		rawDataList[idx].Url = idUrlMap[dataItemId]
	}
	return
}

func (bo *rawDataViewBo) DivideDataView(params dvvb.DivideRawDataViewParams) (res dvvb.DivideRawDataViewResult, err error) {
	dataItemList, err := dvrepo.DataViewRepo.GetAllDataViewItems(bo.dataViewDo.ID)
	if err != nil {
		return dvvb.DivideRawDataViewResult{}, err
	}
	ratio := make([]int, 0)
	for _, param := range params.RawDataViewParamList {
		ratio = append(ratio, int(param.Ratio))
	}
	dataItemPartation := collection.DivideItems(dataItemList, ratio)
	for idx, param := range params.RawDataViewParamList {
		items := dataItemPartation[idx]
		subRawDataView := rawDataViewBo{
			dataViewBo: dataViewBo{
				dataViewDo: dvdo.DataViewDo{
					ID:          uuid.New(),
					Name:        param.Name,
					ViewType:    bo.dataViewDo.ViewType,
					Description: param.Description,
					RawDataType: bo.dataViewDo.RawDataType,
				},
			},
		}
		id, err := subRawDataView.Create()
		if err != nil {
			return dvvb.DivideRawDataViewResult{}, err
		}
		err = subRawDataView.CreateDataViewItems(items)
		if err != nil {
			return dvvb.DivideRawDataViewResult{}, err
		}
		res.RawDataViewResultList = append(res.RawDataViewResultList, dvvb.EachRawDataViewResult{
			Name:       param.Name,
			DataViewId: id.String(),
			ItemCount:  int32(len(items)),
		})
	}
	return
}

func (bo *rawDataViewBo) CreateDataViewItems(itemIdList []uuid.UUID) error {
	return dvrepo.DataViewRepo.CreateDataViewItemsIgnoreConflict(bo.dataViewDo.ID, itemIdList)
}

func (bo *rawDataViewBo) GetStatistics() (dvvb.DataViewStatistics, error) {
	count, err := dvrepo.DataViewRepo.GetDataViewItemCount(bo.dataViewDo.ID)
	if err != nil {
		return dvvb.DataViewStatistics{}, err
	}
	dataItemIdList, err := dvrepo.DataViewRepo.GetAllDataViewItems(bo.dataViewDo.ID)
	if err != nil {
		return dvvb.DataViewStatistics{}, err
	}
	totalSize, err := rawdatasvc.GetTotalDataSize(dataItemIdList, bo.dataViewDo.RawDataType)
	if err != nil {
		return dvvb.DataViewStatistics{}, err
	}
	return dvvb.DataViewStatistics{
		ItemCount:     int32(count),
		TotalDataSize: totalSize,
	}, nil
}
