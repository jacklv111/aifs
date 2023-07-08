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
	annotationtemplatetype "github.com/jacklv111/aifs/pkg/annotation-template-type"
	annotempbo "github.com/jacklv111/aifs/pkg/annotation-template/bo"
	datamodule "github.com/jacklv111/aifs/pkg/data-view/data-module"
	catebo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/category/bo"
	cocobo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/coco-type/bo"
	ocrbo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/ocr/bo"
	p3dbo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/points-3d/bo"
	rgbdbo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/rgbd/bo"
	segmasksbo "github.com/jacklv111/aifs/pkg/data-view/data-module/annotation/segmentation-masks/bo"
	basicbo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/bo"
	basicvb "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/value-object"
	"github.com/jacklv111/aifs/pkg/store/manager"
	storevb "github.com/jacklv111/aifs/pkg/store/value-object"
)

// LoadFromRemote 数据文件来自其他节点，通过网络进行上传。
//
//	@param input
//	@param annoTempBo
//	@return []uuid.UUID
//	@return error
func LoadFromRemote(input basicvb.UploadAnnotationParam, annoTempBo annotempbo.AnnotationTemplateBoInterface) ([]uuid.UUID, error) {
	annoList, err := genBoList(input, annoTempBo)
	if err != nil {
		return nil, err
	}

	if err := validate(annoTempBo, annoList, input.RawDataIdChecker); err != nil {
		return nil, err
	}

	// batch save meta
	result, err := saveMetadata(annoTempBo.GetType(), annoList)
	if err != nil {
		return nil, err
	}

	// save the annotation data
	storeParams, err := getStoreParamRemote(annoTempBo.GetType(), annoList)

	if err != nil {
		return nil, err
	}

	err = manager.StoreMgr.Upload(storeParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getStoreParamRemote(annoTempType string, dataList []basicbo.AnnotationData) (storevb.UploadParams, error) {
	var params storevb.UploadParams
	params.DataType = datamodule.ANNOTATION
	for _, data := range dataList {

		switch annoTempType {
		case annotationtemplatetype.COCO_TYPE, annotationtemplatetype.RGBD, annotationtemplatetype.SEGMENTATION_MASKS, annotationtemplatetype.POINTS_3D:
			params.AddItem(data.GetId(), data.GetReader(), annoTempType)

		case annotationtemplatetype.CATEGORY, annotationtemplatetype.OCR:
			// do nothing
		default:
			return params, fmt.Errorf("getStoreParamRemote cant handle annotation type %s", annoTempType)
		}
	}
	return params, nil
}

func genBoList(input basicvb.UploadAnnotationParam, annoTempBo annotempbo.AnnotationTemplateBoInterface) ([]basicbo.AnnotationData, error) {
	annoList := make([]basicbo.AnnotationData, 0)
	switch annoTempBo.GetType() {
	case annotationtemplatetype.SEGMENTATION_MASKS:
		for rawDataId, reader := range input.DataFileMap {
			segMasksBo := segmasksbo.BuildWithReader(uuid.MustParse(rawDataId), reader, input.DataFileNameMap[rawDataId], annoTempBo)
			err := segMasksBo.LoadFromBuffer()
			if err != nil {
				return nil, err
			}
			annoList = append(annoList, segMasksBo)
		}
	case annotationtemplatetype.CATEGORY:
		// 分类的标注比较简单，也比较小，都放在一个文件中。这里是放在 filemeta 中。
		scanner := bufio.NewScanner(input.FileMeta)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.Trim(line, " \t\n\r")
			items := strings.Split(line, " ")
			rawDataId := uuid.MustParse(items[0])
			annoData := uuid.MustParse(items[1])
			cateBo := catebo.BuildWithAnnoData(rawDataId, annoTempBo.GetId(), annoData)
			annoList = append(annoList, cateBo)
		}
	case annotationtemplatetype.OCR:
		scanner := bufio.NewScanner(input.FileMeta)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			items := strings.Split(line, " ")
			rawDataId := uuid.MustParse(items[0])
			annoData := items[1]
			ocrBo := ocrbo.BuildWithAnnoData(rawDataId, annoTempBo.GetId(), annoData)
			annoList = append(annoList, ocrBo)
		}
	case annotationtemplatetype.RGBD:
		for rawDataId, reader := range input.DataFileMap {
			rgbdBo := rgbdbo.BuildWithReader(uuid.MustParse(rawDataId), reader, input.DataFileNameMap[rawDataId], annoTempBo)
			err := rgbdBo.LoadFromBuffer()
			if err != nil {
				return nil, err
			}
			annoList = append(annoList, rgbdBo)
		}
	case annotationtemplatetype.COCO_TYPE:
		for rawDataId, reader := range input.DataFileMap {
			cocoBo := cocobo.BuildWithReader(uuid.MustParse(rawDataId), reader, input.DataFileNameMap[rawDataId], annoTempBo)
			err := cocoBo.LoadFromBuffer()
			if err != nil {
				return nil, err
			}
			annoList = append(annoList, cocoBo)
		}
	case annotationtemplatetype.POINTS_3D:
		for rawDataId, reader := range input.DataFileMap {
			points3dBo := p3dbo.BuildWithReader(uuid.MustParse(rawDataId), reader, input.DataFileNameMap[rawDataId], annoTempBo)
			err := points3dBo.LoadFromBuffer()
			if err != nil {
				return nil, err
			}
			annoList = append(annoList, points3dBo)
		}
	default:
		return nil, fmt.Errorf("genBoList cant handle annotation type %s", annoTempBo.GetType())
	}
	return annoList, nil
}
