/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"io"
	"path/filepath"

	"github.com/google/uuid"
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	annotempbo "github.com/jacklv111/aifs/pkg/annotation-template/bo"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

func BuildWithReader(rawDataId uuid.UUID, reader io.ReadSeeker, fileName string, annoTemp annotempbo.AnnotationTemplateBoInterface) basicbo.AnnotationData {
	annoId := uuid.New()
	return &SegmentationMasksBo{
		AnnotationDataImpl: basicbo.AnnotationDataImpl{
			DataBaseImpl: basicbo.DataBaseImpl{
				DataItemDo: basicdo.DataItemDo{ID: annoId, Type: annotationtemplatetype.SEGMENTATION_MASKS, Name: filepath.Base(fileName)},
				ReadSeeker: reader,
			},
			AnnotationDo: do.AnnotationDo{
				ID:                   annoId,
				DataItemId:           rawDataId,
				AnnotationTemplateId: annoTemp.GetId(),
			},
		},
		AnnoTemp: annoTemp,
	}
}
