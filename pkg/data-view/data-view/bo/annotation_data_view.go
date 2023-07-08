/*
 * Created on Sat Jul 08 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"fmt"

	"github.com/google/uuid"
	annotempmgr "github.com/jacklv111/aifs/app/apigin/manager/annotation-template"
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	annodo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/repo"
	annosvc "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/service"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	dvdo "github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	dvrepo "github.com/jacklv111/aifs/pkg/data-view/data-view/repo"
	vb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
	storemgr "github.com/jacklv111/aifs/pkg/store/manager"
)

type annotationViewBo struct {
	dataViewBo
}

func (bo *annotationViewBo) Create() (uuid.UUID, error) {
	if bo.dataViewDo.RelatedDataViewId == uuid.Nil {
		return uuid.Nil, fmt.Errorf("annotation type data view must have a related raw-data data view")
	}
	if bo.dataViewDo.AnnotationTemplateId == uuid.Nil {
		return uuid.Nil, fmt.Errorf("annotation type data view must have a annotation template id")
	}
	exists, err := dvrepo.DataViewRepo.ExistsById(bo.dataViewDo.RelatedDataViewId)
	if err != nil {
		return uuid.Nil, err
	}
	if !exists {
		return uuid.Nil, fmt.Errorf("related raw-data data view %s doesn't exist", bo.dataViewDo.ID)
	}

	return bo.dataViewBo.Create()
}

func (bo *annotationViewBo) UploadAnnotations(input basicvb.UploadAnnotationParam) (err error) {
	if err = bo.loadDataViewDo(); err != nil {
		return
	}

	// upload annotations
	var dataItemIdList []uuid.UUID
	dataItemIdList, err = annosvc.UploadAnnotations(input, bo.dataViewDo.AnnotationTemplateId)
	if err != nil {
		return
	}

	// save data view items
	var dataViewItemDoList []dvdo.DataViewItemDo
	for _, data := range dataItemIdList {
		dataViewItemDoList = append(dataViewItemDoList, dvdo.DataViewItemDo{DataViewId: bo.dataViewDo.ID, DataItemId: data})
	}
	err = dvrepo.DataViewRepo.CreateAnnotationDataViewItems(dataViewItemDoList, bo.dataViewDo.AnnotationTemplateId)

	return
}

func (bo *annotationViewBo) GetAllAnnotationLocations() (result vb.AnnotationLocationResult, err error) {
	result.DataViewId = bo.dataViewDo.ID.String()
	result.ViewType = bo.dataViewDo.ViewType
	result.AnnotationTemplateId = bo.dataViewDo.AnnotationTemplateId.String()
	dataItemIdList, err := dvrepo.DataViewRepo.GetAllDataViewItems(bo.dataViewDo.ID)
	if err != nil {
		return result, err
	}
	var annoList []annodo.AnnotationDo
	result.DataItemDoList, annoList, err = repo.AnnotationRepo.GetByIdList(dataItemIdList)
	if err != nil {
		return result, err
	}
	result.AnnoDoMap = make(map[uuid.UUID]annodo.AnnotationDo)
	for _, data := range annoList {
		result.AnnoDoMap[data.ID] = data
	}

	result.LocationMap, err = storemgr.StoreMgr.GetByIdList(dataItemIdList)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (bo *annotationViewBo) GetAllAnnoData() (result vb.AnnotationData, err error) {
	// validate
	annoTempType, err := annotempmgr.AnnotationTemplateMgr.GetTypeById(bo.dataViewDo.AnnotationTemplateId)
	if err != nil {
		return result, err
	}
	if annoTempType != annotationtemplatetype.CATEGORY && annoTempType != annotationtemplatetype.OCR {
		return result, fmt.Errorf("annotation template type %s can't get annotation data", annoTempType)
	}

	result.DataViewId = bo.dataViewDo.ID.String()
	result.ViewType = bo.dataViewDo.ViewType
	result.AnnotationTemplateId = bo.dataViewDo.AnnotationTemplateId.String()

	dataItemIdList, err := dvrepo.DataViewRepo.GetAllDataViewItems(bo.dataViewDo.ID)
	if err != nil {
		return result, err
	}
	var annoList []annodo.AnnotationDo
	result.DataItemDoList, annoList, err = repo.AnnotationRepo.GetByIdList(dataItemIdList)
	if err != nil {
		return result, err
	}
	result.AnnoDoMap = make(map[uuid.UUID]annodo.AnnotationDo)
	for _, data := range annoList {
		result.AnnoDoMap[data.ID] = data
	}

	rawDataLabelList, err := repo.AnnotationRepo.GetRawDataLabelByIdList(dataItemIdList)
	if err != nil {
		return result, err
	}

	result.RawDataLabelMap = make(map[uuid.UUID][]annodo.RawDataLabelDo)
	for _, rawDataLabel := range rawDataLabelList {
		result.RawDataLabelMap[rawDataLabel.AnnotationId] = append(result.RawDataLabelMap[rawDataLabel.AnnotationId], rawDataLabel)
	}
	return result, nil
}

func (bo *annotationViewBo) GetRelateRawDataViewId() uuid.UUID {
	return bo.dataViewDo.RelatedDataViewId
}

func (bo *annotationViewBo) GetAnnotationTemplateId() uuid.UUID {
	return bo.dataViewDo.AnnotationTemplateId
}

func (bo *annotationViewBo) GetAnnotationList(offset int, limit int, rawDataIdList []string, labelId string) (annoList []vb.AnnotationItem, err error) {
	dataItemList, err := dvrepo.DataViewRepo.GetAnnotationViewItems(bo.dataViewDo.ID, offset, limit, rawDataIdList, labelId)
	if err != nil {
		return nil, err
	}
	idList := make([]uuid.UUID, 0)

	for _, data := range dataItemList {
		idList = append(idList, data.ID)
	}

	idUrlMap, err := storemgr.StoreMgr.GetUrlListByIdList(idList)
	if err != nil {
		return nil, err
	}

	rawDataLabelList, err := repo.AnnotationRepo.GetRawDataLabelByIdList(idList)
	if err != nil {
		return nil, err
	}
	rawDataLabelMap := make(map[uuid.UUID][]string)
	for _, rawDataLabel := range rawDataLabelList {
		rawDataLabelMap[rawDataLabel.AnnotationId] = append(rawDataLabelMap[rawDataLabel.AnnotationId], rawDataLabel.LabelId.String())
	}

	annoList = make([]vb.AnnotationItem, len(idList))
	for idx, dataItem := range dataItemList {
		annoList[idx].RawDataId = dataItem.DataItemId.String()
		annoList[idx].Url = idUrlMap[dataItem.ID]
		annoList[idx].Labels = rawDataLabelMap[dataItem.ID]
		annoList[idx].DataItemId = dataItem.ID.String()
	}

	return
}

func (bo *annotationViewBo) FilterAnnotationsInDataView(rawDataViewId uuid.UUID) (annotationViewId uuid.UUID, err error) {
	destAnnoView := annotationViewBo{
		dataViewBo: dataViewBo{
			dataViewDo: dvdo.DataViewDo{
				ID:                   uuid.New(),
				Name:                 bo.dataViewDo.Name,
				ViewType:             bo.dataViewDo.ViewType,
				RelatedDataViewId:    rawDataViewId,
				AnnotationTemplateId: bo.dataViewDo.AnnotationTemplateId,
			},
		},
	}
	destAnnoViewId, err := destAnnoView.Create()
	if err != nil {
		return uuid.Nil, err
	}
	err = dvrepo.DataViewRepo.FilterAnnotationsByRawData(bo.dataViewDo.ID, rawDataViewId, destAnnoViewId)
	if err != nil {
		return uuid.Nil, err
	}
	return destAnnoViewId, nil
}

func (bo *annotationViewBo) GetStatistics() (vb.DataViewStatistics, error) {
	count, err := dvrepo.DataViewRepo.GetDataViewItemCount(bo.dataViewDo.ID)
	if err != nil {
		return vb.DataViewStatistics{}, err
	}
	dataItemIdList, err := dvrepo.DataViewRepo.GetAllDataViewItems(bo.dataViewDo.ID)
	if err != nil {
		return vb.DataViewStatistics{}, err
	}
	rawDataLabelList, err := repo.AnnotationRepo.GetRawDataLabelByIdList(dataItemIdList)
	if err != nil {
		return vb.DataViewStatistics{}, err
	}
	labelCount := int32(len(rawDataLabelList))
	rawDataLabelMap := make(map[uuid.UUID][]uuid.UUID)
	for _, rl := range rawDataLabelList {
		rawDataLabelMap[rl.LabelId] = append(rawDataLabelMap[rl.LabelId], rl.LabelId)
	}
	labelDistribution := make([]vb.LabelDistribution, 0)
	for labelId, labelList := range rawDataLabelMap {
		labelDistribution = append(labelDistribution, vb.LabelDistribution{
			LabelId: labelId.String(),
			Count:   int32(len(labelList)),
			Ratio:   float32(len(labelList)) / float32(labelCount),
		})
	}
	return vb.DataViewStatistics{
		ItemCount:         int32(count),
		LabelCount:        labelCount,
		LabelDistribution: labelDistribution,
	}, nil
}
