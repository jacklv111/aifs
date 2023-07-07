/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	manager "github.com/jacklv111/aifs/app/apigin/manager/annotation-template"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	datamodule "github.com/jacklv111/aifs/pkg/data-view/data-module"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	"github.com/jacklv111/aifs/pkg/data-view/data-view/bo"
	"github.com/jacklv111/aifs/pkg/data-view/data-view/repo"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
)

type dataViewMgrImpl struct {
}

var DataViewMgr dataViewMgrInterface

func init() {
	DataViewMgr = &dataViewMgrImpl{}
}

func (mgr *dataViewMgrImpl) Create(req openapi.CreateDataViewRequest) (result openapi.CreateDataViewSuccessResp, err error) {
	dataViewBo, err := bo.BuildFromCreateDataViewRequest(req)
	if err != nil {
		return result, err
	}

	var id uuid.UUID
	if id, err = dataViewBo.Create(); err != nil {
		return
	}
	result.DataViewId = id.String()
	return
}

func (mgr *dataViewMgrImpl) GetDetailsById(id uuid.UUID) (result openapi.DataViewDetails, err error) {
	dataViewBo, err := bo.BuildWithId(id)
	if err != nil {
		return result, err
	}

	var details dvvb.DataViewDetails
	details, err = dataViewBo.GetDetails()
	if err != nil {
		return
	}

	var annoType string
	if details.ViewType == datamodule.ANNOTATION {
		annoType, err = manager.AnnotationTemplateMgr.GetTypeById(details.AnnotationTemplateId)
		if err != nil {
			return result, err
		}
	}
	result = assembleToDataViewDetails(details, annoType)
	return
}

func (mgr *dataViewMgrImpl) GetStatisticsById(id uuid.UUID) (result openapi.DataViewStatistics, err error) {
	dataViewBo, err := bo.BuildWithId(id)
	if err != nil {
		return result, err
	}

	var statistics dvvb.DataViewStatistics
	switch dv := dataViewBo.(type) {
	case bo.RawDataViewBoInterface:
		statistics, err = dv.GetStatistics()
	case bo.AnnotationViewBoInterface:
		statistics, err = dv.GetStatistics()
	default:
		err = fmt.Errorf("data view type %s is not supported to get statistics", dataViewBo.GetViewType())
	}
	result = assembleToDataViewStatistics(statistics)
	return
}

func (mgr *dataViewMgrImpl) Delete(id uuid.UUID) error {
	dataViewBo, err := bo.BuildWithId(id)
	if err != nil {
		return err
	}

	return dataViewBo.Delete()
}

func (mgr *dataViewMgrImpl) HardDelete(id uuid.UUID) error {
	dataViewBo, err := bo.BuildWithId(id)
	if err != nil {
		return err
	}

	if dataZipView, ok := dataViewBo.(bo.DatasetZipViewBoInterface); ok {
		if dataZipView.IsCompleted() {
			return dataZipView.HardDelete()
		}
		viewList := dataZipView.GetDataViewIdList()
		for _, id := range viewList {

			bo, err := bo.BuildWithId(id)
			if err != nil {
				return err
			}
			err = bo.HardDelete()
			if err != nil {
				return err
			}
		}
		if dataZipView.GetAnnotationTemplateId() == uuid.Nil {
			return nil
		}
		return manager.AnnotationTemplateMgr.DeleteById(dataZipView.GetAnnotationTemplateId())
	} else {
		return dataViewBo.HardDelete()
	}
}

func (mgr *dataViewMgrImpl) GetList(offset int, limit int, dataViewIdList []string, dataViewName string) ([]openapi.DataViewListItem, error) {
	options := dvvb.DataViewListQueryOptions{Offset: offset, Limit: limit, DataViewIdList: dataViewIdList, DataViewName: dataViewName}
	doList, err := repo.DataViewRepo.GetList(options)
	if err != nil {
		return nil, err
	}
	// get annotation type for annotation data view
	annoTypeMap := make(map[uuid.UUID]string, 0)
	for _, data := range doList {
		if data.ViewType == datamodule.ANNOTATION {
			annoTypeMap[data.ID], err = manager.AnnotationTemplateMgr.GetTypeById(data.AnnotationTemplateId)
			if err != nil {
				return nil, err
			}
		}
	}
	return assembleToDataViewListItemList(doList, annoTypeMap), nil
}

func (mgr *dataViewMgrImpl) DeleteDataViewItem(dataViewId uuid.UUID, dataViewItemIdList []uuid.UUID) error {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return err
	}

	err = dataViewBo.DeleteDataItem(dataViewItemIdList)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *dataViewMgrImpl) UploadAnnotationToDataView(dataViewId uuid.UUID, fileMeta io.Reader, dataFileMap map[string]io.ReadSeeker, dataFileNameMap map[string]string) error {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return err
	}
	uploadAnnotationParams := basicvb.UploadAnnotationParam{
		FileMeta:        fileMeta,
		DataFileMap:     dataFileMap,
		DataFileNameMap: dataFileNameMap,
		RawDataIdChecker: func(dataItemIdList []uuid.UUID) error {
			res, err := repo.DataViewRepo.GetInvalidDataItems(dataViewBo.(bo.AnnotationViewBoInterface).GetRelateRawDataViewId(), dataItemIdList)
			if err != nil {
				return err
			}
			if len(res) != 0 {
				return fmt.Errorf("invalid raw data id %v", res)
			}
			return nil
		},
	}
	return dataViewBo.(bo.AnnotationViewBoInterface).UploadAnnotationV2(uploadAnnotationParams)
}

func (mgr *dataViewMgrImpl) UploadRawDataToDataView(dataViewId uuid.UUID, fileMeta io.Reader, dataFileMap map[string]io.ReadSeeker) error {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return err
	}
	uploadRawDataParams := basicvb.UploadRawDataParam{
		FileMeta:    fileMeta,
		DataFileMap: dataFileMap,
	}
	return dataViewBo.(bo.RawDataViewBoInterface).UploadRawDataV2(uploadRawDataParams)
}

func (mgr *dataViewMgrImpl) UploadModelDataToDataView(dataViewId uuid.UUID, pairs map[string]string, dataFileMap map[string]io.ReadSeeker) error {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return err
	}
	uploadModelDataParams := basicvb.UploadModelParams{
		Pairs:       pairs,
		DataFileMap: dataFileMap,
	}
	return dataViewBo.(bo.ModelViewBoInterface).UploadModelData(uploadModelDataParams)
}

func (mgr *dataViewMgrImpl) UploadDatasetZipToDataView(dataViewId uuid.UUID, file io.Reader, fileName string) error {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return err
	}
	uploadDatasetZipParams := basicvb.UploadDatasetZipParams{
		File:     file,
		FileName: fileName,
	}
	return dataViewBo.(bo.DatasetZipViewBoInterface).UploadDatasetZip(uploadDatasetZipParams)
}

func (mgr *dataViewMgrImpl) DownloadAll(dataViewIdList []uuid.UUID) (resp openapi.DownloadResponse, err error) {
	// validate
	invalidIds, err := repo.DataViewRepo.GetInvalidId(dataViewIdList)
	if err != nil {
		return
	}
	if len(invalidIds) > 0 {
		err = fmt.Errorf("invalid data view ids: [%s]", invalidIds)
		return
	}

	rootDir := filepath.Join("/Users/lvyubin/work/test/download", uuid.New().String())

	for _, id := range dataViewIdList {
		var dataViewBo bo.DataViewBoInterface
		dataViewBo, err = bo.BuildWithId(id)
		if err != nil {
			return resp, err
		}
		err = dataViewBo.DownloadAll(dvvb.DownloadQuery{Directory: rootDir})
		if err != nil {
			return
		}
	}
	return openapi.DownloadResponse{Directory: rootDir}, nil
}

func (mgr *dataViewMgrImpl) GetAllAnnotationLocations(dataViewId uuid.UUID) (resp openapi.AnnotationViewLocations, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}

	res, err := dataViewBo.(bo.AnnotationViewBoInterface).GetAllAnnotationLocations()

	if err != nil {
		return resp, err
	}

	return assembleAnnotationLocations(res), nil
}

func (mgr *dataViewMgrImpl) GetAllRawDataLocations(dataViewId uuid.UUID) (resp openapi.RawDataViewLocations, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}

	res, err := dataViewBo.(bo.RawDataViewBoInterface).GetAllRawDataLocations()

	if err != nil {
		return resp, err
	}

	return assembleRawDataLocations(res), nil
}

func (mgr *dataViewMgrImpl) GetAllAnnotationData(dataViewId uuid.UUID) (resp openapi.AnnotationViewData, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}

	res, err := dataViewBo.(bo.AnnotationViewBoInterface).GetAllAnnoData()

	if err != nil {
		return resp, err
	}

	return assembleAnnotationData(res), nil
}

func (mgr *dataViewMgrImpl) GetRawDataHashList(dataViewId uuid.UUID, offset int, limit int) (resp []openapi.RawDataHashListInner, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}

	res, err := dataViewBo.(bo.RawDataViewBoInterface).GetHashList(offset, limit)

	if err != nil {
		return resp, err
	}
	return assembleRawDataHashList(res), nil
}

func (mgr *dataViewMgrImpl) GetAllModelDataLocations(dataViewId uuid.UUID) (resp openapi.ModelDataViewLocations, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}

	res, err := dataViewBo.(bo.ModelViewBoInterface).GetModelDataLocations()

	if err != nil {
		return resp, err
	}

	return assembleModelDataLocations(res), nil
}

func (mgr *dataViewMgrImpl) UpdateDatasetZipView(dataViewId uuid.UUID, req openapi.UpdateDatasetZipRequest) (err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return err
	}
	params := dvvb.UpdateDatasetZipParams{
		Progress:              float64(req.Progress),
		Status:                req.Status,
		RawDataViewId:         req.RawDataViewId,
		AnnotationViewId:      req.AnnotationViewId,
		TrainRawDataViewId:    req.TrainRawDataViewId,
		TrainAnnotationViewId: req.TrainAnnotationViewId,
		ValRawDataViewId:      req.ValRawDataViewId,
		ValAnnotationViewId:   req.ValAnnotationViewId,
		AnnotationTemplateId:  req.AnnotationTemplateId,
	}

	return dataViewBo.(bo.DatasetZipViewBoInterface).UpdateDatasetZipView(params)
}

func (mgr *dataViewMgrImpl) GetDatasetZipLocationInDataView(dataViewId uuid.UUID) (resp openapi.DatasetZipLocation, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}

	res, err := dataViewBo.(bo.DatasetZipViewBoInterface).GetDatasetZipLocation()

	if err != nil {
		return resp, err
	}

	return assembleDatasetZipLocations(res), nil
}

func (mgr *dataViewMgrImpl) UploadFileToDataView(dataViewId uuid.UUID, file io.Reader, fileName string) error {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return err
	}
	uploadArtifactFileParams := basicvb.UploadArtifactFileParams{
		File:     file,
		FileName: fileName,
	}
	return dataViewBo.(bo.ArtifactViewBoInterface).UploadArtifactFile(uploadArtifactFileParams)
}

func (mgr *dataViewMgrImpl) GetArtifactLocationsInDataView(dataViewId uuid.UUID) (resp openapi.ArtifactLocations, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}

	res, err := dataViewBo.(bo.ArtifactViewBoInterface).GetArtifactLocations()

	if err != nil {
		return resp, err
	}

	return assembleArtifactLocations(res), nil
}

func (mgr *dataViewMgrImpl) GetRawDataList(dataViewId uuid.UUID, offset int, limit int, rawDataIdList []string, excludedAnnoViewId, includedAnnotationViewId string) (resp openapi.GetRawDataInDataView200Response, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}
	if _, ok := dataViewBo.(bo.RawDataViewBoInterface); !ok {
		return resp, fmt.Errorf("data view is not raw data view")
	}
	if excludedAnnoViewId != "" && includedAnnotationViewId != "" {
		return resp, fmt.Errorf("excludedAnnoViewId and includedAnnotationViewId can not be set at the same time")
	}
	resp.RawDataType = openapi.RawDataType(dataViewBo.(bo.RawDataViewBoInterface).GetRawDataType())
	res, err := dataViewBo.(bo.RawDataViewBoInterface).GetRawDataList(offset, limit, rawDataIdList, excludedAnnoViewId, includedAnnotationViewId)
	if err != nil {
		return resp, err
	}
	for _, data := range res {
		resp.RawDataList = append(resp.RawDataList, openapi.RawDataListItem{
			RawDataId: data.RawDataId,
			Url:       data.Url,
		})
	}
	return resp, nil
}

func (mgr *dataViewMgrImpl) GetAnnotationList(dataViewId uuid.UUID, offset int, limit int, rawDataIdList []string, labelId string) (resp openapi.GetAnnotationsInDataView200Response, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}
	if _, ok := dataViewBo.(bo.AnnotationViewBoInterface); !ok {
		return resp, fmt.Errorf("data view is not annotation view")
	}
	resp.AnnotationTemplateId = dataViewBo.(bo.AnnotationViewBoInterface).GetAnnotationTemplateId().String()
	resp.AnnotationTemplateType, err = manager.AnnotationTemplateMgr.GetTypeById(dataViewBo.(bo.AnnotationViewBoInterface).GetAnnotationTemplateId())
	if err != nil {
		return resp, err
	}
	res, err := dataViewBo.(bo.AnnotationViewBoInterface).GetAnnotationList(offset, limit, rawDataIdList, labelId)
	if err != nil {
		return resp, err
	}
	for _, data := range res {
		resp.AnnotationList = append(resp.AnnotationList, openapi.AnnotationListItem{
			RawDataId:  data.RawDataId,
			DataItemId: data.DataItemId,
			Url:        data.Url,
			Labels:     data.Labels,
		})
	}
	return resp, nil
}

func (mgr *dataViewMgrImpl) DivideDataView(dataViewId uuid.UUID, req []openapi.DivideRawDataDataViewRequestInner) (resp []openapi.DivideRawDataDataViewResponseInner, err error) {
	dataViewBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return resp, err
	}
	if _, ok := dataViewBo.(bo.RawDataViewBoInterface); !ok {
		return resp, fmt.Errorf("data view is not raw data view")
	}
	params := dvvb.DivideRawDataViewParams{}
	for _, data := range req {
		params.RawDataViewParamList = append(params.RawDataViewParamList, dvvb.EachRawDataViewParam{
			Name:        data.Name,
			Description: data.Description,
			Ratio:       data.Ratio,
		})
	}
	res, err := dataViewBo.(bo.RawDataViewBoInterface).DivideDataView(params)
	if err != nil {
		return resp, err
	}
	for _, data := range res.RawDataViewResultList {
		resp = append(resp, openapi.DivideRawDataDataViewResponseInner{
			Name:       data.Name,
			DataViewId: data.DataViewId,
			ItemCount:  data.ItemCount,
		})
	}
	return resp, nil
}

func (mgr *dataViewMgrImpl) FilterAnnotationsInDataView(annotationViewId uuid.UUID, rawDataViewId uuid.UUID) (resp openapi.FilterAnnotationsInDataViewResponse, err error) {
	dataViewBo, err := bo.BuildWithId(annotationViewId)
	if err != nil {
		return resp, err
	}
	if _, ok := dataViewBo.(bo.AnnotationViewBoInterface); !ok {
		return resp, fmt.Errorf("data view is not annotation view")
	}
	res, err := dataViewBo.(bo.AnnotationViewBoInterface).FilterAnnotationsInDataView(rawDataViewId)
	if err != nil {
		return resp, err
	}
	resp.AnnotationViewId = res.String()
	return resp, nil
}

func (mgr *dataViewMgrImpl) MergeDataViews(req openapi.MergeDataViewsRequest) (resp openapi.MergeDataViewsSuccessResp, err error) {
	dataViewIdList := make([]uuid.UUID, 0)
	for _, dataViewId := range req.DataViewIdList {
		dataViewIdList = append(dataViewIdList, uuid.MustParse(dataViewId))
	}
	resBo, err := bo.BuildFromBoType(dataViewIdList[0], req.Name, req.Description)
	if err != nil {
		return resp, err
	}
	for _, dataViewId := range dataViewIdList {
		err = resBo.MergeDataItems(dataViewId)
		if err != nil {
			return resp, err
		}
	}
	id, err := resBo.Create()
	if err != nil {
		return resp, err
	}
	resp.DataViewId = id.String()
	return resp, nil
}

func (mgr *dataViewMgrImpl) MergeDataViewsToCurrent(dataViewId uuid.UUID, req openapi.MergeDataViewsRequest) error {
	dataViewIdList := make([]uuid.UUID, 0)
	for _, dataViewId := range req.DataViewIdList {
		dataViewIdList = append(dataViewIdList, uuid.MustParse(dataViewId))
	}
	curBo, err := bo.BuildWithId(dataViewId)
	if err != nil {
		return err
	}

	for _, dataViewId := range dataViewIdList {
		err = curBo.MergeDataItems(dataViewId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mgr *dataViewMgrImpl) MoveDataViewItems(srcDataViewId, dstDataViewId string) error {
	return repo.DataViewRepo.MoveTo(uuid.MustParse(srcDataViewId), uuid.MustParse(dstDataViewId))
}
