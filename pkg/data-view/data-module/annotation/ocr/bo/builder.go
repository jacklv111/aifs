/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"database/sql"

	"github.com/google/uuid"
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

func BuildWithAnnoData(rawDataId uuid.UUID, annoTempId uuid.UUID, annoData string) basicbo.AnnotationData {
	annotationId := uuid.New()
	return &OcrBo{
		AnnotationDataImpl: basicbo.AnnotationDataImpl{
			DataBaseImpl: basicbo.DataBaseImpl{
				DataItemDo: basicdo.DataItemDo{ID: annotationId, Type: annotationtemplatetype.OCR},
			},
			AnnotationDo: do.AnnotationDo{
				ID:                   annotationId,
				DataItemId:           rawDataId,
				AnnotationTemplateId: annoTempId,
				TextData:             sql.NullString{String: annoData, Valid: true},
			},
		},
	}
}
