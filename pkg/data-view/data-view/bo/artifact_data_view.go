/*
 * Created on Sat Jul 08 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"github.com/google/uuid"
	datamodule "github.com/jacklv111/aifs/pkg/data-view/data-module"
	artifactbo "github.com/jacklv111/aifs/pkg/data-view/data-module/artifact/bo"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	dvrepo "github.com/jacklv111/aifs/pkg/data-view/data-view/repo"
	dvvb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
	storemgr "github.com/jacklv111/aifs/pkg/store/manager"
	storevb "github.com/jacklv111/aifs/pkg/store/value-object"
)

type artifactViewBo struct {
	dataViewBo
}

func (bo *artifactViewBo) UploadArtifactFile(input basicvb.UploadArtifactFileParams) (err error) {
	if err = bo.loadDataViewDo(); err != nil {
		return
	}

	artifactBo := artifactbo.BuildFromBuffer(input.File, input.FileName)

	// save meta
	err = artifactBo.Create()
	if err != nil {
		return
	}

	// save data
	var storeParams storevb.UploadParams
	storeParams.DataType = datamodule.ARTIFACT
	storeParams.AddItem(artifactBo.GetId(), artifactBo.GetReader(), artifactBo.GetName())

	err = storemgr.StoreMgr.Upload(storeParams)
	if err != nil {
		return err
	}

	// save data view items
	return dvrepo.DataViewRepo.CreateDataViewItemsIgnoreConflict(bo.dataViewDo.ID, []uuid.UUID{artifactBo.ID})
}

func (bo *artifactViewBo) GetArtifactLocations() (result dvvb.ArtifactLocationResult, err error) {
	result.DataViewId = bo.dataViewDo.ID.String()
	result.ViewType = bo.dataViewDo.ViewType
	result.DataItemDoList, result.LocationMap, err = bo.getDataItemAndLocations()
	return
}
