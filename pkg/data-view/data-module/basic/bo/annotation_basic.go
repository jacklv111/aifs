/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"fmt"
	"io"

	"github.com/google/uuid"
	annoDo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/do"
	basicDo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
)

type RawDataLabelList []annoDo.RawDataLabelDo
type LocalPath string

// annotation base
type AnnotationDataImpl struct {
	DataBaseImpl
	annoDo.AnnotationDo
	RawDataLabelList
}

func (bo *AnnotationDataImpl) GetAnnotationTemplateId() uuid.UUID {
	return bo.AnnotationTemplateId
}
func (bo *AnnotationDataImpl) GetRawDataId() uuid.UUID {
	return bo.AnnotationDo.DataItemId
}

// data base
type DataBaseImpl struct {
	basicDo.DataItemDo
	LocalPath
	io.ReadSeeker
	io.WriterAt
}

func (bo *DataBaseImpl) GetId() uuid.UUID {
	return bo.DataItemDo.ID
}

func (bo *DataBaseImpl) GetLocalPath() string {
	return string(bo.LocalPath)
}

func (bo *DataBaseImpl) GetReader() io.Reader {
	return bo.ReadSeeker
}

func (bo *DataBaseImpl) ResetReader() error {
	_, err := bo.ReadSeeker.Seek(0, io.SeekStart)
	return err
}

func (bo *DataBaseImpl) GetWriterAt() io.WriterAt {
	return bo.WriterAt
}

func (bo *DataBaseImpl) LoadFromBuffer() error {
	return fmt.Errorf("LoadFromBuffer is not implemented")
}

func (bo *DataBaseImpl) LoadFromLocal() error {
	return fmt.Errorf("LoadFromLocal is not implemented")
}
