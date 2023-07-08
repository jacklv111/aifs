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
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	"github.com/jacklv111/aifs/pkg/data-view/data-view/do"
)

// 构造 bo
func buildAnnotationViewFromCreateDataViewRequest(req openapi.CreateDataViewRequest) AnnotationViewBoInterface {
	relatedDataViewId, _ := uuid.Parse(req.RelatedDataViewId)
	annotationTemplateId, _ := uuid.Parse(req.AnnotationTemplateId)
	return &annotationViewBo{
		dataViewBo: dataViewBo{
			dataViewDo: do.DataViewDo{
				ID:                   uuid.New(),
				Name:                 req.DataViewName,
				ViewType:             string(req.ViewType),
				RelatedDataViewId:    relatedDataViewId,
				AnnotationTemplateId: annotationTemplateId,
				Description:          req.Description,
			},
		},
	}
}

func buildAnnotationViewWithBo(bo dataViewBo) AnnotationViewBoInterface {
	return &annotationViewBo{
		dataViewBo: bo,
	}
}
