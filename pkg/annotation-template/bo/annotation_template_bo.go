/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"fmt"

	"github.com/google/uuid"
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	"github.com/jacklv111/aifs/pkg/annotation-template/do"
	"github.com/jacklv111/aifs/pkg/annotation-template/repo"
	vb "github.com/jacklv111/aifs/pkg/annotation-template/value-object"
	"github.com/jacklv111/common-sdk/collection/mapset"
)

type annotationTemplateBo struct {
	annotationTemplateDo    do.AnnotationTemplateDo
	annotationTemplateExtDo do.AnnotationTemplateExtDo
	labelDoList             []do.LabelDo
	// key is label id
	labelIdMap    map[uuid.UUID]do.LabelDo
	labelColorMap map[int32]do.LabelDo
	wordSet       mapset.Set[string]
	isSynced      bool
}

func (bo *annotationTemplateBo) Create() (uuid.UUID, error) {
	// validate
	if bo.GetType() == annotationtemplatetype.OCR {
		if bo.annotationTemplateExtDo.WordList.IsEmpty() {
			return uuid.Nil, fmt.Errorf("%s type should have word list", annotationtemplatetype.OCR)
		}
	}

	err := repo.AnnotationTemplateRepo.Create(bo.annotationTemplateDo, bo.annotationTemplateExtDo, bo.labelDoList)
	if err != nil {
		return uuid.Nil, err
	}
	return bo.annotationTemplateDo.ID, nil
}

func (bo *annotationTemplateBo) GetDetails() (*vb.AnnotationTemplateDetails, error) {
	if bo.annotationTemplateDo.ID == uuid.Nil {
		return nil, fmt.Errorf("annotation template id is null, get details failed")
	}
	var err error
	bo.annotationTemplateDo, bo.annotationTemplateExtDo, bo.labelDoList, err = repo.AnnotationTemplateRepo.GetById(bo.annotationTemplateDo.ID)
	if err != nil {
		return nil, err
	}
	return &vb.AnnotationTemplateDetails{
		AnnotationTemplateDo:    bo.annotationTemplateDo,
		AnnotationTemplateExtDo: bo.annotationTemplateExtDo,
		LabelDoList:             bo.labelDoList,
	}, nil
}

func (bo *annotationTemplateBo) Delete() error {
	if bo.annotationTemplateDo.ID == uuid.Nil {
		return fmt.Errorf("annotation template id is null, delete failed")
	}
	err := repo.AnnotationTemplateRepo.Delete(bo.annotationTemplateDo.ID)
	if err != nil {
		return err
	}

	bo.annotationTemplateDo.ID = uuid.Nil
	return nil
}

func (bo *annotationTemplateBo) Update() error {
	if bo.annotationTemplateDo.ID == uuid.Nil {
		return fmt.Errorf("annotation template id is null, delete failed")
	}

	err := repo.AnnotationTemplateRepo.Update(bo.annotationTemplateDo, bo.annotationTemplateExtDo, bo.labelDoList)
	if err != nil {
		return err
	}

	return nil
}

func (bo *annotationTemplateBo) Sync() (err error) {
	if bo.isSynced {
		return
	}
	bo.annotationTemplateDo, bo.annotationTemplateExtDo, bo.labelDoList, err = repo.AnnotationTemplateRepo.GetById(bo.annotationTemplateDo.ID)
	if err != nil {
		return err
	}
	bo.initialize()
	bo.isSynced = true
	return
}

func (bo *annotationTemplateBo) Copy() (AnnotationTemplateBoInterface, error) {
	err := bo.Sync()
	if err != nil {
		return nil, err
	}
	newBo := annotationTemplateBo{
		annotationTemplateDo:    bo.annotationTemplateDo,
		annotationTemplateExtDo: bo.annotationTemplateExtDo,
		labelDoList:             bo.labelDoList,
	}
	newAnnoTempId := uuid.New()
	newBo.annotationTemplateDo.ID = newAnnoTempId
	newBo.annotationTemplateExtDo.AnnotationTemplateId = newAnnoTempId
	for idx := range newBo.labelDoList {
		newBo.labelDoList[idx].ID = uuid.New()
		newBo.labelDoList[idx].AnnotationTemplateId = newAnnoTempId
	}

	bo.initialize()
	return &newBo, nil
}

func (bo *annotationTemplateBo) HasLabel(idList []uuid.UUID) bool {
	if len(idList) == 0 {
		return true
	}
	if bo.labelIdMap == nil {
		return false
	}
	for _, id := range idList {
		if _, ok := bo.labelIdMap[id]; !ok {
			return false
		}
	}
	return true
}

func (bo *annotationTemplateBo) GetType() string {
	return bo.annotationTemplateDo.Type
}

func (bo *annotationTemplateBo) GetId() uuid.UUID {
	return bo.annotationTemplateDo.ID
}

func (bo *annotationTemplateBo) GetLabelIdByColor(color int32) uuid.UUID {
	labelDo, ok := bo.labelColorMap[color]
	if !ok {
		return uuid.Nil
	}
	return labelDo.ID
}

func (bo *annotationTemplateBo) initialize() {
	bo.labelIdMap = make(map[uuid.UUID]do.LabelDo)
	bo.labelColorMap = make(map[int32]do.LabelDo)
	for _, data := range bo.labelDoList {
		bo.labelIdMap[data.ID] = data
		bo.labelColorMap[data.Color] = data
	}

	if !bo.annotationTemplateExtDo.WordList.IsEmpty() {
		bo.wordSet = mapset.NewSet[string]()
		for _, word := range bo.annotationTemplateExtDo.WordList {
			bo.wordSet.Add(word)
		}
	}
}
