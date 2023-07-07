/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	"github.com/jacklv111/aifs/pkg/annotation-template/bo"
	"github.com/jacklv111/aifs/pkg/annotation-template/repo"
	vb "github.com/jacklv111/aifs/pkg/annotation-template/value-object"
)

type annotationTemplateMgrImpl struct {
}

func (mgr *annotationTemplateMgrImpl) Create(req openapi.CreateAnnotationTemplateRequest) (result openapi.CreateAnnoTemplateSuccessResp, err error) {
	hasType := annotationtemplatetype.HasAnnotationTemplateType(req.Type)

	if !hasType {
		return result, fmt.Errorf("annotation template type %s does not exist", req.Type)
	}
	annoTemplateBo := bo.BuildFromCreateAnnotationTemplateRequest(req)

	var id uuid.UUID
	if id, err = annoTemplateBo.Create(); err != nil {
		return result, err
	}
	result.AnnotationTemplateId = id.String()
	return result, nil
}

func (mgr *annotationTemplateMgrImpl) GetList(offset int, limit int, annoTempIdList []string) ([]openapi.AnnotationTemplateListItem, error) {
	options := vb.ListQueryOptions{Offset: offset, Limit: limit, AnnoTemplateIdList: annoTempIdList}
	doList, err := repo.AnnotationTemplateRepo.GetList(options)
	if err != nil {
		return nil, err
	}
	return assembleToAnnotationTemplateList(doList), nil
}

func (mgr *annotationTemplateMgrImpl) GetDetailsById(id uuid.UUID) (*openapi.AnnotationTemplateDetails, error) {
	annoTemplateBo := bo.BuildWithId(id)
	details, err := annoTemplateBo.GetDetails()
	if err != nil {
		return nil, err
	}

	return assembleToAnnotationTemplateDetails(details), nil
}

func (mgr *annotationTemplateMgrImpl) DeleteById(id uuid.UUID) error {
	annoTemplateBo := bo.BuildWithId(id)
	err := annoTemplateBo.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (mgr *annotationTemplateMgrImpl) Update(req openapi.UpdateAnnotationTemplateRequest) error {
	hasType := annotationtemplatetype.HasAnnotationTemplateType(req.Type)

	if !hasType {
		return fmt.Errorf("annotation template type %s does not exist", req.Type)
	}

	annoTemplateBo := bo.BuildFromUpdateAnnotationTemplateRequest(req)
	if err := annoTemplateBo.Update(); err != nil {
		return err
	}
	return nil
}

func (mgr *annotationTemplateMgrImpl) ExistsById(id uuid.UUID) (bool, error) {
	has, err := repo.AnnotationTemplateRepo.ExistsById(id)
	if err != nil {
		return false, err
	}
	return has, nil
}

func (mgr *annotationTemplateMgrImpl) GetTypeById(id uuid.UUID) (string, error) {
	annoTempType, err := repo.AnnotationTemplateRepo.GetTypeById(id)
	if err != nil {
		return "", err
	}
	return annoTempType, nil
}

func (mgr *annotationTemplateMgrImpl) CopyAnnotationTemplate(id uuid.UUID) (resp openapi.CopyAnnotationTemplate200Response, err error) {
	annoTemplateBo := bo.BuildWithId(id)
	newBo, err := annoTemplateBo.Copy()
	if err != nil {
		return resp, err
	}
	id, err = newBo.Create()
	if err != nil {
		return resp, err
	}
	resp.AnnotationTemplateId = id.String()
	return resp, nil
}

var AnnotationTemplateMgr annotationTemplateMgrInterface

func init() {
	AnnotationTemplateMgr = &annotationTemplateMgrImpl{}
}
