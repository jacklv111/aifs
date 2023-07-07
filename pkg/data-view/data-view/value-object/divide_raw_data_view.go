/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type DivideRawDataViewParams struct {
	RawDataViewParamList []EachRawDataViewParam
}

type EachRawDataViewParam struct {
	Name        string
	Description string
	Ratio       int32
}

type DivideRawDataViewResult struct {
	RawDataViewResultList []EachRawDataViewResult
}

type EachRawDataViewResult struct {
	Name       string
	DataViewId string
	ItemCount  int32
}
