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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/artifact/constant"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

func BuildFromBuffer(reader io.Reader, name string) *ArtifactFileBo {
	dataItemId := uuid.New()
	return &ArtifactFileBo{
		DataBaseImpl: bo.DataBaseImpl{
			DataItemDo: basicdo.DataItemDo{
				ID:   dataItemId,
				Type: constant.ARTIFACT_FILE,
				Name: name,
			},
		},
		reader: reader,
	}
}
