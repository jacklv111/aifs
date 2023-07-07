/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	"github.com/jacklv111/aifs/pkg/annotation-template/do"
	vb "github.com/jacklv111/aifs/pkg/annotation-template/value-object"
)

func assembleToAnnotationTemplateList(doList []do.ListItem) []openapi.AnnotationTemplateListItem {
	var result []openapi.AnnotationTemplateListItem
	for _, data := range doList {
		result = append(result, openapi.AnnotationTemplateListItem{
			Id:         data.Id,
			Name:       data.Name,
			CreateAt:   data.CreateAt,
			Type:       data.Type,
			LabelCount: data.LabelCount,
		})
	}
	return result
}

func assembleToAnnotationTemplateDetails(details *vb.AnnotationTemplateDetails) *openapi.AnnotationTemplateDetails {
	var labels []openapi.Label
	for _, data := range details.LabelDoList {
		labels = append(labels, openapi.Label{
			Id:                data.ID.String(),
			Name:              data.Name,
			Color:             data.Color,
			SuperCategoryName: data.SuperCategoryName,
			KeyPointDef:       data.KeyPointDef,
			KeyPointSkeleton:  data.KeyPointSkeleton,
		})
	}
	return &openapi.AnnotationTemplateDetails{
		Id:       details.AnnotationTemplateDo.ID.String(),
		Name:     details.AnnotationTemplateDo.Name,
		CreateAt: details.AnnotationTemplateDo.CreateAt,
		UpdateAt: details.AnnotationTemplateDo.UpdateAt,
		Type:     details.AnnotationTemplateDo.Type,
		WordList: details.AnnotationTemplateExtDo.WordList,
		Labels:   labels,
	}
}
