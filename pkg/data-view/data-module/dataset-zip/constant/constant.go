/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package constant

// zip 文件类型
const (
	DATASET_ZIP_FILE = "dataset-zip-file"
)

func IsDatasetZipFile(dataType string) bool {
	return dataType == DATASET_ZIP_FILE
}

// zip 解压后的文件组织形式
const (
	IMAGE_CLASSIFICATION        = "image-classification"
	RGBD_BOUNDING_BOX_2D_AND_3D = "rgbd-bounding-box-2d-and-3d"
	IMAGE_SEGMENTATION_MASKS    = "image-segmentation-masks"
	OCR                         = "ocr"
	COCO                        = "coco"
	RAW_DATA_IMAGES             = "raw-data-images"
	SAM                         = "sam"
	POINTS_3D_ZIP               = "points-3d-zip"
)

var zipFormatMap map[string]string

func IsZipFormat(zipFormat string) bool {
	_, ok := zipFormatMap[zipFormat]
	return ok
}

func GetZipFormatList() []string {
	res := make([]string, 0)
	for k := range zipFormatMap {
		res = append(res, k)
	}
	return res
}

func init() {
	zipFormatMap = map[string]string{
		IMAGE_CLASSIFICATION:        "",
		RGBD_BOUNDING_BOX_2D_AND_3D: "",
		IMAGE_SEGMENTATION_MASKS:    "",
		OCR:                         "",
		COCO:                        "",
		RAW_DATA_IMAGES:             "",
		SAM:                         "",
		POINTS_3D_ZIP:               "",
	}
}
