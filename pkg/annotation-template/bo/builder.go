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
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	"github.com/jacklv111/aifs/pkg/annotation-template/do"
)

// 构造 bo
func BuildFromCreateAnnotationTemplateRequest(req openapi.CreateAnnotationTemplateRequest) AnnotationTemplateBoInterface {
	switch req.Type {
	case annotationtemplatetype.COCO_TYPE:
		return BuildFromCreateAnnotationTemplateRequestCocoType(req)
	case annotationtemplatetype.OCR:
		return BuildFromCreateAnnotationTemplateRequestOcr(req)
	default:
		return BuildFromCreateAnnotationTemplateRequestDefault(req)
	}
}

func BuildWithId(id uuid.UUID) AnnotationTemplateBoInterface {
	return &annotationTemplateBo{
		annotationTemplateDo: do.AnnotationTemplateDo{ID: id},
	}
}

func BuildFromUpdateAnnotationTemplateRequest(req openapi.UpdateAnnotationTemplateRequest) AnnotationTemplateBoInterface {
	switch req.Type {
	case annotationtemplatetype.COCO_TYPE:
		return BuildFromUpdateAnnotationTemplateRequestCocoType(req)
	case annotationtemplatetype.OCR:
		return BuildFromUpdateAnnotationTemplateRequestOcr(req)
	default:
		return BuildFromUpdateAnnotationTemplateRequestDefault(req)
	}
}
