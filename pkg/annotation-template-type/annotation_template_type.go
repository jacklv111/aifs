/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package annotationtemplatetype

import (
	"sort"

	"github.com/jacklv111/common-sdk/utils"
)

// the annotation template type that aifs support
const (
	CATEGORY           = "category"
	COCO_TYPE          = "coco-type"
	OCR                = "ocr"
	RGBD               = "rgbd"
	SEGMENTATION_MASKS = "segmentation-masks"
	POINTS_3D          = "points-3d"
)

type annoTempTypeMap map[string]string

func HasAnnotationTemplateType(key string) bool {
	_, ok := annoTempType[key]
	return ok
}

func GetAnnotationTemplateList(offset int, limit int) []string {
	return annoTempList[utils.Min(offset, size):utils.Min(offset+limit, size)]
}

var annoTempType annoTempTypeMap
var annoTempList []string
var size int

func init() {
	annoTempType = map[string]string{
		CATEGORY:           "",
		COCO_TYPE:          "",
		OCR:                "",
		RGBD:               "",
		SEGMENTATION_MASKS: "",
		POINTS_3D:          "",
	}
	for key := range annoTempType {
		annoTempList = append(annoTempList, key)
	}
	sort.Slice(annoTempList, func(i, j int) bool {
		return annoTempList[i] < annoTempList[j]
	})
	size = len(annoTempList)
}
