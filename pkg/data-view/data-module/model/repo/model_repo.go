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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/model/do"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/database"
	"gorm.io/gorm"
)

type modelRepoImpl struct {
}

var ModelRepo modelRepoInterface

func init() {
	ModelRepo = &modelRepoImpl{}
}

func (repo *modelRepoImpl) CreateBatch(dataItemDoList []basicdo.DataItemDo, modelExts []do.ModelExtDo) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dataItemDoList).Error; err != nil {
			return err
		}
		if err := tx.Create(&modelExts).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *modelRepoImpl) DeleteBatch(idList []uuid.UUID) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		err := collection.BatchRange(idList, constant.BATCH_SIZE, func(batch []uuid.UUID) error {
			if err := tx.Where("id in (?)", batch).Delete(&basicdo.DataItemDo{}).Error; err != nil {
				return err
			}

			if err := tx.Where("id in (?)", batch).Delete(&do.ModelExtDo{}).Error; err != nil {
				return err
			}
			return nil
		})

		return err
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *modelRepoImpl) GetByIdList(idList []uuid.UUID) (dataItemDoList []basicdo.DataItemDo, modelExts []do.ModelExtDo, err error) {
	err = collection.BatchRange(idList, constant.BATCH_SIZE, func(batch []uuid.UUID) error {
		var tempDataItemDoList []basicdo.DataItemDo
		if err = database.Db.Where("id in (?)", batch).Find(&tempDataItemDoList).Error; err != nil {
			return err
		}
		var tempList []do.ModelExtDo
		if err = database.Db.Where("id in (?)", batch).Find(&tempList).Error; err != nil {
			return err
		}
		dataItemDoList = append(dataItemDoList, tempDataItemDoList...)
		modelExts = append(modelExts, tempList...)
		return nil
	})
	return
}
