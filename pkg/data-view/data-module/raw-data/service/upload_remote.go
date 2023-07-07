/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/google/uuid"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/basic/service"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/constant"
	imagebo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/bo"
	p3dbo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/bo"
	p3dvb "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/points-3d/value-object"
	rgbdbo "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/bo"
	rgbdvb "github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/rgbd/value-object"
	"github.com/jacklv111/aifs/pkg/store/manager"
	storevb "github.com/jacklv111/aifs/pkg/store/value-object"
)

func LoadFromRemote(input basicvb.UploadRawDataParam, rawDataType string) ([]uuid.UUID, error) {
	rawDataList, err := genBoList(input, rawDataType)
	if err != nil {
		return nil, err
	}

	// batch save meta
	result, err := saveMetadata(rawDataType, rawDataList)
	if err != nil {
		return nil, err
	}

	// save the annotation data
	storeParams := getStoreParamRemote(rawDataType, rawDataList)

	err = manager.StoreMgr.Upload(storeParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func genBoList(input basicvb.UploadRawDataParam, rawDataType string) ([]basicbo.DataInterface, error) {
	rawDataList := make([]basicbo.DataInterface, 0)
	switch rawDataType {
	case constant.IMAGE:
		for fileName, reader := range input.DataFileMap {
			imageBo := imagebo.BuildWithReader(reader, fileName)
			err := imageBo.LoadFromBuffer()
			if err != nil {
				return nil, err
			}
			rawDataList = append(rawDataList, imagebo)
		}

	case constant.RGBD:
		scanner := bufio.NewScanner(input.FileMeta)
		// key: file name; value: raw data meta
		metaMap := make(map[string]rgbdvb.RgbdRawDataMeta)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.Trim(line, " \t\n\r")
			meta := rgbdvb.RgbdRawDataMeta{}
			err := meta.ParseFromString(line)
			if err != nil {
				return nil, err
			}
			metaMap[meta.FileName] = meta
		}
		for fileName, reader := range input.DataFileMap {
			rgbdBo := rgbdbo.BuildWithReader(reader, fileName, metaMap[fileName])
			rawDataList = append(rawDataList, rgbdBo)
		}
	case constant.POINTS_3D:
		scanner := bufio.NewScanner(input.FileMeta)
		// key: file name; value: raw data meta
		metaMap := make(map[string]p3dvb.Points3DMeta)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.Trim(line, " \t\n\r")
			meta := p3dvb.Points3DMeta{}
			err := meta.ParseFromString(line)
			if err != nil {
				return nil, err
			}
			metaMap[meta.FileName] = meta
		}
		for fileName, reader := range input.DataFileMap {
			point3DBo := p3dbo.BuildWithReader(reader, fileName, metaMap[fileName])
			rawDataList = append(rawDataList, point3DBo)
		}
	default:
		return nil, fmt.Errorf("genBoList cant handle annotation type %s", rawDataType)
	}
	return rawDataList, nil
}

func getStoreParamRemote(rawDataType string, rawDataList []basicbo.DataInterface) storevb.UploadParams {
	var params storevb.UploadParams
	params.DataType = service.RAW_DATA
	for _, data := range rawDataList {
		params.AddItem(data.GetId(), data.GetReader(), rawDataType)
	}
	return params
}
