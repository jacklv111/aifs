/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"fmt"

	"github.com/google/uuid"
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	annotempbo "github.com/jacklv111/aifs/pkg/annotation-template/bo"
	catebo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/category/bo"
	cocobo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/coco-type/bo"
	ocrbo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/ocr/bo"
	p3dbo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/points-3d/bo"
	rgbdbo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/rgbd/bo"
	segmasksbo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/segmentation-masks/bo"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
)

func UploadAnnotations(input basicvb.UploadAnnotationParam, annotationTemplateId uuid.UUID) ([]uuid.UUID, error) {
	annoTempBo := annotempbo.BuildWithId(annotationTemplateId)
	if err := annoTempBo.Sync(); err != nil {
		return nil, err
	}
	return LoadFromRemote(input, annoTempBo)
}

func saveMetadata(annoType string, dataList []basicbo.AnnotationData) ([]uuid.UUID, error) {
	switch annoType {
	case annotationtemplatetype.CATEGORY:
		return catebo.CreateBatch(dataList)
	case annotationtemplatetype.COCO_TYPE:
		return cocobo.CreateBatch(dataList)
	case annotationtemplatetype.OCR:
		return ocrbo.CreateBatch(dataList)
	case annotationtemplatetype.RGBD:
		return rgbdbo.CreateBatch(dataList)
	case annotationtemplatetype.SEGMENTATION_MASKS:
		return segmasksbo.CreateBatch(dataList)
	case annotationtemplatetype.POINTS_3D:
		return p3dbo.CreateBatch(dataList)
	default:
		return nil, fmt.Errorf("saveMetadata cant handle annotation type %s", annoType)
	}
}

func validate(annoTempBo annotempbo.AnnotationTemplateBoInterface, dataList []basicbo.AnnotationData, rawDataIdChecker func(dataItemIdList []uuid.UUID) error) error {
	for _, anno := range dataList {
		if !annoTempBo.HasLabel(anno.GetLabels()) {
			return fmt.Errorf("annotation has labels %v, some labels cant be recognized", anno.GetLabels())
		}
		if annoTempBo.GetId() != anno.GetAnnotationTemplateId() {
			return fmt.Errorf("it is not allowed to upload annotation with annotation template %s to data view with annotation template id %s", anno.GetAnnotationTemplateId(), annoTempBo.GetId())
		}
	}
	if rawDataIdChecker != nil {
		rawDataIdList := basicbo.GetRawDataIdList(dataList)
		err := rawDataIdChecker(rawDataIdList)
		if err != nil {
			return err
		}
	}
	return nil
}
