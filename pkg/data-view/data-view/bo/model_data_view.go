/*
 * Created on Sat Jul 08 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"database/sql"
	"strconv"

	"github.com/google/uuid"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/model/service"
	dvdo "github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	dvrepo "github.com/jacklv111/aifs/pkg/data-view/data-view/repo"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
)

type modelViewBo struct {
	dataViewBo
}

func (bo *modelViewBo) UploadModelData(input basicvb.UploadModelParams) (err error) {
	if err = bo.loadDataViewDo(); err != nil {
		return
	}

	dataItemList, err := dvrepo.DataViewRepo.GetAllDataViewItems(bo.dataViewDo.ID)
	if err != nil {
		return err
	}

	if err = updateIfNecessary(bo.dataViewDo.ID, input.Pairs); err != nil {
		return err
	}

	addList, deleteList, err := service.UploadModelData(input.DataFileMap, dataItemList)
	if err != nil {
		return
	}

	// save data view items
	err = dvrepo.DataViewRepo.CreateDataViewItemsIgnoreConflict(bo.dataViewDo.ID, addList)
	if err != nil {
		return err
	}
	err = dvrepo.DataViewRepo.DeleteDataViewItem(bo.dataViewDo.ID, deleteList)
	if err != nil {
		return err
	}

	return
}

func (bo *modelViewBo) GetModelDataLocations() (result dvvb.ModelLocationResult, err error) {
	result.DataViewId = bo.dataViewDo.ID.String()
	result.ViewType = bo.dataViewDo.ViewType
	result.DataItemDoList, result.LocationMap, err = bo.getDataItemAndLocations()
	return
}

// private -------------------------------------------------------------------------------------------------

func updateIfNecessary(dataViewId uuid.UUID, pairs map[string]string) (err error) {
	data := dvdo.DataViewDo{ID: dataViewId}
	shouldUpdate := false
	for key, val := range pairs {
		if key == "progress" {
			data.Progress, err = strconv.ParseFloat(val, 32)
			if err != nil {
				return err
			}
			shouldUpdate = true
		}
		if key == "commitId" {
			data.CommitId = sql.NullString{String: val, Valid: true}
			shouldUpdate = true
		}
	}
	if shouldUpdate {
		return dvrepo.DataViewRepo.Updates(data)
	}
	return nil
}
