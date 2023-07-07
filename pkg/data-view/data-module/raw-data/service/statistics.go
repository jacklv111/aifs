/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"fmt"

	"github.com/google/uuid"
	rawdataconst "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/repo"
)

func GetTotalDataSize(dataItemIdList []uuid.UUID, rawDataType string) (int64, error) {
	switch rawDataType {
	case rawdataconst.IMAGE:
		return repo.ImageRawDataRepo.GetTotalDataSize(dataItemIdList)
	}
	return 0, fmt.Errorf("raw data type %s is not supported to get total data size", rawDataType)
}
