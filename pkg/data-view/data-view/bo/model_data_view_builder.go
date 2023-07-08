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
func buildModelViewFromCreateDataViewRequest(req openapi.CreateDataViewRequest) ModelViewBoInterface {
	return &modelViewBo{
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

func buildModelViewWithBo(bo dataViewBo) ModelViewBoInterface {
	return &modelViewBo{
		dataViewBo: bo,
	}
}
