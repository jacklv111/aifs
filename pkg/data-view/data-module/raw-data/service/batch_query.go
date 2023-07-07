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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	imagerepo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/repo"
	p3drepo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/repo"
	rgbdrepo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/repo"
)

func GetHashList(dataItemIdList []uuid.UUID, rawDataType string) ([]do.IdHash, error) {
	switch rawDataType {
	case constant.IMAGE:
		return imagerepo.ImageRawDataRepo.GetHashList(dataItemIdList)
	case constant.RGBD:
		return rgbdrepo.RgbdRawDataRepo.GetHashList(dataItemIdList)
	case constant.POINTS_3D:
		return p3drepo.P3dRawDataRepo.GetHashList(dataItemIdList)
	}
	return nil, fmt.Errorf("cant handle rawd data type %s", rawDataType)
}
