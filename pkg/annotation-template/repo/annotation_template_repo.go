/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"database/sql"

	"github.com/google/uuid"
	annotationtemplate "github.com/jacklv111/aifs/pkg/annotation-template"
	"github.com/jacklv111/aifs/pkg/annotation-template/do"
	valueobject "github.com/jacklv111/aifs/pkg/annotation-template/value-object"
	"github.com/jacklv111/common-sdk/database"
	utilerror "github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"gorm.io/gorm"
)

type annotationTemplateRepoImpl struct {
}

func (repo *annotationTemplateRepoImpl) Create(annoTemplateDo do.AnnotationTemplateDo, annoTemplateExtDo do.AnnotationTemplateExtDo, labelList []do.LabelDo) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&annoTemplateDo).Error; err != nil {
			return err
		}

		if !annoTemplateExtDo.IsEmpty() {
			if err := tx.Create(&annoTemplateExtDo).Error; err != nil {
				return err
			}
		}

		// use default create batch
		if err := tx.Create(labelList).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *annotationTemplateRepoImpl) GetList(options valueobject.ListQueryOptions) ([]do.ListItem, error) {
	var result []do.ListItem

	fullQuery := database.Db.Model(do.AnnotationTemplateDo{}).
		Joins("left join labels on annotation_templates.id=labels.annotation_template_id").
		Group("annotation_templates.id").
		Select("annotation_templates.*, count(labels.id) as label_count").
		Offset(options.Offset).
		Limit(options.Limit)

	if options.AnnoTemplateIdList != nil {
		fullQuery.Where("annotation_templates.id in ?", options.AnnoTemplateIdList)
	}

	err := fullQuery.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *annotationTemplateRepoImpl) GetById(annoTempId uuid.UUID) (annoTempDo do.AnnotationTemplateDo, annoTempExtDo do.AnnotationTemplateExtDo, labelDoList []do.LabelDo, err error) {
	res := database.Db.First(&annoTempDo, annoTempId)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			err = annotationtemplate.ErrAnnotationTemplateNotFound
			return
		}
		err = res.Error
		return
	}

	err = database.Db.Where("annotation_template_id = ?", annoTempId).First(&annoTempExtDo).Error
	if err != nil {
		if err.Error() == utilerror.RECORD_NOT_FOUND {
			log.Infof("annotation template id %s has no ext data", annoTempId.String())
		} else {
			return
		}
	}

	err = database.Db.Where(&do.LabelDo{AnnotationTemplateId: annoTempId}).Find(&labelDoList).Error
	if err != nil {
		return
	}
	return
}

func (repo *annotationTemplateRepoImpl) GetTypeById(annoTempId uuid.UUID) (string, error) {
	var annoTempDo do.AnnotationTemplateDo
	res := database.Db.First(&annoTempDo, annoTempId)
	if res.Error != nil {
		if res.Error.Error() == utilerror.RECORD_NOT_FOUND {
			err := annotationtemplate.ErrAnnotationTemplateNotFound
			return "", err
		}
		err := res.Error
		return "", err
	}
	return annoTempDo.Type, nil
}

func (repo *annotationTemplateRepoImpl) ExistsById(annoTempId uuid.UUID) (bool, error) {
	var count int64
	err := database.Db.Model(&do.AnnotationTemplateDo{}).Where("id = ?", annoTempId.String()).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (repo *annotationTemplateRepoImpl) Delete(annoTempId uuid.UUID) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&do.AnnotationTemplateDo{ID: annoTempId}).Error; err != nil {
			return err
		}

		if err := tx.Where("annotation_template_id = ?", annoTempId).Delete(&do.AnnotationTemplateExtDo{}).Error; err != nil {
			return err
		}

		// batch delete labels
		if err := tx.Where("annotation_template_id = ?", annoTempId.String()).Delete(&do.LabelDo{}).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *annotationTemplateRepoImpl) Update(destAnnoTemplateDo do.AnnotationTemplateDo, destAnnoTempExtDo do.AnnotationTemplateExtDo, destLabelList []do.LabelDo) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		srcAnnoTemplateDo, srcAnnoTempExtDo, srcLabelList, err := repo.GetById(destAnnoTemplateDo.ID)
		if err != nil {
			// TODO: 这里需要区分 error，如果是 not found，那就是 bad request，如果是 db 的问题，那就是 5XX
			return err
		}
		if err := tx.Model(&srcAnnoTemplateDo).Updates(destAnnoTemplateDo).Error; err != nil {
			return err
		}

		// 如果某个 annotation template 在创建的时候就没有某个 ext 数据，那之后都不会有
		if !srcAnnoTempExtDo.IsEmpty() {
			if err := tx.Model(&srcAnnoTempExtDo).Updates(destAnnoTempExtDo).Error; err != nil {
				return err
			}
		}

		if err = updateLabels(srcLabelList, destLabelList, tx); err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: false})
}

func getLabelIdMap(labelList []do.LabelDo) (res map[string]do.LabelDo) {
	res = make(map[string]do.LabelDo)
	for _, data := range labelList {
		res[data.ID.String()] = data
	}
	return
}

func updateLabels(srcLabelList []do.LabelDo, destLabelList []do.LabelDo, tx *gorm.DB) error {
	srcLabelIdMap := getLabelIdMap(srcLabelList)
	destLabelIdMap := getLabelIdMap(destLabelList)
	// dest 中没有，src 中有，删除
	for _, data := range srcLabelList {
		if _, ok := destLabelIdMap[data.ID.String()]; !ok {
			// 硬删除
			tx.Unscoped().Delete(&data)
		}
	}

	for _, destLabel := range destLabelList {
		id := destLabel.ID.String()
		// dest 中有，src 中没有
		if srcLabel, ok := srcLabelIdMap[id]; !ok {
			if err := tx.Create(&destLabel).Error; err != nil {
				return err
			}
		} else { // 都有，更新
			// 这里会全列更新
			if err := tx.Model(&srcLabel).Updates(destLabelIdMap[id]).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

var AnnotationTemplateRepo annotationTemplateRepoInterface

func init() {
	AnnotationTemplateRepo = &annotationTemplateRepoImpl{}
}
