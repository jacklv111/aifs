/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type UpdateDatasetZipParams struct {
	Progress float64

	Status string

	RawDataViewId string

	AnnotationViewId string

	TrainRawDataViewId string

	TrainAnnotationViewId string

	ValRawDataViewId string

	ValAnnotationViewId string

	AnnotationTemplateId string
}
