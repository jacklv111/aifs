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
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/database"
	"gorm.io/gorm"
)

type annotationRepoImpl struct {
}

var AnnotationRepo annotationRepoInterface

func init() {
	AnnotationRepo = &annotationRepoImpl{}
}

func (repo *annotationRepoImpl) CreateBatch(dataItemDoList []basicdo.DataItemDo, annoDoList []annodo.AnnotationDo, rawDataLabelList []annodo.RawDataLabelDo) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dataItemDoList).Error; err != nil {
			return err
		}
		if err := tx.Create(&annoDoList).Error; err != nil {
			return err
		}
		if rawDataLabelList != nil {
			if err := tx.Create(&rawDataLabelList).Error; err != nil {
				return err
			}
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *annotationRepoImpl) DeleteBatch(idList []uuid.UUID) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		err := collection.BatchRange(idList, BATCH_SIZE, func(batch []uuid.UUID) error {
			if err := tx.Where("id in (?)", batch).Delete(&basicdo.DataItemDo{}).Error; err != nil {
				return err
			}

			if err := tx.Where("id in (?)", batch).Delete(&annodo.AnnotationDo{}).Error; err != nil {
				return err
			}

			if err := tx.Where("annotation_id in (?)", batch).Delete(&annodo.RawDataLabelDo{}).Error; err != nil {
				return err
			}
			return nil
		})

		return err
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *annotationRepoImpl) GetByIdList(idList []uuid.UUID) (dataItemDoList []basicdo.DataItemDo, annoDoList []annodo.AnnotationDo, err error) {
	err = collection.BatchRange(idList, BATCH_SIZE, func(batch []uuid.UUID) error {
		var tempDataItemDoList []basicdo.DataItemDo
		if err = database.Db.Where("id in (?)", batch).Find(&tempDataItemDoList).Error; err != nil {
			return err
		}
		var tempAnnoDoList []annodo.AnnotationDo
		if err = database.Db.Where("id in (?)", batch).Find(&tempAnnoDoList).Error; err != nil {
			return err
		}
		dataItemDoList = append(dataItemDoList, tempDataItemDoList...)
		annoDoList = append(annoDoList, tempAnnoDoList...)
		return nil
	})
	return
}

func (repo *annotationRepoImpl) GetRawDataLabelByIdList(idList []uuid.UUID) (rawDataLabelList []annodo.RawDataLabelDo, err error) {
	err = collection.BatchRange(idList, BATCH_SIZE, func(batch []uuid.UUID) error {
		var temp []annodo.RawDataLabelDo
		if err = database.Db.Where("annotation_id in (?)", batch).Find(&temp).Error; err != nil {
			return err
		}
		rawDataLabelList = append(rawDataLabelList, temp...)
		return nil
	})
	return
}
