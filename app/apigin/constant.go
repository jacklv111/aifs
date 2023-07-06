/*
 * Created on Thu Jul 06 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package apigin

import "math"

const (
	OFFSET_STR = "offset"
	LIMIT_STR  = "limit"
)

const (
	OFFSET_MIN                    = 0
	OFFSET_MAX                    = math.MaxInt
	LIMIT_MIN                     = 1
	LIMIT_MAX                     = 50
	OFFSET_DEFAULT                = 0
	LIMIT_DEFAULT                 = 10
	ANNOTATION_TEMPLATE_TYPE_NAME = "annotationTemplateTypeName"
	ANNOTATION_TEMPLATE_ID_LIST   = "annotationTemplateIdList"
	ANNOTATION_TEMPLATE_ID        = "annotationTemplateId"
	DATA_VIEW_ID                  = "dataViewId"
	DATA_VIEW_ID_LIST             = "dataViewIdList"
	DATA_VIEW_ITEM_ID_LIST        = "dataViewItemIdList"
	DATA_VIEW_NAME                = "dataViewName"
	RAW_DATA_ID_LIST              = "rawDataIdList"
	LABEL_ID                      = "labelId"
	EXCLUDED_ANNO_VIEW_ID         = "excludedAnnotationViewId"
	INCLUDED_ANNO_VIEW_ID         = "includedAnnotationViewId"
	X_FILE_NAME                   = "X-File-Name"
)
