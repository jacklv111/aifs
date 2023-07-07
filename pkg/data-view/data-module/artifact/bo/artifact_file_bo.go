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
	"github.com/jacklv111/aifs/pkg/data-view/data-module/artifact/repo"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
)

type ArtifactFileBo struct {
	bo.DataBaseImpl
	reader io.Reader
}

func (bo *ArtifactFileBo) GetName() string {
	return bo.Name
}

func (bo *ArtifactFileBo) GetId() uuid.UUID {
	return bo.ID
}

func (bo *ArtifactFileBo) Create() error {
	return repo.ArtifactRepo.Create(bo.DataItemDo)
}

func (bo *ArtifactFileBo) GetReader() io.Reader {
	return bo.reader
}
