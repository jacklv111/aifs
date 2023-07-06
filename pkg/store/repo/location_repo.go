/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"github.com/google/uuid"
	"github.com/jacklv111/aifs/pkg/store/do"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/collection/mapset"
	"github.com/jacklv111/common-sdk/database"
)

type LocationRepoImpl struct {
}

var LocationRepo locationRepoInterface

func init() {
	LocationRepo = &LocationRepoImpl{}
}

func (repo *LocationRepoImpl) Create(doList []do.LocationDo) error {
	return database.Db.Create(&doList).Error
}

func (repo *LocationRepoImpl) FindExistedUkey(ukList []do.LocationUkey) (mapset.Set[do.LocationUkey], error) {
	existed := make([]do.LocationUkey, 0)

	err := collection.BatchRange(do.GetTupleList(ukList), BATCH_SIZE, func(batch [][]interface{}) error {
		var temp []do.LocationUkey
		err := database.Db.Table(do.TABLE_LOCATION).
			Select("data_item_id", "name", "environment").
			Where("(data_item_id, name, environment) in ?", batch).
			Find(&temp).Error
		if err != nil {
			return err
		}
		existed = append(existed, temp...)
		return nil
	})

	if err != nil {
		return nil, err
	}
	idSet := mapset.NewSet(existed...)
	return idSet, nil
}

func (repo *LocationRepoImpl) FindByIdList(idList []uuid.UUID) (result map[do.LocationUkey]do.LocationDo, err error) {
	var resultList []do.LocationDo
	collection.BatchRange(idList, BATCH_SIZE, func(batch []uuid.UUID) error {
		var temp []do.LocationDo
		err := database.Db.Table(do.TABLE_LOCATION).Where("data_item_id in ?", batch).Find(&temp).Error
		if err != nil {
			return err
		}
		resultList = append(resultList, temp...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	result = make(map[do.LocationUkey]do.LocationDo)
	for _, data := range resultList {
		result[do.LocationUkey{DataItemId: data.DataItemId, Name: data.Name, Environment: data.Environment}] = data
	}
	return
}

func (repo *LocationRepoImpl) DeleteByUk(ukList []do.LocationUkey) error {
	err := collection.BatchRange(do.GetTupleList(ukList), BATCH_SIZE, func(batch [][]interface{}) error {
		err := database.Db.
			Where("(data_item_id, name, environment) in ?", batch).
			Delete(&do.LocationDo{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
