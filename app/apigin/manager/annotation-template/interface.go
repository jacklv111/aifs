/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	"github.com/google/uuid"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type annotationTemplateMgrInterface interface {
	GetList(offset int, limit int, annoTempIdList []string) ([]openapi.AnnotationTemplateListItem, error)
	Create(openapi.CreateAnnotationTemplateRequest) (openapi.CreateAnnoTemplateSuccessResp, error)
	GetDetailsById(id uuid.UUID) (*openapi.AnnotationTemplateDetails, error)
	DeleteById(id uuid.UUID) error
	Update(req openapi.UpdateAnnotationTemplateRequest) error
	ExistsById(id uuid.UUID) (bool, error)
	GetTypeById(id uuid.UUID) (string, error)
	CopyAnnotationTemplate(id uuid.UUID) (resp openapi.CopyAnnotationTemplate200Response, err error)
}
