/*
 * Created on Sat Jul 08 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	dvdo "github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
	"github.com/jacklv111/common-sdk/collection"
	"github.com/jacklv111/common-sdk/collection/mapset"
	"github.com/jacklv111/common-sdk/database"
	"github.com/jacklv111/common-sdk/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type dataViewRepoImpl struct {
}

var DataViewRepo dataViewRepoInterface

func init() {
	DataViewRepo = &dataViewRepoImpl{}
}

func (repo *dataViewRepoImpl) ExistsById(dataViewId uuid.UUID) (bool, error) {
	var count int64
	err := database.Db.Model(&dvdo.DataViewDo{}).Where("id = ?", dataViewId.String()).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (repo *dataViewRepoImpl) Create(data dvdo.DataViewDo) error {
	return database.Db.Create(&data).Error
}

func (repo *dataViewRepoImpl) GetList(options dvvb.DataViewListQueryOptions) (res []dvdo.DataViewDo, err error) {
	query := database.Db.Offset(options.Offset).Limit(options.Limit)
	if options.HasDataViewList() {
		query.Where("id in (?)", options.DataViewIdList)
	}
	if options.HasNameFilter() {
		query.Where("name like ?", fmt.Sprintf("%%%s%%", options.DataViewName))
	}
	err = query.Find(&res).Error
	return
}

func (repo *dataViewRepoImpl) GetById(dataViewId uuid.UUID) (res dvdo.DataViewDo, err error) {
	err = database.Db.First(&res, dataViewId).Error
	return
}

func (repo *dataViewRepoImpl) GetDataViewItemCount(dataViewId uuid.UUID) (count int64, err error) {
	err = database.Db.Model(&dvdo.DataViewItemDo{}).Where("data_view_id = ?", dataViewId.String()).Count(&count).Error
	return
}

func (repo *dataViewRepoImpl) SoftDelete(dataViewId uuid.UUID) error {
	if err := database.Db.Delete(&dvdo.DataViewDo{ID: dataViewId}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *dataViewRepoImpl) HardDelete(dataViewId uuid.UUID) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Delete(&dvdo.DataViewDo{ID: dataViewId}).Error; err != nil {
			return err
		}
		if err := tx.Where("data_view_id = ?", dataViewId.String()).Delete(&dvdo.DataViewItemDo{}).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *dataViewRepoImpl) DeleteDataViewItem(dataViewId uuid.UUID, dataViewItemIdList []uuid.UUID) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("data_view_id = ?", dataViewId.String()).Where("data_item_id in ?", dataViewItemIdList).Delete(&dvdo.DataViewItemDo{})
		if res.Error != nil {
			return res.Error
		}
		if int(res.RowsAffected) != len(dataViewItemIdList) {
			return fmt.Errorf("some items do not exist")
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *dataViewRepoImpl) CreateDataViewItemsIgnoreConflict(dataViewId uuid.UUID, itemIdList []uuid.UUID) error {
	dataItemList := make([]dvdo.DataViewItemDo, 0)
	for _, data := range itemIdList {
		dataItemList = append(dataItemList, dvdo.DataViewItemDo{DataViewId: dataViewId, DataItemId: data})
	}
	result := database.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&dataItemList)
	log.Infof("create raw data view items, insert %d new tuples and %d tuples conflict", result.RowsAffected, len(itemIdList)-int(result.RowsAffected))
	return result.Error
}

func (repo *dataViewRepoImpl) CreateAnnotationDataViewItems(itemList []dvdo.DataViewItemDo, annotationTemplateId uuid.UUID) error {
	if len(itemList) == 0 {
		return nil
	}
	dataViewId := itemList[0].DataViewId
	annoIdList := dvdo.GetDataItemIdList(itemList)
	return database.Db.Transaction(func(tx *gorm.DB) error {
		// 如果结果太大，可以使用临时表进行优化
		var annoIdToBeDeleted []uuid.UUID

		err := collection.BatchRange(annoIdList, BATCH_SIZE, func(batch []uuid.UUID) error {
			var temp []uuid.UUID
			// 这些 raw data 的 annotation 被改变了
			annoChangedRawDataId := tx.Table(annodo.TABLE_ANNOTATION).Select("annotations.data_item_id").Where("annotations.id in ?", batch)
			// 找出原来对这些 raw data 进行标注的数据
			res := tx.Table(dvdo.TABLE_DATA_VIEW_ITEM).
				Select("annotations.id").
				Joins("left join annotations on data_view_items.data_item_id=annotations.id").
				Where("data_view_items.data_view_id = ?", dataViewId).
				Where("annotations.annotation_template_id = ?", annotationTemplateId).
				Where("annotations.data_item_id in (?)", annoChangedRawDataId).
				Find(&temp)

			if res.Error != nil {
				return nil
			}
			annoIdToBeDeleted = append(annoIdToBeDeleted, temp...)
			return nil
		})

		if err != nil {
			return err
		}

		log.Infof("%d items need to be deleted", len(annoIdToBeDeleted))

		collection.BatchRange(annoIdToBeDeleted, BATCH_SIZE, func(batch []uuid.UUID) error {
			res := tx.Where("data_item_id in (?)", batch).Where("data_view_id = ?", dataViewId).Delete(&dvdo.DataViewItemDo{})
			return res.Error
		})

		res := tx.Create(&itemList)
		if res.Error != nil {
			return nil
		}

		log.Infof("%d items are created", res.RowsAffected)

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *dataViewRepoImpl) GetAllDataViewItems(dataViewId uuid.UUID) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	err := database.Db.Model(dvdo.DataViewItemDo{}).Select("data_item_id").Where("data_view_id = ?", dataViewId).Find(&idList).Error
	return idList, err
}

func (repo *dataViewRepoImpl) GetInvalidId(dataViewIdList []uuid.UUID) ([]uuid.UUID, error) {
	var existedIds []uuid.UUID
	err := database.Db.Model(&dvdo.DataViewDo{}).Select("id").Where("id in (?)", dataViewIdList).Find(&existedIds).Error
	if err != nil {
		return nil, err
	}
	idSet := mapset.NewSet(existedIds...)

	var invalidIds []uuid.UUID
	for _, id := range dataViewIdList {
		if !idSet.Contains(id) {
			invalidIds = append(invalidIds, id)
		}
	}

	return invalidIds, nil
}

func (repo *dataViewRepoImpl) GetInvalidDataItems(dataViewId uuid.UUID, dataItemIdList []uuid.UUID) (invalidDataItemIdList []uuid.UUID, err error) {
	var existedIds []uuid.UUID
	var batchRes []uuid.UUID
	err = collection.BatchRange(dataItemIdList, BATCH_SIZE, func(batch []uuid.UUID) error {
		err = database.Db.Model(&dvdo.DataViewItemDo{}).Select("data_item_id").Where("data_view_id = (?) and data_item_id in (?)", dataViewId, batch).Find(&batchRes).Error
		if err != nil {
			return err
		}
		existedIds = append(existedIds, batchRes...)
		return nil
	})
	if err != nil {
		return
	}

	idSet := mapset.NewSet(existedIds...)
	for _, id := range dataItemIdList {
		if !idSet.Contains(id) {
			invalidDataItemIdList = append(invalidDataItemIdList, id)
		}
	}
	return
}

func (repo *dataViewRepoImpl) GetDataViewItems(dataViewId uuid.UUID, offset int, limit int) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	err := database.Db.Model(dvdo.DataViewItemDo{}).Select("data_item_id").Where("data_view_id = ?", dataViewId).Offset(offset).Limit(limit).Find(&idList).Error
	return idList, err
}

func (repo *dataViewRepoImpl) Updates(data dvdo.DataViewDo) error {
	return database.Db.Model(&data).Updates(data).Error
}

func (repo *dataViewRepoImpl) GetRawDataViewItems(dataViewId uuid.UUID, offset int, limit int, rawDataIdList []string, excludedAnnoViewId, includedAnnotationViewId string) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	query := database.Db.Model(dvdo.DataViewItemDo{}).Select("data_item_id").Offset(offset).Limit(limit).Where("data_view_id = ?", dataViewId)
	if len(rawDataIdList) > 0 {
		query.Where("data_item_id in (?)", rawDataIdList)
	}
	if excludedAnnoViewId != "" {
		query.Where("data_item_id not in (?)",
			database.Db.Table(dvdo.TABLE_DATA_VIEW_ITEM).
				Select("annotations.data_item_id").
				Joins("left join annotations on data_view_items.data_item_id=annotations.id").
				Where("data_view_items.data_view_id = ?", excludedAnnoViewId))
	}
	if includedAnnotationViewId != "" {
		query.Where("data_item_id in (?)",
			database.Db.Table(dvdo.TABLE_DATA_VIEW_ITEM).
				Select("annotations.data_item_id").
				Joins("left join annotations on data_view_items.data_item_id=annotations.id").
				Where("data_view_items.data_view_id = ?", includedAnnotationViewId))
	}
	err := query.Find(&idList).Error
	return idList, err
}

func (repo *dataViewRepoImpl) GetAnnotationViewItems(dataViewId uuid.UUID, offset int, limit int, rawDataIdList []string, labelId string) ([]annodo.AnnotationDo, error) {
	var annoList []annodo.AnnotationDo
	query := database.Db.Table(dvdo.TABLE_DATA_VIEW_ITEM).
		Select("annotations.id as id, annotations.data_item_id as data_item_id").
		Joins("left join annotations on data_view_items.data_item_id=annotations.id")

	if len(rawDataIdList) > 0 {
		query.Where("annotations.data_item_id in (?)", rawDataIdList)
	}
	if labelId != "" {
		query.Joins("left join raw_data_labels on annotations.id=raw_data_labels.annotation_id")
		query.Where("raw_data_labels.label_id = ?", labelId)
		query.Distinct()
	}

	query.Where("data_view_items.data_view_id = ?", dataViewId).Offset(offset).Limit(limit)
	err := query.Find(&annoList).Error
	return annoList, err
}

func (repo *dataViewRepoImpl) FilterAnnotationsByRawData(srcAnnoViewId, rawDataViewId, destAnnoViewId uuid.UUID) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err := tx.Table(dvdo.TABLE_DATA_VIEW_ITEM).Where("data_view_id = ?", destAnnoViewId).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("destAnnoViewId is not empty, it should be empty")
		}

		query := tx.Table(dvdo.TABLE_DATA_VIEW_ITEM).
			Select("annotations.id as id").
			Joins("left join annotations on data_view_items.data_item_id=annotations.id").
			Where("data_view_items.data_view_id = ?", srcAnnoViewId)
		rawDataList := tx.Table(dvdo.TABLE_DATA_VIEW_ITEM).
			Select("data_item_id").
			Where("data_view_id = ?", rawDataViewId)
		query = query.Where("annotations.data_item_id in (?)", rawDataList)
		var annoIdList []uuid.UUID
		err = query.Find(&annoIdList).Error
		if err != nil {
			return err
		}
		if len(annoIdList) == 0 {
			return nil
		}
		dataItemList := make([]dvdo.DataViewItemDo, 0)
		for _, id := range annoIdList {
			dataItemList = append(dataItemList, dvdo.DataViewItemDo{
				DataViewId: destAnnoViewId,
				DataItemId: id,
			})
		}
		return tx.Create(&dataItemList).Error
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: false})
}

func (repo *dataViewRepoImpl) MergeTo(toViewId, fromViewId uuid.UUID) (err error) {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		query := tx.Table(dvdo.TABLE_DATA_VIEW_ITEM).
			Select("data_item_id").
			Where("data_view_id = ?", fromViewId)
		var dataItemIdList []uuid.UUID
		err = query.Find(&dataItemIdList).Error
		if err != nil {
			return err
		}
		if len(dataItemIdList) == 0 {
			return nil
		}
		dataItemList := make([]dvdo.DataViewItemDo, 0)
		for _, id := range dataItemIdList {
			dataItemList = append(dataItemList, dvdo.DataViewItemDo{
				DataViewId: toViewId,
				DataItemId: id,
			})
		}
		return tx.Create(&dataItemList).Error
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: false})
}

func (repo *dataViewRepoImpl) MoveTo(srcDataViewId, dstDataViewId uuid.UUID) (err error) {
	return database.Db.Model(&dvdo.DataViewItemDo{}).Where("data_view_id = ?", srcDataViewId).Update("data_view_id", dstDataViewId).Error
}
