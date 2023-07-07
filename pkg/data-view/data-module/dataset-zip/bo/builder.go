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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/dataset-zip/constant"
)

func BuildFromBuffer(reader io.Reader, name string) *DatasetZipBo {
	dataItemId := uuid.New()
	return &DatasetZipBo{
		DataBaseImpl: bo.DataBaseImpl{
			DataItemDo: basicdo.DataItemDo{
				ID:   dataItemId,
				Type: constant.DATASET_ZIP_FILE,
				Name: name,
			},
		},
		reader: reader,
	}
}
