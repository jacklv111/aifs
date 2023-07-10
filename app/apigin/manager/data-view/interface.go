/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	"io"

	"github.com/google/uuid"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type dataViewMgrInterface interface {
	Create(openapi.CreateDataViewRequest) (openapi.CreateDataViewSuccessResp, error)
	GetDetailsById(uuid.UUID) (result openapi.DataViewDetails, err error)
	GetStatisticsById(id uuid.UUID) (result openapi.DataViewStatistics, err error)
	Delete(uuid.UUID) error
	GetList(offset int, limit int, dataViewIdList []string, dataViewName string) ([]openapi.DataViewListItem, error)
	DeleteDataViewItem(dataViewId uuid.UUID, dataViewItemIdList []uuid.UUID) error
	UploadRawDataToDataView(dataViewId uuid.UUID, fileMeta io.Reader, dataFileMap map[string]io.ReadSeeker) error
	UploadAnnotationToDataView(dataViewId uuid.UUID, fileMeta io.Reader, dataFileMap map[string]io.ReadSeeker, dataFileNameMap map[string]string) error
	HardDelete(id uuid.UUID) error
	GetAllAnnotationLocations(dataViewId uuid.UUID) (resp openapi.AnnotationViewLocations, err error)
	GetAllRawDataLocations(dataViewId uuid.UUID) (resp openapi.RawDataViewLocations, err error)
	GetAllAnnotationData(dataViewId uuid.UUID) (resp openapi.AnnotationViewData, err error)
	GetAllModelDataLocations(dataViewId uuid.UUID) (resp openapi.ModelDataViewLocations, err error)
	GetRawDataHashList(dataViewId uuid.UUID, offset int, limit int) (resp []openapi.RawDataHashListInner, err error)
	UploadModelDataToDataView(dataViewId uuid.UUID, pairs map[string]string, dataFileMap map[string]io.ReadSeeker) error
	UploadDatasetZipToDataView(dataViewId uuid.UUID, file io.Reader, fileName string) error
	UpdateDatasetZipView(dataViewId uuid.UUID, req openapi.UpdateDatasetZipRequest) (err error)
	GetDatasetZipLocationInDataView(dataViewId uuid.UUID) (resp openapi.DatasetZipLocation, err error)
	UploadFileToDataView(dataViewId uuid.UUID, file io.Reader, fileName string) error
	GetArtifactLocationsInDataView(dataViewId uuid.UUID) (resp openapi.ArtifactLocations, err error)
	GetRawDataList(dataViewId uuid.UUID, offset int, limit int, rawDataIdList []string, excludedAnnoViewId, includedAnnotationViewId string) (resp openapi.GetRawDataInDataView200Response, err error)
	GetAnnotationList(dataViewId uuid.UUID, offset int, limit int, rawDataIdList []string, labelId string) (resp openapi.GetAnnotationsInDataView200Response, err error)
	DivideDataView(dataViewId uuid.UUID, req []openapi.DivideRawDataDataViewRequestInner) (resp []openapi.DivideRawDataDataViewResponseInner, err error)
	FilterAnnotationsInDataView(annotationViewId uuid.UUID, rawDataViewId uuid.UUID) (resp openapi.FilterAnnotationsInDataViewResponse, err error)
	MergeDataViews(openapi.MergeDataViewsRequest) (resp openapi.MergeDataViewsSuccessResp, err error)
	MergeDataViewsToCurrent(dataViewId uuid.UUID, req openapi.MergeDataViewsRequest) error
	MoveDataViewItems(srcDataViewId, dstDataViewId string) error
}
