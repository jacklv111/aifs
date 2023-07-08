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
func buildRawDataViewFromCreateDataViewRequest(req openapi.CreateDataViewRequest) RawDataViewBoInterface {
	return &rawDataViewBo{
		dataViewBo: dataViewBo{
			dataViewDo: do.DataViewDo{
				ID:          uuid.New(),
				Name:        req.DataViewName,
				ViewType:    string(req.ViewType),
				RawDataType: string(req.RawDataType),
				Description: req.Description,
			},
		},
	}
}

func buildRawDataViewWithBo(bo dataViewBo) RawDataViewBoInterface {
	return &rawDataViewBo{
		dataViewBo: bo,
	}
}
