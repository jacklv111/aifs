/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"database/sql"

	"github.com/google/uuid"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/model/constant"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/database"
	"gorm.io/gorm"
)

type BasicDataRepoImpl struct {
}

var BasicDataRepo BasicDataRepoInterface

func init() {
	BasicDataRepo = &BasicDataRepoImpl{}
}

func (repo *BasicDataRepoImpl) GetNameAndType(idList []uuid.UUID) (res []basicdo.DataItemDo, err error) {
	err = collection.BatchRange(idList, BATCH_SIZE, func(batch []uuid.UUID) error {
		var temp []basicdo.DataItemDo
		err := database.Db.Select("id", "name", "type").Where("id in (?)", batch).Find(&temp).Error
		if err != nil {
			return err
		}
		res = append(res, temp...)
		return nil
	})
	return
}

func (repo *BasicDataRepoImpl) CreateBatch(dataItemDoList []basicdo.DataItemDo) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dataItemDoList).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *BasicDataRepoImpl) DeleteBatch(idList []uuid.UUID) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		err := collection.BatchRange(idList, constant.BATCH_SIZE, func(batch []uuid.UUID) error {
			if err := tx.Where("id in (?)", batch).Delete(&basicdo.DataItemDo{}).Error; err != nil {
				return err
			}
			return nil
		})

		return err
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *BasicDataRepoImpl) GetByIdList(idList []uuid.UUID) (dataItemDoList []basicdo.DataItemDo, err error) {
	err = collection.BatchRange(idList, constant.BATCH_SIZE, func(batch []uuid.UUID) error {
		var tempDataItemDoList []basicdo.DataItemDo
		if err = database.Db.Where("id in (?)", batch).Find(&tempDataItemDoList).Error; err != nil {
			return err
		}

		dataItemDoList = append(dataItemDoList, tempDataItemDoList...)
		return nil
	})
	return
}

func (repo *BasicDataRepoImpl) Create(dataItemDo basicdo.DataItemDo) error {
	return database.Db.Create(&dataItemDo).Error
}
