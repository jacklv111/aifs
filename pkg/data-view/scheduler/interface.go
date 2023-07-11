/*
 * Created on Tue Jul 11 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package scheduler

import (
	"time"

	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	annosvc "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/service"
	artifactconst "github.com/jacklv111/aifs/pkg/data-view/data-module/artifact/constant"
	artifactsvc "github.com/jacklv111/aifs/pkg/data-view/data-module/artifact/service"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	dszconst "github.com/jacklv111/aifs/pkg/data-view/data-module/dataset-zip/constant"
	dszsvc "github.com/jacklv111/aifs/pkg/data-view/data-module/dataset-zip/service"
	modelconst "github.com/jacklv111/aifs/pkg/data-view/data-module/model/constant"
	modelsvc "github.com/jacklv111/aifs/pkg/data-view/data-module/model/service"
	rawdataconst "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	rawdatasvc "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/service"
	dvdo "github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	"github.com/jacklv111/common-sdk/database"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/scheduler"
	"github.com/jacklv111/common-sdk/scheduler/shedlock"
	"gorm.io/gorm"
)

const (
	BATCH_SIZE = 3000
)

func Start() {
	log.Info("data item gc scheduler")

	scheduleConfig := scheduler.ScheduleConfig{
		Name:         "data-item-gc",
		Interval:     time.Hour * 24,
		InitialDelay: time.Hour,
		Runnable:     houseKeeping,
	}
	shedlockConfig := shedlock.ShedlockConfig{
		Enabled:        true,
		Name:           "data-item-gc-lock",
		LockAtLeastFor: time.Hour * 24,
		LockAtMostFor:  time.Hour * 24,
	}
	scheduler.Schedule(shedlockConfig, scheduleConfig)
}

// houseKeeping 清理掉没有被任何 view（包括软删除的） 看到的数据
func houseKeeping() {
	log.Info("housekeeping starts")
	results := make([]basicdo.DataItemDo, BATCH_SIZE)
	inDataView := database.Db.Table(dvdo.TABLE_DATA_VIEW_ITEM).Where("data_view_items.data_item_id = data_items.id")
	database.Db.Select("id", "type").Where("not exists (?)", inDataView).FindInBatches(&results, BATCH_SIZE, func(tx *gorm.DB, batch int) error {
		resMap := make(map[string][]basicdo.DataItemDo, 0)
		for _, data := range results {
			resMap[data.Type] = append(resMap[data.Type], data)
		}
		for key, val := range resMap {
			log.Infof("clear %d %s type data", len(val), key)
			if rawdataconst.HasRawDataType(key) {
				err := rawdatasvc.DeleteRawData(val)
				log.Errorf("error occurred when housekeeping delete raw data [%s]", err)
			}
			if annotationtemplatetype.HasAnnotationTemplateType(key) {
				err := annosvc.DeleteAnnotations(val)
				if err != nil {
					log.Errorf("error occurred when housekeeping delete annotation [%s]", err)
				}
			}
			if modelconst.IsModelFile(key) {
				err := modelsvc.DeleteModelData(val)
				if err != nil {
					log.Errorf("error occurred when housekeeping delete model [%s]", err)
				}
			}
			if dszconst.IsDatasetZipFile(key) {
				err := dszsvc.DeleteDatasetZipData(val)
				if err != nil {
					log.Errorf("error occurred when housekeeping delete dataset zip [%s]", err)
				}
			}
			if artifactconst.IsArtifactFile(key) {
				err := artifactsvc.DeleteArtifactData(val)
				if err != nil {
					log.Errorf("error occurred when housekeeping delete artifact [%s]", err)
				}
			}
		}
		return nil
	})

}
