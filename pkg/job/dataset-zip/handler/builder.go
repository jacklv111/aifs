/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package handler

import (
	"fmt"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
)

func BuildHandler(client *aifsclientgo.APIClient, dataView *aifsclientgo.DataViewDetails, datasetDir string) (Handler, error) {
	switch *dataView.ZipFormat {
	case aifsclientgo.IMAGE_CLASSIFICATION:
		return NewClassificationHandler(client, dataView, datasetDir), nil
	case aifsclientgo.RGBD_BOUNDING_BOX_2D_AND_3D:
		return NewRgbdHandler(client, dataView, datasetDir), nil
	case aifsclientgo.IMAGE_SEGMENTATION_MASKS:
		return NewImageSegmentationMasksHandler(client, dataView, datasetDir), nil
	case aifsclientgo.COCO:
		return NewCocoHandler(client, dataView, datasetDir), nil
	case aifsclientgo.OCR:
		return NewOcrHandler(client, dataView, datasetDir), nil
	case aifsclientgo.RAW_DATA_IMAGES:
		return NewRawDataImagesHandler(client, dataView, datasetDir), nil
	case aifsclientgo.SAM:
		return NewSamHandler(client, dataView, datasetDir), nil
	case aifsclientgo.POINTS_3D_ZIP:
		return NewPoints3DHandler(client, dataView, datasetDir), nil
	default:
		return nil, fmt.Errorf("unzip dataset zip view can't handle zip format %s", *dataView.ZipFormat)
	}
}
