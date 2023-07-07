/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package datamodule

// data view 的类型，即 ai 数据类型
const (
	RAW_DATA    = "raw-data"
	ANNOTATION  = "annotation"
	MODEL       = "model"
	DATASET_ZIP = "dataset-zip"
	ARTIFACT    = "artifact"
)

type viewTypeClass map[string]string

func GetList() []string {
	var keyList []string
	for k := range viewType {
		keyList = append(keyList, k)
	}
	return keyList
}

var viewType viewTypeClass

func init() {
	viewType = map[string]string{
		RAW_DATA:    "",
		ANNOTATION:  "",
		MODEL:       "",
		DATASET_ZIP: "",
		ARTIFACT:    "",
	}
}
