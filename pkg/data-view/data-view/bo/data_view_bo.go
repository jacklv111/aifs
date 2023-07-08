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
	dv "github.com/jacklv111/aifs/pkg/data-view"
	datamodule "github.com/jacklv111/aifs/pkg/data-view/data-module"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	basicrepo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/repo"
	dvdo "github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	dvrepo "github.com/jacklv111/aifs/pkg/data-view/data-view/repo"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
	"github.com/jacklv111/aifs/pkg/store/manager"
	storevb "github.com/jacklv111/aifs/pkg/store/value-object"
	"github.com/jacklv111/common-sdk/errors"
)

type dataViewBo struct {
	dataViewDo dvdo.DataViewDo
	// 懒加载
	hasDataViewDo bool
}

func (bo *dataViewBo) GetId() uuid.UUID {
	return bo.dataViewDo.ID
}

func (bo *dataViewBo) GetViewType() string {
	return bo.dataViewDo.ViewType
}

func (bo *dataViewBo) Create() (uuid.UUID, error) {
	if bo.dataViewDo.Name == "" {
		return uuid.Nil, fmt.Errorf("data view name is null, create failed")
	}
	return bo.dataViewDo.ID, dvrepo.DataViewRepo.Create(bo.dataViewDo)
}

func (bo *dataViewBo) Delete() error {
	if bo.dataViewDo.ID == uuid.Nil {
		return fmt.Errorf("data view id is null, delete failed")
	}
	err := dvrepo.DataViewRepo.SoftDelete(bo.dataViewDo.ID)
	if err != nil {
		return err
	}
	bo.dataViewDo.ID = uuid.Nil
	return nil
}

func (bo *dataViewBo) GetDetails() (res dvvb.DataViewDetails, err error) {
	res = dvvb.DataViewDetails{DataViewDo: bo.dataViewDo}
	return
}

func (bo *dataViewBo) DeleteDataItem(idList []uuid.UUID) error {
	if bo.dataViewDo.ID == uuid.Nil {
		return fmt.Errorf("data view id is null, delete failed")
	}
	err := dvrepo.DataViewRepo.DeleteDataViewItem(bo.dataViewDo.ID, idList)
	if err != nil {
		return err
	}
	return nil
}

func (bo *dataViewBo) HardDelete() error {
	return dvrepo.DataViewRepo.HardDelete(bo.dataViewDo.ID)
}

// common private func
func (bo *dataViewBo) loadDataViewDo() (err error) {
	if bo.hasDataViewDo {
		return
	}
	bo.dataViewDo, err = dvrepo.DataViewRepo.GetById(bo.dataViewDo.ID)
	if err != nil {
		if err.Error() == errors.RECORD_NOT_FOUND {
			err = dv.ErrDataViewNotFound
		}
		return
	}
	bo.hasDataViewDo = true
	return
}

func (bo *dataViewBo) MergeDataItems(other uuid.UUID) error {
	otherBo := &dataViewBo{dataViewDo: dvdo.DataViewDo{ID: other}}
	err := otherBo.loadDataViewDo()
	if err != nil {
		return err
	}

	if !bo.isSameType(otherBo) {
		return fmt.Errorf("data view type is not same, add failed")
	}
	return dvrepo.DataViewRepo.MergeTo(bo.dataViewDo.ID, otherBo.GetId())
}

func (bo *dataViewBo) getDataItemAndLocations() (dataItemDoList []basicdo.DataItemDo, locationMap map[uuid.UUID]storevb.LocationResult, err error) {
	// get data item in data view
	dataItemIdList, err := dvrepo.DataViewRepo.GetAllDataViewItems(bo.dataViewDo.ID)
	if err != nil {
		return nil, nil, err
	}
	// get data item name and type
	dataItemDoList, err = basicrepo.BasicDataRepo.GetNameAndType(dataItemIdList)
	if err != nil {
		return nil, nil, err
	}
	// get locations
	locationMap, err = manager.StoreMgr.GetByIdList(dataItemIdList)
	if err != nil {
		return nil, nil, err
	}
	return dataItemDoList, locationMap, nil
}

func (bo *dataViewBo) isSameType(other *dataViewBo) bool {
	if bo.dataViewDo.ViewType != other.dataViewDo.ViewType {
		return false
	}
	switch bo.dataViewDo.ViewType {
	case datamodule.RAW_DATA:
		return bo.dataViewDo.RawDataType == other.dataViewDo.RawDataType
	case datamodule.ANNOTATION:
		return bo.dataViewDo.AnnotationTemplateId == other.dataViewDo.AnnotationTemplateId
	default:
		return false
	}
}
