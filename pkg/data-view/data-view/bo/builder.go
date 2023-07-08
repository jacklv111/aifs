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
	"time"

	"github.com/google/uuid"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	datamodule "github.com/jacklv111/aifs/pkg/data-view/data-module"
	"github.com/jacklv111/aifs/pkg/data-view/data-view/do"
)

// 构造 bo
func BuildFromCreateDataViewRequest(req openapi.CreateDataViewRequest) (DataViewBoInterface, error) {
	switch req.ViewType {
	case datamodule.RAW_DATA:
		return buildRawDataViewFromCreateDataViewRequest(req), nil
	case datamodule.ANNOTATION:
		return buildAnnotationViewFromCreateDataViewRequest(req), nil
	case datamodule.MODEL:
		return buildModelViewFromCreateDataViewRequest(req), nil
	case datamodule.DATASET_ZIP:
		return buildDatasetZipViewFromCreateDataViewRequest(req), nil
	case datamodule.ARTIFACT:
		return buildArtifactViewFromCreateDataViewRequest(req), nil
	}
	return nil, fmt.Errorf("error view type %s, view type should be in %v", req.ViewType, datamodule.GetList())
}

func BuildWithId(id uuid.UUID) (DataViewBoInterface, error) {
	dataViewBo := dataViewBo{dataViewDo: do.DataViewDo{ID: id}}
	err := dataViewBo.loadDataViewDo()
	if err != nil {
		return nil, err
	}
	switch dataViewBo.dataViewDo.ViewType {
	case datamodule.RAW_DATA:
		return buildRawDataViewWithBo(dataViewBo), nil
	case datamodule.ANNOTATION:
		return buildAnnotationViewWithBo(dataViewBo), nil
	case datamodule.MODEL:
		return buildModelViewWithBo(dataViewBo), nil
	case datamodule.DATASET_ZIP:
		return buildDatasetZipWithBo(dataViewBo), nil
	case datamodule.ARTIFACT:
		return buildArtifactWithBo(dataViewBo), nil
	}
	return nil, fmt.Errorf("error type %s", dataViewBo.dataViewDo.ViewType)
}

func BuildFromBoType(id uuid.UUID, name, desc string) (DataViewBoInterface, error) {
	fromBo := dataViewBo{dataViewDo: do.DataViewDo{ID: id}}
	err := fromBo.loadDataViewDo()
	if err != nil {
		return nil, err
	}
	destBo := dataViewBo{dataViewDo: fromBo.dataViewDo}
	destBo.dataViewDo.ID = uuid.New()
	destBo.dataViewDo.Name = name
	destBo.dataViewDo.Description = desc
	destBo.dataViewDo.CreateAt = time.Now().Unix()
	destBo.dataViewDo.UpdateAt = time.Now().Unix()
	return &destBo, nil
}
