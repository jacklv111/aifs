/*
 * Created on Sat Jul 08 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"github.com/google/uuid"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	dvdo "github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type dataViewRepoInterface interface {
	ExistsById(dataViewId uuid.UUID) (bool, error)
	Create(dvdo.DataViewDo) error
	GetList(dvvb.DataViewListQueryOptions) ([]dvdo.DataViewDo, error)
	GetById(dataViewId uuid.UUID) (dvdo.DataViewDo, error)
	SoftDelete(annoTempId uuid.UUID) error
	HardDelete(annoTempId uuid.UUID) error
	GetDataViewItemCount(dataViewId uuid.UUID) (count int64, err error)
	DeleteDataViewItem(dataViewId uuid.UUID, dataViewItemIdList []uuid.UUID) error
	CreateDataViewItemsIgnoreConflict(dataViewId uuid.UUID, itemIdList []uuid.UUID) error
	CreateAnnotationDataViewItems(itemList []dvdo.DataViewItemDo, annotationTemplateId uuid.UUID) error
	GetAllDataViewItems(dataViewId uuid.UUID) ([]uuid.UUID, error)
	GetDataViewItems(dataViewId uuid.UUID, offset int, limit int) ([]uuid.UUID, error)
	GetInvalidId(dataViewIdList []uuid.UUID) ([]uuid.UUID, error)
	GetInvalidDataItems(dataViewId uuid.UUID, dataItemIdList []uuid.UUID) (invalidDataItemIdList []uuid.UUID, err error)
	Updates(data dvdo.DataViewDo) error
	GetRawDataViewItems(dataViewId uuid.UUID, offset int, limit int, rawDataIdList []string, excludedAnnoViewId, includedAnnotationViewId string) ([]uuid.UUID, error)
	GetAnnotationViewItems(dataViewId uuid.UUID, offset int, limit int, rawDataIdList []string, labelId string) ([]annodo.AnnotationDo, error)
	FilterAnnotationsByRawData(srcAnnoViewId, rawDataViewId, destAnnoViewId uuid.UUID) error
	MergeTo(toViewId, fromViewId uuid.UUID) (err error)
	MoveTo(srcDataViewId, dstDataViewId uuid.UUID) (err error)
}
