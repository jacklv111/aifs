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
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	rawdataconst "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	imagebo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/bo"
	p3dbo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/bo"
	rgbdbo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/bo"
)

// UploadRawData 上传 raw data，数据在其他节点
//
//	@param input
//	@param rawDataType
//	@return []uuid.UUID
//	@return error
func UploadRawData(input basicvb.UploadRawDataParam, rawDataType string) ([]uuid.UUID, error) {
	return LoadFromRemote(input, rawDataType)
}

func saveMetadata(rawDataType string, rawDataList []basicbo.DataInterface) ([]uuid.UUID, error) {
	switch rawDataType {
	case rawdataconst.IMAGE:
		err := imagebo.FixDataItemId(rawDataList)
		if err != nil {
			return nil, err
		}
		return imagebo.CreateBatch(rawDataList)
	case rawdataconst.RGBD:
		err := rgbdbo.FixDataItemId(rawDataList)
		if err != nil {
			return nil, err
		}
		return rgbdbo.CreateBatch(rawDataList)
	case rawdataconst.POINTS_3D:
		err := p3dbo.FixDataItemId(rawDataList)
		if err != nil {
			return nil, err
		}
		return p3dbo.CreateBatch(rawDataList)
	}
	return nil, fmt.Errorf("saveMetadata cant handle raw data type %s", rawDataType)
}
