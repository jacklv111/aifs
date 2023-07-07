/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"github.com/google/uuid"
	"github.com/jacklv111/aifs/pkg/annotation-template/do"
	valueobject "github.com/jacklv111/aifs/pkg/annotation-template/value-object"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type annotationTemplateRepoInterface interface {
	// 查询 annotation template list.
	GetList(valueobject.ListQueryOptions) ([]do.ListItem, error)
	// 创建 annotation template
	Create(annoTemplateDo do.AnnotationTemplateDo, annoTemplateExtDo do.AnnotationTemplateExtDo, labelList []do.LabelDo) error
	// 通过 id 获取 annotation template 数据
	GetById(annoTempId uuid.UUID) (do.AnnotationTemplateDo, do.AnnotationTemplateExtDo, []do.LabelDo, error)
	// 通过 id 删除 annotation template 数据
	Delete(annoTempId uuid.UUID) error
	// 更新 annotation template，传入期望的目标状态
	Update(annoTemplateDo do.AnnotationTemplateDo, annoTemplateExtDo do.AnnotationTemplateExtDo, labelList []do.LabelDo) error
	// 判断 annotation template id 是否存在
	ExistsById(annoTempId uuid.UUID) (bool, error)
	// 获取 annotation template type
	GetTypeById(annoTempId uuid.UUID) (string, error)
}
