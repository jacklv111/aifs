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
	"fmt"

	"github.com/google/uuid"
	datamodule "github.com/jacklv111/aifs/pkg/data-view/data-module"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	zipbo "github.com/jacklv111/aifs/pkg/data-view/data-module/dataset-zip/bo"
	zipconst "github.com/jacklv111/aifs/pkg/data-view/data-module/dataset-zip/constant"
	dvrepo "github.com/jacklv111/aifs/pkg/data-view/data-view/repo"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
	storemgr "github.com/jacklv111/aifs/pkg/store/manager"
	storevb "github.com/jacklv111/aifs/pkg/store/value-object"
)

type datasetZipViewBo struct {
	dataViewBo
}

func (bo *datasetZipViewBo) Create() (uuid.UUID, error) {
	if !zipconst.IsZipFormat(bo.dataViewDo.ZipFormat.String) {
		return uuid.Nil, fmt.Errorf("invalid zip format %s, zip format should be %v", bo.dataViewDo.ZipFormat.String, zipconst.GetZipFormatList())
	}

	return bo.dataViewBo.Create()
}

func (bo *datasetZipViewBo) UploadDatasetZip(input basicvb.UploadDatasetZipParams) (err error) {
	if err = bo.loadDataViewDo(); err != nil {
		return
	}
	dataItemList, err := dvrepo.DataViewRepo.GetAllDataViewItems(bo.dataViewDo.ID)
	if err != nil {
		return err
	}
	if len(dataItemList) > 0 {
		return fmt.Errorf("you have already uploaded dataset zip file")
	}

	datasetZipBo := zipbo.BuildFromBuffer(input.File, input.FileName)

	// save meta
	err = datasetZipBo.Create()
	if err != nil {
		return
	}

	// save data
	storeParams, err := getStoreParamRemote(datasetZipBo)

	if err != nil {
		return err
	}

	err = storemgr.StoreMgr.Upload(storeParams)
	if err != nil {
		return err
	}

	// save data view items
	return dvrepo.DataViewRepo.CreateDataViewItemsIgnoreConflict(bo.dataViewDo.ID, []uuid.UUID{datasetZipBo.ID})
}

func (bo *datasetZipViewBo) UpdateDatasetZipView(params dvvb.UpdateDatasetZipParams) error {
	bo.dataViewDo.Progress = params.Progress
	bo.dataViewDo.Status = sql.NullString{String: params.Status, Valid: true}
	bo.dataViewDo.TrainRawDataViewId = sql.NullString{String: params.TrainRawDataViewId, Valid: params.TrainRawDataViewId != ""}
	bo.dataViewDo.TrainAnnotationViewId = sql.NullString{String: params.TrainAnnotationViewId, Valid: params.TrainAnnotationViewId != ""}
	bo.dataViewDo.ValRawDataViewId = sql.NullString{String: params.ValRawDataViewId, Valid: params.TrainRawDataViewId != ""}
	bo.dataViewDo.ValAnnotationViewId = sql.NullString{String: params.ValAnnotationViewId, Valid: params.TrainAnnotationViewId != ""}
	bo.dataViewDo.RawDataViewId = sql.NullString{String: params.RawDataViewId, Valid: params.RawDataViewId != ""}
	bo.dataViewDo.AnnotationViewId = sql.NullString{String: params.AnnotationViewId, Valid: params.AnnotationViewId != ""}
	if params.AnnotationTemplateId != "" {
		bo.dataViewDo.AnnotationTemplateId = uuid.MustParse(params.AnnotationTemplateId)
	}
	return dvrepo.DataViewRepo.Updates(bo.dataViewDo)
}

func (bo *datasetZipViewBo) GetDatasetZipLocation() (result dvvb.DatasetZipLocationResult, err error) {
	result.DataViewId = bo.dataViewDo.ID.String()
	result.ViewType = bo.dataViewDo.ViewType
	result.DataItemDoList, result.LocationMap, err = bo.getDataItemAndLocations()
	return
}

func (bo *datasetZipViewBo) IsCompleted() bool {
	return bo.dataViewDo.Progress+(1e-8) >= 100.0
}

func (bo *datasetZipViewBo) GetDataViewIdList() []uuid.UUID {
	res := make([]uuid.UUID, 0)
	if bo.dataViewDo.TrainRawDataViewId.Valid {
		res = append(res, uuid.MustParse(bo.dataViewDo.TrainRawDataViewId.String))
	}

	if bo.dataViewDo.TrainAnnotationViewId.Valid {
		res = append(res, uuid.MustParse(bo.dataViewDo.TrainAnnotationViewId.String))
	}

	if bo.dataViewDo.ValRawDataViewId.Valid {
		res = append(res, uuid.MustParse(bo.dataViewDo.ValRawDataViewId.String))
	}

	if bo.dataViewDo.ValAnnotationViewId.Valid {
		res = append(res, uuid.MustParse(bo.dataViewDo.ValAnnotationViewId.String))
	}
	res = append(res, bo.dataViewDo.ID)
	return res
}

func (bo *datasetZipViewBo) GetAnnotationTemplateId() uuid.UUID {
	return bo.dataViewDo.AnnotationTemplateId
}

func getStoreParamRemote(bo *zipbo.DatasetZipBo) (storevb.UploadParams, error) {
	var params storevb.UploadParams
	params.DataType = datamodule.DATASET_ZIP
	params.AddItem(bo.GetId(), bo.GetReader(), bo.GetName())
	return params, nil
}
