/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type DataViewStatistics struct {
	ItemCount int32 `json:"itemCount,omitempty"`

	// the number of labels in the annotation data view.
	LabelCount int32 `json:"labelCount,omitempty"`

	// the distribution of labels in the annotation data view.
	LabelDistribution []LabelDistribution `json:"labelDistribution,omitempty"`

	// the total size of the data in the data view.
	TotalDataSize int64 `json:"totalDataSize,omitempty"`
}

type LabelDistribution struct {
	// the label id
	LabelId string `json:"labelId,omitempty"`

	// the number of the label in the data view.
	Count int32 `json:"count,omitempty"`

	// the ratio of the label in the data view.
	Ratio float32 `json:"ratio,omitempty"`
}
