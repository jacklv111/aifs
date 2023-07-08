/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/uuid"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	rgbddo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/repo"
	"github.com/jacklv111/common-sdk/flatbuffer/raw-data/go/RawData/Rgbd"
	"github.com/jacklv111/common-sdk/utils"
)

type RgbdRawDataBo struct {
	basicbo.DataBaseImpl
	rgbdExtDo     rgbddo.RgbdExtDo
	imageFilePath string
	depthFilePath string
	calibFilePath string
}

func (bo *RgbdRawDataBo) LoadFromLocal() error {
	// image ext
	meta, err := utils.GetImageMeta(bo.imageFilePath)
	if err != nil {
		return err
	}
	bo.rgbdExtDo.ImageHeight = meta.Height
	bo.rgbdExtDo.ImageWidth = meta.Width
	bo.rgbdExtDo.ImageSize = meta.Size

	// depth ext
	meta, err = utils.GetImageMeta(bo.depthFilePath)
	if err != nil {
		return err
	}
	bo.rgbdExtDo.DepthHeight = meta.Height
	bo.rgbdExtDo.DepthWidth = meta.Width
	bo.rgbdExtDo.DepthSize = meta.Size

	// initial size of the buffer (here 1024 bytes), which will grow automatically if needed
	builder := flatbuffers.NewBuilder(1024)
	imageBytes, err := os.ReadFile(bo.imageFilePath)
	if err != nil {
		return err
	}
	depthBytes, err := os.ReadFile(bo.depthFilePath)
	if err != nil {
		return err
	}
	imageOffset := builder.CreateByteVector(imageBytes)
	depthOffset := builder.CreateByteVector(depthBytes)

	file, err := os.Open(bo.calibFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if err != nil {
		return err
	}

	scanner.Split(bufio.ScanWords)
	extrinsics := make([]float64, 9)
	intrinsics := make([]float64, 9)
	for i := 0; i < 9; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("data %s has wrong calib data", bo.DataItemDo.Name)
		}
		extrinsics[i], err = strconv.ParseFloat(scanner.Text(), 32)

		if err != nil {
			return err
		}
	}
	for i := 0; i < 9; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("data %s has wrong calib data", bo.DataItemDo.Name)
		}
		intrinsics[i], err = strconv.ParseFloat(scanner.Text(), 32)

		if err != nil {
			return err
		}
	}
	Rgbd.CalibStartExtrinsicsVector(builder, 9)
	for i := 8; i >= 0; i-- {
		builder.PrependFloat32(float32(extrinsics[i]))
	}
	extrinsicsOffset := builder.EndVector(9)
	Rgbd.CalibStartIntrinsicsVector(builder, 9)
	for i := 8; i >= 0; i-- {
		builder.PrependFloat32(float32(intrinsics[i]))
	}
	intrinsicsOffset := builder.EndVector(9)

	Rgbd.CalibStart(builder)
	Rgbd.CalibAddExtrinsics(builder, extrinsicsOffset)
	Rgbd.CalibAddIntrinsics(builder, intrinsicsOffset)
	calibOffset := Rgbd.CalibEnd(builder)

	Rgbd.RgbdDataStart(builder)
	Rgbd.RgbdDataAddImage(builder, imageOffset)
	Rgbd.RgbdDataAddDepth(builder, depthOffset)
	Rgbd.RgbdDataAddCalib(builder, calibOffset)
	rgbdData := Rgbd.RgbdDataEnd(builder)

	builder.Finish(rgbdData)
	buf := builder.FinishedBytes()

	os.WriteFile(string(bo.LocalPath), buf, 0777)

	hashStr, err := utils.GetFileSha256Bytes(buf)
	if err != nil {
		return nil
	}
	bo.rgbdExtDo.Sha256 = hashStr

	return nil
}

// FixDataItemId 对于已经存在的 raw data，使用其已经存在的 id
//
//	@param imageList
//	@return error
func FixDataItemId(dataList []basicbo.DataInterface) error {
	var rgbdExtDoList []rgbddo.RgbdExtDo
	for _, data := range dataList {
		rgbdRawDataBo := data.(*RgbdRawDataBo)
		rgbdExtDoList = append(rgbdExtDoList, rgbdRawDataBo.rgbdExtDo)
	}
	existedHashMap, err := repo.RgbdRawDataRepo.FindExistedByHash(rgbddo.GetHashList(rgbdExtDoList))
	if err != nil {
		return err
	}
	for _, data := range dataList {
		rgbdRawDataBo := data.(*RgbdRawDataBo)
		id, ok := existedHashMap[rgbdRawDataBo.rgbdExtDo.Sha256]
		if !ok {
			continue
		}
		rgbdRawDataBo.DataItemDo.ID = id
		rgbdRawDataBo.rgbdExtDo.ID = id
	}
	return nil
}

// CreateBatch 批量插入数据
//
//	@param imageList
//	@return []uuid.UUID
//	@return error
func CreateBatch(imageList []basicbo.DataInterface) ([]uuid.UUID, error) {
	var idList []uuid.UUID
	var dataItemDoList []basicdo.DataItemDo
	var rgbdExtDoList []rgbddo.RgbdExtDo

	for _, data := range imageList {
		rgbdRawDataBo := data.(*RgbdRawDataBo)
		idList = append(idList, rgbdRawDataBo.DataItemDo.ID)
		dataItemDoList = append(dataItemDoList, rgbdRawDataBo.DataItemDo)
		rgbdExtDoList = append(rgbdExtDoList, rgbdRawDataBo.rgbdExtDo)
	}

	err := repo.RgbdRawDataRepo.CreateBatch(dataItemDoList, rgbdExtDoList)
	if err != nil {
		return nil, err
	}

	return idList, nil
}
