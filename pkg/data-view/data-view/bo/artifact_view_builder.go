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
	"github.com/jacklv111/aifs/app/apigin/view-object/openapi"
	"github.com/jacklv111/aifs/pkg/data-view/data-view/do"
)

// 构造 bo
func buildArtifactViewFromCreateDataViewRequest(req openapi.CreateDataViewRequest) ArtifactViewBoInterface {
	return &artifactViewBo{
		dataViewBo: dataViewBo{
			dataViewDo: do.DataViewDo{
				ID:          uuid.New(),
				Name:        req.DataViewName,
				ViewType:    string(req.ViewType),
				Description: req.Description,
			},
		},
	}
}

func buildArtifactWithBo(bo dataViewBo) ArtifactViewBoInterface {
	return &artifactViewBo{
		dataViewBo: bo,
	}
}
