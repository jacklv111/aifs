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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/dataset-zip/repo"
)

type DatasetZipBo struct {
	bo.DataBaseImpl
	// server 端不进行读取操作
	reader io.Reader
}

func (bo *DatasetZipBo) GetName() string {
	return bo.Name
}

func (bo *DatasetZipBo) GetId() uuid.UUID {
	return bo.ID
}

func (bo *DatasetZipBo) Create() error {
	return repo.DatasetZipRepo.Create(bo.DataItemDo)
}

func (bo *DatasetZipBo) GetReader() io.Reader {
	return bo.reader
}
