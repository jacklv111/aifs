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

	"github.com/google/uuid"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	"github.com/jacklv111/aifs/pkg/data-view/data-view/do"
)

// 构造 bo
func buildDatasetZipViewFromCreateDataViewRequest(req openapi.CreateDataViewRequest) DatasetZipViewBoInterface {
	var annotationTemplateId uuid.UUID
	if req.AnnotationTemplateId != "" {
		annotationTemplateId = uuid.MustParse(req.AnnotationTemplateId)
	}
	return &datasetZipViewBo{
		dataViewBo: dataViewBo{
			dataViewDo: do.DataViewDo{
				ID:                   uuid.New(),
				Name:                 req.DataViewName,
				ViewType:             string(req.ViewType),
				ZipFormat:            sql.NullString{String: string(req.ZipFormat), Valid: true},
				RawDataViewId:        sql.NullString{String: req.RawDataViewId, Valid: req.RawDataViewId != ""},       // 有可能为空
				AnnotationViewId:     sql.NullString{String: req.AnnotationViewId, Valid: req.AnnotationViewId != ""}, // 有可能为空
				AnnotationTemplateId: annotationTemplateId,
				Description:          req.Description,
			},
		},
	}
}

func buildDatasetZipWithBo(bo dataViewBo) DatasetZipViewBoInterface {
	return &datasetZipViewBo{
		dataViewBo: bo,
	}
}
