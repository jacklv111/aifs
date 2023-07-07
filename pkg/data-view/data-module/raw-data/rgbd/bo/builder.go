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
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	rawdatatype "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	rgbddo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/do"
	vb "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/value-object"
)

func BuildWithReader(reader io.ReadSeeker, fileName string, meta vb.RgbdRawDataMeta) basicbo.DataInterface {
	rawDataId := uuid.New()
	return &RgbdRawDataBo{
		DataBaseImpl: basicbo.DataBaseImpl{
			DataItemDo: basicdo.DataItemDo{ID: rawDataId, Type: rawdatatype.RGBD, Name: fileName},
			ReadSeeker: reader,
		},
		rgbdExtDo: rgbddo.RgbdExtDo{
			ID:          rawDataId,
			Sha256:      meta.Sha256,
			ImageSize:   meta.ImageSize,
			ImageWidth:  int32(meta.ImageWidth),
			ImageHeight: int32(meta.ImageHeight),
			DepthSize:   meta.DepthSize,
			DepthWidth:  int32(meta.DepthWidth),
			DepthHeight: int32(meta.DepthHeight),
		},
	}
}
