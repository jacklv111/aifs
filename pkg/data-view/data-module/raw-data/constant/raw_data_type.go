/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package constant

// raw data type
const (
	IMAGE     = "image"
	TEXT      = "text"
	RGBD      = "rgbd"
	POINTS_3D = "points-3d"
	// video, audio
)

var rawDataTypeMap map[string]string

func init() {
	rawDataTypeMap = map[string]string{
		IMAGE:     "",
		RGBD:      "",
		TEXT:      "",
		POINTS_3D: "",
	}
}

func HasRawDataType(typeStr string) bool {
	_, ok := rawDataTypeMap[typeStr]
	return ok
}

func GetRawDataTypeList() []string {
	var keyList []string
	for k := range rawDataTypeMap {
		keyList = append(keyList, k)
	}
	return keyList
}
