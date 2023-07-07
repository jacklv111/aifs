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
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	"github.com/jacklv111/aifs/pkg/annotation-template/do"
)

func BuildFromCreateAnnotationTemplateRequestCocoType(req openapi.CreateAnnotationTemplateRequest) AnnotationTemplateBoInterface {
	annoTempId := uuid.New()
	var labelDoList []do.LabelDo
	for _, label := range req.Labels {
		labelDoList = append(labelDoList, do.LabelDo{
			ID:                   uuid.New(),
			AnnotationTemplateId: annoTempId,
			Name:                 label.Name,
			Color:                label.Color,
			SuperCategoryName:    label.SuperCategoryName,
			KeyPointDef:          do.KeyPointDefType(label.KeyPointDef),
			KeyPointSkeleton:     do.KeyPointSkeletonType(label.KeyPointSkeleton),
			CoverImageUrl:        label.CoverImageUrl,
		})
	}

	return &annotationTemplateBo{
		annotationTemplateDo: do.AnnotationTemplateDo{
			ID:          annoTempId,
			Type:        req.Type,
			Name:        req.Name,
			Description: req.Description,
		},
		labelDoList: labelDoList,
	}
}

func BuildFromUpdateAnnotationTemplateRequestCocoType(req openapi.UpdateAnnotationTemplateRequest) AnnotationTemplateBoInterface {
	var labelDoList []do.LabelDo
	annoTempId := uuid.MustParse(req.Id)
	for _, label := range req.Labels {
		labelDo := do.LabelDo{
			Name:             label.Name,
			Color:            label.Color,
			KeyPointDef:      do.KeyPointDefType(label.KeyPointDef),
			KeyPointSkeleton: do.KeyPointSkeletonType(label.KeyPointSkeleton),
			CoverImageUrl:    label.CoverImageUrl,
		}
		// 有 id
		if label.Id != "" {
			labelDo.ID = uuid.MustParse(label.Id)
		} else {
			labelDo.ID = uuid.New()
		}
		labelDo.AnnotationTemplateId = annoTempId
		labelDoList = append(labelDoList, labelDo)
	}
	return &annotationTemplateBo{
		annotationTemplateDo: do.AnnotationTemplateDo{
			ID:          annoTempId,
			Type:        req.Type,
			Name:        req.Name,
			Description: req.Description,
		},
		labelDoList: labelDoList,
	}
}
