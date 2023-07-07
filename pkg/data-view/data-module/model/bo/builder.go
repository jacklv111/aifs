/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"io"

	"github.com/google/uuid"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/model/constant"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/model/do"
)

func BuildFromBuffer(reader io.ReadSeeker, name string) *ModelDataBo {
	dataItemId := uuid.New()
	return &ModelDataBo{
		DataBaseImpl: bo.DataBaseImpl{
			DataItemDo: basicdo.DataItemDo{
				ID:   dataItemId,
				Type: constant.MODEL_FILE,
				Name: name,
			},
			ReadSeeker: reader,
		},
		ModelExtDo: do.ModelExtDo{ID: dataItemId},
	}
}
