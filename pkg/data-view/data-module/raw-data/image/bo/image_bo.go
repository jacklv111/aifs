/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"bytes"
	"io"

	"github.com/google/uuid"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	imagedo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/repo"
	"github.com/jacklv111/common-sdk/utils"
)

// TODO: thumbnail, image score
type ImageRawDataBo struct {
	bo.DataBaseImpl
	imageExtDo   imagedo.ImageExtDo
	imageScoreDo imagedo.ImageScoreDo
}

func (bo *ImageRawDataBo) LoadFromLocal() error {
	// image ext
	meta, err := utils.GetImageMeta(bo.GetLocalPath())
	if err != nil {
		return err
	}
	bo.imageExtDo.Height = meta.Height
	bo.imageExtDo.Width = meta.Width
	bo.imageExtDo.Size = meta.Size

	hashStr, err := utils.GetFileSha256FromFile(bo.GetLocalPath())
	if err != nil {
		return err
	}
	bo.imageExtDo.Sha256 = hashStr

	return nil
}

func (bo *ImageRawDataBo) LoadFromBuffer() error {
	imageBytes, err := io.ReadAll(bo.GetReader())
	if err != nil {
		return err
	}

	image, err := utils.ReadImage(bytes.NewReader(imageBytes))
	if err == nil {
		// image ext
		bo.imageExtDo.Height = int32(image.Bounds().Dy())
		bo.imageExtDo.Width = int32(image.Bounds().Dx())
		bo.imageExtDo.Size = int64(len(imageBytes))
	}

	hashStr, err := utils.GetFileSha256Bytes(imageBytes)
	if err != nil {
		return err
	}
	bo.imageExtDo.Sha256 = hashStr

	return bo.ResetReader()
}

// FixDataItemId 对于已经存在的 raw data，使用其已经存在的 id
//
//	@param imageList
//	@return error
func FixDataItemId(imageList []bo.DataInterface) error {
	var imageExtDoList []imagedo.ImageExtDo
	for _, data := range imageList {
		imageRawDataBo := data.(*ImageRawDataBo)
		imageExtDoList = append(imageExtDoList, imageRawDataBo.imageExtDo)
	}
	existedHashMap, err := repo.ImageRawDataRepo.FindExistedByHash(imagedo.GetHashList(imageExtDoList))
	if err != nil {
		return err
	}
	for _, data := range imageList {
		imageRawDataBo := data.(*ImageRawDataBo)
		id, ok := existedHashMap[imageRawDataBo.imageExtDo.Sha256]
		if !ok {
			continue
		}
		imageRawDataBo.DataItemDo.ID = id
		imageRawDataBo.imageExtDo.ID = id
		imageRawDataBo.imageScoreDo.ID = id
	}
	return nil
}

// CreateBatch 批量插入数据
//
//	@param imageList
//	@return []uuid.UUID
//	@return error
func CreateBatch(imageList []bo.DataInterface) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	var dataItemDoList []basicdo.DataItemDo
	var imageExtDoList []imagedo.ImageExtDo
	var imageScoreDoList []imagedo.ImageScoreDo

	for _, data := range imageList {
		imageRawDataBo := data.(*ImageRawDataBo)
		idList = append(idList, imageRawDataBo.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, imageRawDataBo.DataItemDo)

		imageExtDoList = append(imageExtDoList, imageRawDataBo.imageExtDo)
		imageScoreDoList = append(imageScoreDoList, imageRawDataBo.imageScoreDo)
	}

	err := repo.ImageRawDataRepo.CreateBatch(dataItemDoList, imageExtDoList, imageScoreDoList)
	if err != nil {
		return nil, err
	}

	return idList, nil
}
