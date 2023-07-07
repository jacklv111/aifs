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
	"path/filepath"

	"github.com/google/uuid"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/do"
)

func BuildWithLocalPath(filePath string) basicbo.DataInterface {
	dataItemId := uuid.New()
	return &ImageRawDataBo{
		DataBaseImpl: basicbo.DataBaseImpl{
			DataItemDo: basicdo.DataItemDo{ID: dataItemId, Type: constant.IMAGE, Name: filepath.Base(filePath)},
			LocalPath:  basicbo.LocalPath(filePath),
		},
		imageExtDo:   do.ImageExtDo{ID: dataItemId},
		imageScoreDo: do.ImageScoreDo{ID: dataItemId},
	}
}

func BuildWithReader(reader io.ReadSeeker, fileName string) basicbo.DataInterface {
	dataItemId := uuid.New()
	return &ImageRawDataBo{
		DataBaseImpl: basicbo.DataBaseImpl{
			DataItemDo: basicdo.DataItemDo{ID: dataItemId, Type: constant.IMAGE, Name: fileName},
			ReadSeeker: reader,
		},
		imageExtDo:   do.ImageExtDo{ID: dataItemId},
		imageScoreDo: do.ImageScoreDo{ID: dataItemId},
	}
}
