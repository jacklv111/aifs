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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/do"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type rgbdRawDataRepoImpl struct {
}

var RgbdRawDataRepo RgbdRawDataRepoInterface

func init() {
	RgbdRawDataRepo = &rgbdRawDataRepoImpl{}
}

func (repo *rgbdRawDataRepoImpl) CreateBatch(dataItemDoList []basicdo.DataItemDo, rgbdExtDoList []do.RgbdExtDo) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		// do nothing if id exists
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&dataItemDoList).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&rgbdExtDoList).Error; err != nil {
			return err
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *rgbdRawDataRepoImpl) DeleteBatch(idList []uuid.UUID) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		err := collection.BatchRange(idList, BATCH_SIZE, func(batch []uuid.UUID) error {
			if err := tx.Where("id in (?)", batch).Delete(&basicdo.DataItemDo{}).Error; err != nil {
				return err
			}

			if err := tx.Where("id in (?)", batch).Delete(&do.RgbdExtDo{}).Error; err != nil {
				return err
			}
			return nil
		})

		return err
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *rgbdRawDataRepoImpl) FindExistedByHash(sha256List []string) (map[string]uuid.UUID, error) {
	var existedItemList []do.RgbdExtDo
	existedSha256Map := make(map[string]uuid.UUID)

	err := collection.BatchRange(sha256List, BATCH_SIZE, func(batch []string) error {
		err := database.Db.Select("id", "sha256").Where("sha256 in ?", batch).FindInBatches(&existedItemList, BATCH_SIZE, func(tx *gorm.DB, batch int) error {
			for _, data := range existedItemList {
				existedSha256Map[data.Sha256] = data.ID
			}
			return nil
		}).Error

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return existedSha256Map, nil
}

func (repo *rgbdRawDataRepoImpl) GetHashList(dataItemIdList []uuid.UUID) (res []basicdo.IdHash, err error) {
	res = make([]basicdo.IdHash, 0)
	temp := make([]basicdo.IdHash, 0)
	err = collection.BatchRange(dataItemIdList, BATCH_SIZE, func(batch []uuid.UUID) error {
		err := database.Db.Model(do.RgbdExtDo{}).Select("id", "sha256").Where("id in ?", batch).Find(&temp).Error
		if err != nil {
			return err
		}
		res = append(res, temp...)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}
