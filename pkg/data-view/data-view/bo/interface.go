/*
 * Created on Sat Jul 08 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"github.com/google/uuid"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type DataViewBoInterface interface {
	GetId() uuid.UUID

	GetViewType() string

	// 创建一个 data view
	Create() (uuid.UUID, error)

	// 软删除 data view
	Delete() error

	// 硬删除
	HardDelete() error

	// 获取 data view 的详情信息
	GetDetails() (dvvb.DataViewDetails, error)

	// 删除 data view 下的 data item
	//  @param idList 需要删除的 data item id 列表
	//  @return error
	DeleteDataItem(idList []uuid.UUID) error

	// MergeDataItems it will check if data views are the same type
	//  @param other
	//  @return error
	MergeDataItems(other uuid.UUID) error
}

type RawDataViewBoInterface interface {
	DataViewBoInterface

	UploadRawData(input basicvb.UploadRawDataParam) error

	GetAllRawDataLocations() (dvvb.RawDataLocationResult, error)

	GetHashList(offset int, limit int) (result []basicdo.IdHash, err error)

	GetRawDataType() string

	GetRawDataList(offset int, limit int, rawDataIdList []string, excludedAnnoViewId, includedAnnotationViewId string) ([]dvvb.RawDataItem, error)

	DivideDataView(dvvb.DivideRawDataViewParams) (dvvb.DivideRawDataViewResult, error)

	CreateDataViewItems(itemIdList []uuid.UUID) error

	GetStatistics() (dvvb.DataViewStatistics, error)
}

type AnnotationViewBoInterface interface {
	DataViewBoInterface

	UploadAnnotations(input basicvb.UploadAnnotationParam) error

	GetAllAnnotationLocations() (dvvb.AnnotationLocationResult, error)

	GetAllAnnoData() (result dvvb.AnnotationData, err error)

	GetRelateRawDataViewId() uuid.UUID

	GetAnnotationTemplateId() uuid.UUID

	GetAnnotationList(offset int, limit int, rawDataIdList []string, labelId string) (annoList []dvvb.AnnotationItem, err error)

	FilterAnnotationsInDataView(rawDataViewId uuid.UUID) (annotationViewId uuid.UUID, err error)

	GetStatistics() (dvvb.DataViewStatistics, error)
}

type ModelViewBoInterface interface {
	DataViewBoInterface

	UploadModelData(input basicvb.UploadModelParams) (err error)

	GetModelDataLocations() (result dvvb.ModelLocationResult, err error)
}

type DatasetZipViewBoInterface interface {
	DataViewBoInterface

	UploadDatasetZip(input basicvb.UploadDatasetZipParams) (err error)

	UpdateDatasetZipView(params dvvb.UpdateDatasetZipParams) error

	GetDatasetZipLocation() (result dvvb.DatasetZipLocationResult, err error)

	IsCompleted() bool

	GetDataViewIdList() []uuid.UUID

	GetAnnotationTemplateId() uuid.UUID
}

type ArtifactViewBoInterface interface {
	DataViewBoInterface

	UploadArtifactFile(input basicvb.UploadArtifactFileParams) (err error)

	GetArtifactLocations() (result dvvb.ArtifactLocationResult, err error)
}
