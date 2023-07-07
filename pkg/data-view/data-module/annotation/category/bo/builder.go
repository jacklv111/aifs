/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"github.com/google/uuid"
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

func BuildWithAnnoData(rawDataId uuid.UUID, annoTempId uuid.UUID, labelId uuid.UUID) basicbo.AnnotationData {
	dataItemId := uuid.New()
	return &CategoryBo{
		AnnotationDataImpl: basicbo.AnnotationDataImpl{
			DataBaseImpl: basicbo.DataBaseImpl{
				DataItemDo: basicdo.DataItemDo{ID: dataItemId, Type: annotationtemplatetype.CATEGORY},
			},
			AnnotationDo:     do.AnnotationDo{ID: dataItemId, AnnotationTemplateId: annoTempId, DataItemId: rawDataId},
			RawDataLabelList: []do.RawDataLabelDo{{AnnotationId: dataItemId, LabelId: labelId, RawDataId: rawDataId}},
		},
	}
}
